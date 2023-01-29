import { get } from 'svelte/store';
import { environment } from './stores';

// Calculate the min and max values of signed integers with n bits
export const intMin = bits => Math.pow(2, bits - 1) * -1;
export const intMax = bits => Math.pow(2, bits - 1) - 1;
export const uintMax = bits => Math.pow(2, bits) - 1;

// Boundaries for some ubiquitous integer types
export const int32 = [ intMin(32), intMax(32) ];
export const int64 = [ intMin(64), intMax(64) ];
export const uint64 = [ 0, uintMax(64) ];

// Get a value from an object with a JSON path, from Webdesq core
export function resolveKeypath(object, path) {
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

// Set a value in an object with a JSON path, from Webdesq core
export function setValue(object, path, value) {
  const parts = path.split('.').flatMap(part => {
    let indexMatch = part.match(/\[\d+\]/g);
    if (indexMatch) {
      // Convert strings to numbers
      const indexes = indexMatch.map(index => Number(index.slice(1, -1)));
      const base = part.slice(0, part.indexOf(indexMatch[0]));
      return base.length ? [ base, ...indexes ] : indexes;
    }
    indexMatch = part.match(/^\d+$/g);
    if (indexMatch) {
      // Convert strings to numbers
      const indexes = indexMatch.map(index => Number(index.slice(1, -1)));
      const base = part.slice(0, part.indexOf(indexMatch[0]));
      return base.length ? [ base, ...indexes ] : indexes;
    }
    return part;
  });
  let result = object;
  while (parts.length) {
    const part = parts.shift();
    if (!parts.length) {
      // No parts left, we can set the value
      result[part] = value;
      break;
    }
    if (!result[part]) {
      // Default value if none is found
      result[part] = (typeof parts[0] === 'number') ? [] : {};
    }
    result = result[part];
  }
  return object;
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

// https://stackoverflow.com/a/14794066
export function isInt(value) {
  if (isNaN(value)) {
    return false;
  }
  const x = parseFloat(value);
  return (x | 0) === x;
}

export function randInt(min, max) {
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
