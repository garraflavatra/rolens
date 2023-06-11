import { Hosts } from '$wails/go/app/App';
import { writable } from 'svelte/store';
import applicationInited from './inited';

const { set, subscribe } = writable();

const update = async() => set(await Hosts());
applicationInited.defer(update);

const hosts = { update, subscribe };
export default hosts;
