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
