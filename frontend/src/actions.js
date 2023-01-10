export function input(node, { json } = { json: false }) {
  const handleInput = () => {
    if (json) {
      try {
        JSON.parse(node.value);
        node.classList.remove('invalid');
      }
      catch {
        node.classList.add('invalid');
      }
    }
  };

  const handleFocus = () => {
    node.select();
  };

  node.addEventListener('focus', handleFocus);
  node.addEventListener('input', handleInput);

  return {
    destroy: () => {
      node.removeEventListener('focus', handleFocus);
      node.removeEventListener('input', handleInput);
    },
  };
}
