import { Environment } from '$wails/go/app/App.js';
import { writable } from 'svelte/store';

const { set, subscribe } = writable({});

async function reload() {
  const newEnv = await Environment();
  set(newEnv);
  return newEnv;
}

reload();

subscribe(env => {
  // @ts-ignore
  document.body.dataset.platform = env?.platform;
});

const environment = { reload, subscribe };
export default environment;
