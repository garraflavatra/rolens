import { derived, writable } from 'svelte/store';
import { Environment, Settings, UpdateSettings, UpdateViewStore, Views } from '../wailsjs/go/app/App';

export const busy = (() => {
  const { update, subscribe } = writable(0);

  subscribe(isBusy => {
    if (isBusy) {
      document.body.classList.add('busy');
    }
    else {
      document.body.classList.remove('busy');
    }
  });

  return {
    start: () => update(v => ++v),
    end: () => update(v => --v),
    subscribe,
  };
})();

export const contextMenu = (() => {
  const { set, subscribe } = writable();

  return {
    show: (evt, menu) => set(Object.keys(menu || {}).length ? {
      position: [ evt.clientX, evt.clientY ],
      items: menu,
    } : undefined),
    hide: () => set(undefined),
    subscribe,
  };
})();

export const connections = writable({});

export const applicationSettings = (() => {
  const { set, subscribe } = writable({});
  const reload = async() => {
    const newSettings = await Settings();
    set(newSettings);
    return newSettings;
  };

  let skipUpdate = true;

  reload();
  subscribe(newSettings => {
    if (skipUpdate) {
      skipUpdate = false;
      return;
    }
    UpdateSettings(JSON.stringify(newSettings || {}));
  });

  return { reload, set, subscribe };
})();

export const environment = (() => {
  const { set, subscribe } = writable({});
  const reload = async() => {
    const newEnv = await Environment();
    set(newEnv);
    return newEnv;
  };

  reload();
  return { reload, subscribe };
})();

export const views = (() => {
  const { set, subscribe } = writable({});

  const reload = async() => {
    const newViewStore = await Views();
    set(newViewStore);
    return newViewStore;
  };

  const forCollection = (hostKey, dbKey, collKey) => {
    let allViews;
    subscribe(v => allViews = v)();

    const entries = Object.entries(allViews).filter(v => (
      v[0] === 'list' || (
        v[1].host === hostKey &&
        v[1].database === dbKey &&
        v[1].collection === collKey
      )
    ));

    return Object.fromEntries(entries);
  };

  let skipUpdate = true;

  reload();
  subscribe(newViewStore => {
    if (skipUpdate) {
      skipUpdate = false;
      return;
    }
    UpdateViewStore(JSON.stringify(newViewStore));
  });

  return { reload, set, subscribe, forCollection };
})();

export const applicationInited = derived([ environment, applicationSettings ], ([ env, settings ], set) => {
  if (env && settings) {
    set(true);

    // Remove loading spinner.
    // document.getElementById('app-loading')?.remove();
  }
}, false);
