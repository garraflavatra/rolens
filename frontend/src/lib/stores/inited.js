import { derived } from 'svelte/store';
import environment from './environment';
import applicationSettings from './settings';

const applicationInited = derived([ environment, applicationSettings ], ([ env, settings ], set) => {
  if (env && settings) {
    set(true);

    // Remove loading spinner.
    document.getElementById('app-loading')?.remove();
  }
}, false);

export default applicationInited;
