import { get, writable } from 'svelte/store';
import { UpdateViewStore, Views } from '../../../wailsjs/go/app/App';

const { set, subscribe } = writable({});
let skipUpdate = true;

async function reload() {
  const newViewStore = await Views();
  set(newViewStore);
  return newViewStore;
}

function forCollection(hostKey, dbKey, collKey)  {
  const allViews = get({ subscribe });
  const entries = Object.entries(allViews).filter(v => (
    v[0] === 'list' || (
      v[1].host === hostKey &&
        v[1].database === dbKey &&
        v[1].collection === collKey
    )
  ));

  return Object.fromEntries(entries);
}

reload();
subscribe(newViewStore => {
  if (skipUpdate) {
    skipUpdate = false;
    return;
  }
  UpdateViewStore(JSON.stringify(newViewStore));
});

const views = { reload, set, subscribe, forCollection };
export default views;
