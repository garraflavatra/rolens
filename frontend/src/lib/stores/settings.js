import { Settings, UpdateSettings } from '$wails/go/app/App';
import { writable } from 'svelte/store';

const { set, subscribe } = writable({});
let skipUpdate = true;

async function reload() {
  const newSettings = await Settings();
  set(newSettings);
  return newSettings;
}

reload();
subscribe(newSettings => {
  if (skipUpdate) {
    skipUpdate = false;
    return;
  }
  UpdateSettings(JSON.stringify(newSettings || {}));
});

const applicationSettings = { reload, set, subscribe };
export default applicationSettings;
