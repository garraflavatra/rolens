import { get, writable } from "svelte/store";
import applicationInited from "./inited";
import { startProgress } from "$lib/progress";
import windowTitle from "./windowtitle";
import {
  DropCollection,
  DropDatabase,
  Hosts,
  OpenCollection,
  OpenConnection,
  OpenDatabase,
  RemoveHost,
  TruncateCollection
} from "$wails/go/app/App";

const { set, subscribe } = writable({});
const getValue = () => get({ subscribe });

async function refresh() {
  const hosts = await Hosts();
  const hostTree = getValue();

  for (const [ hostKey, hostDetails ] of Object.entries(hosts)) {
    hostTree[hostKey] = hostTree[hostKey] || {};
    const host = hostTree[hostKey];
    host.key = hostKey;
    host.name = hostDetails.name;
    host.uri = hostDetails.uri;

    host.open = async function() {
      const progress = startProgress(`Connecting to "${hostKey}"…`);
      const { databases, status, systemInfo } = await OpenConnection(hostKey);
      host.status = status;
      host.systemInfo = systemInfo;
      host.databases = host.databases || {};

      for (const dbKey of databases || []) {
        host.databases[dbKey] = host.databases[dbKey] || {};
        const database = host.databases[dbKey];
        database.key = dbKey;
        database.collections = database.collections || {};

        database.open = async function() {
          const progress = startProgress(`Opening database "${dbKey}"…`);
          const { collections, stats } = await OpenDatabase(hostKey, dbKey);
          database.stats = stats;

          for (const collKey of collections || []) {
            database.collections[collKey] = database.collections[collKey] || {};
            const collection = database.collections[collKey];
            collection.key = collKey;

            collection.open = async function() {
              const progress = startProgress(`Opening database "${dbKey}"…`);
              const stats = await OpenCollection(hostKey, dbKey, collKey);
              collection.stats = stats;
              await refresh();
              progress.end();
            }

            collection.truncate = async function() {
              const progress = startProgress(`Truncating collection "${collKey}"…`);
              await TruncateCollection(hostKey, dbKey, collKey);
              await refresh();
              progress.end();
            };

            collection.drop = async function() {
              const progress = startProgress(`Dropping collection "${collKey}"…`);
              const success = await DropCollection(hostKey, dbKey, collKey);

              if (success) {
                await refresh();
              }

              progress.end();
            };
          };

          await refresh();
          progress.end();
          windowTitle.setSegments(dbKey, host.name, 'Rolens');
        };

        database.drop = async function() {
          const progress = startProgress(`Dropping database "${dbKey}"…`);
          const success = await DropDatabase(hostKey, dbKey);

          if (success) {
            await refresh();
          }

          progress.end();
        };
      }

      await refresh();
      progress.end();
    };

    host.remove = async function() {
      await RemoveHost(hostKey);
      await refresh();
    };
  }

  set(hostTree);
}

applicationInited.defer(refresh);

const hostTree = { refresh, subscribe, get: getValue };
export default hostTree;
