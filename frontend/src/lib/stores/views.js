import dialogs from '$lib/dialogs';
import ViewConfigDialog from '$organisms/connection/collection/dialogs/viewconfig.svelte';
import { UpdateViewStore, Views } from '$wails/go/app/App';
import { get, writable } from 'svelte/store';

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

function openConfig(collection, firstItem = {}) {
  const dialog = dialogs.new(ViewConfigDialog, { collection, firstItem });
  return dialog;
}

reload();
subscribe(newViewStore => {
  if (skipUpdate) {
    skipUpdate = false;
    return;
  }
  UpdateViewStore(JSON.stringify(newViewStore));
});

const views = { reload, set, subscribe, forCollection, openConfig };
export default views;
