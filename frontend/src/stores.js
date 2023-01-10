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
