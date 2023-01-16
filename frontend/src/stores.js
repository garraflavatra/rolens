import { writable } from 'svelte/store';

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
    show: (evt, menu) => set(menu ? {
      position: [ evt.clientX, evt.clientY ],
      items: menu,
    } : undefined),
    hide: () => set(undefined),
    subscribe,
  };
})();

export const connections = writable({});
