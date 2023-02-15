import environment from '$lib/stores/environment';
import { get } from 'svelte/store';

export function controlKeyDown(event) {
  const env = get(environment);
  // @ts-ignore
  if (env?.platform === 'darwin') {
    return event?.metaKey;
  }
  else {
    return event?.ctrlKey;
  }
}
