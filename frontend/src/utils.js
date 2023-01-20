import { get } from 'svelte/store';
import { environment } from './stores';

export function resolveKeypath(object, path) {
  // Get a value from an object with a JSON path, from Webdesq core

  const parts = path.split('.').flatMap(part => {
    const indexMatch = part.match(/\[\d+\]/g);
    if (indexMatch) {
      // Convert strings to numbers
      const indexes = indexMatch.map(index => Number(index.slice(1, -1)));
      const base = part.slice(0, part.indexOf(indexMatch[0]));
      return base.length ? [ base, ...indexes ] : indexes;
    }
    return part;
  });

  let result = object;
  while (result && parts.length) {
    result = result[parts.shift()];
  }

  return result;
}

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
