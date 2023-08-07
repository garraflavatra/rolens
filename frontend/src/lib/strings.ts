export function capitalise(string = ''): string {
  const capitalised = string.charAt(0).toUpperCase() + string.slice(1);
  return capitalised;
}

export function jsonLooseParse<T>(json: string): T {
  const obj: T = new Function(`return (${json})`)();
  return obj;
}

export function convertLooseJson(json: any) {
  const j = JSON.stringify(jsonLooseParse(json));
  return j;
}

export function looseJsonIsValid(json: string): boolean {
  try {
    jsonLooseParse(json);
    return true;
  }
  catch {
    return false;
  }
}

export function stringCouldBeID(string: string) {
  return /^[a-zA-Z0-9_-]{1,}$/.test(string);
}
