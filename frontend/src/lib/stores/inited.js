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
    listeners.push(listener)
  }
};

const { subscribe } = derived([ environment, applicationSettings ], ([ env, settings ], set) => {
  if (alreadyInited) {
    return;
  }

  if (env && settings) {
    set(true);
    alreadyInited = true;

    // Remove loading spinner.
    document.getElementById('app-loading')?.remove();

    // Call hooks
    listeners.forEach(l => l());
  }
}, false);

const applicationInited = { defer, subscribe };

export default applicationInited;
