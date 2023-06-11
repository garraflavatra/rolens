import { startProgress } from '$lib/progress';
import {
  DropCollection,
  DropDatabase,
  Hosts,
  OpenCollection,
  OpenConnection,
  OpenDatabase,
  RemoveHost,
  RenameCollection,
  TruncateCollection
} from '$wails/go/app/App';
import { EnterText } from '$wails/go/ui/UI';
import { get, writable } from 'svelte/store';
import applicationInited from './inited';
import windowTitle from './windowtitle';
import dialogs from '$lib/dialogs';
import HostDetailDialog from '$organisms/connection/host/dialogs/hostdetail.svelte';

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
      const { databases: dbNames, status, systemInfo } = await OpenConnection(hostKey);
      host.status = status;
      host.systemInfo = systemInfo;
      host.databases = host.databases || {};

      for (const dbKey of dbNames) {
        host.databases[dbKey] = host.databases[dbKey] || {};
      }

      for (const [ dbKey, database ] of Object.entries(host.databases)) {
        database.key = dbKey;
        database.collections = database.collections || {};

        database.open = async function() {
          const progress = startProgress(`Opening database "${dbKey}"…`);
          const { collections: collNames, stats } = await OpenDatabase(hostKey, dbKey);
          database.stats = stats;

          for (const collKey of collNames) {
            database.collections[collKey] = database.collections[collKey] || {};
          }

          for (const [ collKey, collection ] of Object.entries(database.collections)) {
            collection.key = collKey;

            collection.open = async function() {
              const progress = startProgress(`Opening database "${dbKey}"…`);
              const stats = await OpenCollection(hostKey, dbKey, collKey);
              collection.stats = stats;
              await refresh();
              progress.end();
            };

            collection.rename = async function() {
              const newCollKey = await EnterText('Rename collection', `Enter a new name for collection ${collKey}.`, collKey);
              if (newCollKey && (newCollKey !== collKey)) {
                const progress = startProgress(`Renaming collection "${collKey}" to "${newCollKey}"…`);
                const ok = await RenameCollection(hostKey, dbKey, collKey, newCollKey);
                await refresh();
                progress.end();
                return ok;
              }
            };

            collection.export = async function() {
              const exportInfo = {
                type: 'export',
                filetype: 'json',
                hostKey,
                dbKey,
                collKeys: [ collKey ],
              };
            };

            collection.dump = async function() {
              const exportInfo = {
                type: 'dump',
                filetype: 'bson',
                hostKey,
                dbKey,
                collKeys: [ collKey ],
              };
            };

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
          }

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

        database.newCollection = async function() {
          const name = await EnterText('Create a collection', 'Note: collections in MongoDB do not exist until they have at least one item. Your new collection will not persist on the server; fill it to have it created.', '');
          if (name) {
            database.collections[name] = {};
            await refresh();
          }
        };
      }

      host.newDatabase = async function() {
        const name = await EnterText('Create a database', 'Enter the database name. Note: databases in MongoDB do not exist until they have a collection and an item. Your new database will not persist on the server; fill it to have it created.', '');
        if (name) {
          host.databases[name] = {};
          await refresh();
        }
      };

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

function newHost() {
  const dialog = dialogs.new(HostDetailDialog, { hostKey: '' });
  return new Promise(resolve => {
    dialog.$on('close', () => {
      refresh().then(resolve);
    });
  });
}

applicationInited.defer(refresh);

const hostTree = { refresh, subscribe, get: getValue, newHost };
export default hostTree;
