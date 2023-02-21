import { StartProgressBar, StopProgressBar } from '$wails/go/ui/UI';
import { writable } from 'svelte/store';

const { update, subscribe } = writable(0);

let timer;
let progressBarShown = false;
subscribe(isBusy => {
  if (isBusy) {
    document.body.classList.add('busy');
    if (!progressBarShown) {
      progressBarShown = true;
      timer = setTimeout(() => StartProgressBar(''), 100);
    }
  }
  else {
    if (timer) {
      clearTimeout(timer);
      timer = undefined;
    }
    progressBarShown = false;
    document.body.classList.remove('busy');
    StopProgressBar();
  }
});

const busy = {
  start: () => update(v => ++v),
  end: () => update(v => --v),
  subscribe,
};

export default busy;
