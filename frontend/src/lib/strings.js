export function capitalise(string = '') {
  const capitalised = string.charAt(0).toUpperCase() + string.slice(1);
  return capitalised;
}

export function jsonLooseParse(json) {
  const obj = new Function(`return (${json})`)();
  return obj;
}

export function convertLooseJson(json) {
  const j = JSON.stringify(jsonLooseParse(json));
  return j;
}

export function looseJsonIsValid(json) {
  try {
    jsonLooseParse(json);
    return true;
  }
  catch {
    return false;
  }
}

export function stringCouldBeID(string) {
  if (looseJsonIsValid(string)) {
    return false;
  }
  return /^[a-zA-Z0-9_-]{1,}$/.test(string);
}
