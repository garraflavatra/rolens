import { WindowSetTitle } from '$wails/runtime/runtime.js';
import { writable } from 'svelte/store';

const { set, subscribe } = writable('Rolens');

subscribe(newTitle => WindowSetTitle(newTitle));

const windowTitle = {
  set,
  setSegments: (...segments) => set(segments.map(s => s.trim()).join(' — ')),
  subscribe,
};

export default windowTitle;
