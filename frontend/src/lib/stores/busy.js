import { writable } from 'svelte/store';

const { update, subscribe } = writable(0);

subscribe(isBusy => {
  if (isBusy) {
    document.body.classList.add('busy');
  }
  else {
    document.body.classList.remove('busy');
  }
});

const busy = {
  start: () => update(v => ++v),
  end: () => update(v => --v),
  subscribe,
};

export default busy;
