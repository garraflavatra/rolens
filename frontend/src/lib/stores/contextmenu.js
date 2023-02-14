import { writable } from 'svelte/store';

const { set, subscribe } = writable();

const contextMenu = {
  show: (evt, menu) => set(Object.keys(menu || {}).length ? {
    position: [ evt.clientX, evt.clientY ],
    items: menu,
  } : undefined),
  hide: () => set(undefined),
  subscribe,
};

export default contextMenu;
