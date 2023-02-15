export function capitalise(string = '') {
  const capitalised = string.charAt(0).toUpperCase() + string.slice(1);
  return capitalised;
}
