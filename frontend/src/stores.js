import { writable } from 'svelte/store';

export const busy = writable(false);
busy.subscribe(isBusy => {
  if (isBusy) {
    document.body.classList.add('busy');
  }
  else {
    document.body.classList.remove('busy');
  }
});

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
