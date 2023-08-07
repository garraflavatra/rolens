export function randInt(min: number, max: number) {
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
