import { get, writable } from 'svelte/store';
import { RemoveQuery, SavedQueries, SaveQuery, UpdateQueries } from '../../../wailsjs/go/app/App';

const { set, subscribe } = writable({});
let skipUpdate = true;

async function reload() {
  const newViewStore = await SavedQueries();
  set(newViewStore);
  return newViewStore;
}

function forCollection(hostKey, dbKey, collKey)  {
  const allViews = get({ subscribe });
  const entries = Object.entries(allViews).filter(v => (
    v[1].hostKey === hostKey &&
    v[1].dbKey === dbKey &&
    v[1].collKey === collKey
  ));

  return Object.fromEntries(entries);
}

async function create(query) {
  const newId = await SaveQuery(JSON.stringify(query));
  await reload();
  return newId;
}

async function remove(id) {
  await RemoveQuery(id);
  await reload();
}

reload();
subscribe(newViewStore => {
  if (skipUpdate) {
    skipUpdate = false;
    return;
  }
  UpdateQueries(JSON.stringify(newViewStore));
});

const queries = { reload, create, remove, set, subscribe, forCollection };
export default queries;
