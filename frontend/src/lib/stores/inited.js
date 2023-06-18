import { derived } from 'svelte/store';
import environment from './environment';
import applicationSettings from './settings';

let alreadyInited = false;

const listeners = [];
const defer = listener => {
  if (alreadyInited) {
    listener();
  }
  else {
    listeners.push(listener);
  }
};

const { subscribe } = derived([ environment, applicationSettings ], ([ env, settings ], set) => {
  if (alreadyInited) {
    return;
  }
  else if (env && settings) {
    Promise.all(listeners.map(l => l())).then(() => {
      set(true);
      alreadyInited = true;
      document.getElementById('app-loading')?.remove();
    });
  }
}, false);

const applicationInited = { defer, subscribe };

export default applicationInited;
