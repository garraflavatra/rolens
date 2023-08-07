import { Environment } from '$wails/go/app/App.js';
import { writable } from 'svelte/store';

const { set, subscribe } = writable({});

async function reload() {
  const newEnv = await Environment();
  set(newEnv);
  return newEnv;
}

reload();

const environment = { reload, subscribe };
export default environment;
