import { canBeObjectId, int32, int64, isInt, uint64 } from './utils';

export function input(node, { autofocus, type, onValid, onInvalid, mandatory } = {
  autofocus: false,
  type: '',
  onValid: () => 0,
  onInvalid: () => 0,
  mandatory: false,
}) {

  const getMessage = () => {
    const checkInteger = () => (isInt(node.value) ? false : 'Value must be an integer');
    const checkNumberBoundaries = boundaries => {
      if (node.value < boundaries[0]) {
        return `Input is too low for type ${type}`;
      }
      else if (node.value > boundaries[1]) {
        return `Input is too high for type ${type}`;
      }
      else {
        return true;
      }
    };

    switch (type) {
      case 'json':
        try {
          JSON.parse(node.value);
          return false;
        }
        catch {
          return 'Invalid JSON';
        }

      case 'int': // int32
        return checkInteger() || checkNumberBoundaries(int32);

      case 'long': // int64
        return checkInteger() || checkNumberBoundaries(int64);

      case 'uint64':
        return checkInteger() || checkNumberBoundaries(uint64);

      case 'string':
        if (mandatory && (!node.value)) {
          return 'This field cannot empty';
        }
        return false;

      case 'objectid':
        return !canBeObjectId(node.value) && 'Invalid string representation of an ObjectId';

      case 'double':
      case 'decimal':
      default:
        return false;
    }
  };

  const handleInput = () => {
    const invalid = getMessage();
    if (invalid) {
      node.classList.add('invalid');
      node.setCustomValidity(invalid);
      node.reportValidity();
      onInvalid?.();
    }
    else {
      node.classList.remove('invalid');
      node.setCustomValidity('');
      node.reportValidity();
      onValid?.();
    }
  };

  const handleFocus = () => {
    node.select();
  };

  node.addEventListener('focus', handleFocus);
  node.addEventListener('input', handleInput);

  if (autofocus) {
    node.focus();
  }

  return {
    destroy: () => {
      node.removeEventListener('focus', handleFocus);
      node.removeEventListener('input', handleInput);
    },
  };
}
