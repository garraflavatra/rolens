import { writable } from 'svelte/store';
import { Environment } from '../../../wailsjs/go/app/App';

const { set, subscribe } = writable({});

async function reload() {
  const newEnv = await Environment();
  set(newEnv);
  return newEnv;
}

reload();

const environment = { reload, subscribe };
export default environment;
