import { writable } from 'svelte/store';
import { Environment } from '../wailsjs/runtime/runtime';
import { Settings, UpdateSettings } from '../wailsjs/go/app/App';

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

  reload();
  subscribe(newSettings => {
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
