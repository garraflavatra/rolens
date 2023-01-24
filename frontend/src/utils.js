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

function randInt(min, max) {
  return Math.round(Math.random() * (max - min) + min);
}

export function randomString(length = 12) {
  const chars = 'qwertyuiopasdfghjklzxcvbnm1234567890';
  let output = '';

  Array(length).fill('').forEach(() => {
    output += chars[randInt(0, chars.length - 1)];
  });

  return output;
}

export function isBsonBuiltin(value) {
  return (
    (typeof value === 'object') &&
    (value !== null) &&
    (typeof value._bsontype === 'string') &&
    (typeof value.inspect === 'function')
  );
}

export function isDate(value) {
  return (value instanceof Date) && !isNaN(value.getTime());
}
