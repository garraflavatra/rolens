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
