import dialogs from '$lib/dialogs';
import { startProgress } from '$lib/progress';
import { EnterText } from '$wails/go/ui/UI';
import { get, writable } from 'svelte/store';
import applicationInited from './inited';
import queries from './queries';
import windowTitle from './windowtitle';

import ExportDialog from '$organisms/connection/collection/dialogs/export.svelte';
import IndexDetailDialog from '$organisms/connection/collection/dialogs/indexdetail.svelte';
import QueryChooserDialog from '$organisms/connection/collection/dialogs/querychooser.svelte';
import DumpDialog from '$organisms/connection/database/dialogs/dump.svelte';
import HostDetailDialog from '$organisms/connection/host/dialogs/hostdetail.svelte';

import {
  CreateIndex,
  DropCollection,
  DropDatabase,
  DropIndex,
  GetIndexes,
  Hosts,
  OpenCollection,
  OpenConnection,
  OpenDatabase,
  PerformDump,
  PerformFindExport,
  RemoveHost,
  RenameCollection,
  TruncateCollection
} from '$wails/go/app/App';

const { set, subscribe } = writable({});
const getValue = () => get({ subscribe });
let hostTreeInited = false;

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

      if (!dbNames) {
        return;
      }

      for (const dbKey of dbNames.sort((a, b) => a.localeCompare(b))) {
        host.databases[dbKey] = host.databases[dbKey] || {};
      }

      for (const [ dbKey, database ] of Object.entries(host.databases)) {
        database.key = dbKey;
        database.hostKey = hostKey;
        database.collections = database.collections || {};

        database.open = async function() {
          const progress = startProgress(`Opening database "${dbKey}"…`);
          const { collections: collNames, stats } = await OpenDatabase(hostKey, dbKey);
          database.stats = stats;

          if (!collNames) {
            return;
          }

          for (const collKey of collNames.sort((a, b) => a.localeCompare(b))) {
            database.collections[collKey] = database.collections[collKey] || {};
          }

          for (const [ collKey, collection ] of Object.entries(database.collections)) {
            collection.key = collKey;
            collection.dbKey = dbKey;
            collection.hostKey = hostKey;
            collection.indexes = collection.indexes || [];

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

            collection.export = function(query) {
              const dialog = dialogs.new(ExportDialog, { collection, query });

              return new Promise(resolve => {
                dialog.$on('export', async event => {
                  const success = await PerformFindExport(hostKey, dbKey, collKey, JSON.stringify(event.detail.exportInfo));
                  if (success) {
                    dialog.$close();
                    resolve();
                  }
                });
              });
            };

            collection.dump = function() {
              const dialog = dialogs.new(DumpDialog, { info: {
                hostKey,
                dbKey,
                collKeys: [ collKey ],
              } });

              return new Promise(resolve => {
                dialog.$on('dump', async event => {
                  const success = await PerformDump(JSON.stringify(event.detail.info));
                  if (success) {
                    dialog.$close();
                    resolve();
                  }
                });
              });
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

            collection.getIndexes = async function() {
              const progress = startProgress(`Retrieving indexes of "${collKey}"…`);
              collection.indexes = [];
              const indexes = await GetIndexes(hostKey, dbKey, collKey);

              for (const indexDetails of indexes) {
                const index = {
                  name: indexDetails.name,
                  background: indexDetails.background || false,
                  unique: indexDetails.unique || false,
                  sparse: indexDetails.sparse || false,
                  model: indexDetails.model,
                };

                index.drop = async function() {
                  const progress = startProgress(`Dropping index ${index.name}…`);
                  const hasBeenDropped = await DropIndex(hostKey, dbKey, collKey, index.name);
                  progress.end();
                  return hasBeenDropped;
                };

                collection.indexes.push(index);
              }

              progress.end();
              return collection.indexes;
            };

            collection.getIndexByName = function(indesName) {
              return collection.indexes.find(idx => idx.name = indesName);
            };

            collection.newIndex = function() {
              const dialog = dialogs.new(IndexDetailDialog, { collection });

              return new Promise(resolve => {
                dialog.$on('create', async event => {
                  const progress = startProgress('Creating index…');
                  const newIndexName = await CreateIndex(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(event.detail.index));

                  if (newIndexName) {
                    dialog.$close();
                  }

                  progress.end();
                  resolve(newIndexName);
                });
              });
            };

            collection.openQueryChooser = function(queryToSave = undefined) {
              const dialog = dialogs.new(QueryChooserDialog, { collection, queryToSave });

              return new Promise(resolve => {
                dialog.$on('select', async event => {
                  dialog.$close();
                  resolve(event.detail.query);
                });

                dialog.$on('create', async event => {
                  const ok = await queries.create(event.detail.query);
                  if (ok) {
                    dialog.$close();
                    resolve(event.detail.query);
                  }
                });
              });
            };
          }

          await refresh();
          progress.end();
          windowTitle.setSegments(dbKey, host.name, 'Rolens');
        };

        database.dump = function() {
          const dialog = dialogs.new(DumpDialog, { info: { hostKey, dbKey } });

          return new Promise(resolve => {
            dialog.$on('dump', async event => {
              const success = await PerformDump(JSON.stringify(event.detail.info));
              if (success) {
                dialog.$close();
                resolve();
              }
            });
          });
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

  hostTreeInited = true;
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

const hostTree = {
  refresh,
  subscribe,
  get: getValue,
  newHost,
  hasBeenInited: () => hostTreeInited,
};

export default hostTree;
