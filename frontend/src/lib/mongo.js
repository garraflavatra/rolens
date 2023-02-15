import { ObjectId } from 'bson';

// Calculate the min and max values of (un)signed integers with n bits
export const intMin = bits => Math.pow(2, bits - 1) * -1;
export const intMax = bits => Math.pow(2, bits - 1) - 1;
export const uintMax = bits => Math.pow(2, bits) - 1;

// Boundaries for some ubiquitous integer types
export const int32 = [ intMin(32), intMax(32) ];
export const int64 = [ intMin(64), intMax(64) ];
export const uint64 = [ 0, uintMax(64) ];

// Input types
export const numericInputTypes = [ 'int', 'long', 'uint64', 'double', 'decimal' ];
export const inputTypes = [ 'string', 'objectid', 'bool', 'date', ...numericInputTypes ];

export function isBsonBuiltin(value) {
  return (
    (typeof value === 'object') &&
    (value !== null) &&
    (typeof value._bsontype === 'string') &&
    (typeof value.inspect === 'function')
  );
}

export function canBeObjectId(value) {
  try {
    new ObjectId(value);
    return true;
  }
  catch {
    return false;
  }
}
