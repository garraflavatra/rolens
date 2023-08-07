<script>
  import { isBsonBuiltin } from '$lib/mongo/index.js';
  import { isDate } from 'date-fns';
  import Grid from './grid.svelte';

  export let data = [];
  export let key = '_id';
  export let activePath = [];
  export let hideObjectIndicators = false;
  export let getRootMenu = () => undefined;
  export let errorTitle = '';
  export let errorDescription = '';
  export let busy = false;
  export let showTypes = true;

  let items = [];

  $: columns = [
    { key: 'key', label: 'Key' },
    { key: 'value', label: 'Value' },
    showTypes ? { key: 'type', label: 'Type' } : undefined,
  ].filter(c => !!c);

  $: if (data) {
    items = [];

    if (Array.isArray(data)) {
      for (const item of data) {
        const newItem = {};
        newItem.key = stringifyValue(item[key]);
        newItem.type = getType(item[key]);
        newItem.children = dissectObject(item);
        newItem.menu = getRootMenu(key, item[key]);
        items = [ ...items, newItem ];
      }
    }
    else {
      items = dissectObject(data);
    }
  }

  function getType(value) {
    if (isBsonBuiltin(value)) {
      return value._bsontype;
    }
    else if (isDate(value)) {
      return 'Date';
    }
    else if (Array.isArray(value)) {
      return `array (${value.length} item${value.length === 1 ? '' : 's'})`;
    }
    else if (typeof value === 'number') {
      if (value.toString().includes('.')) {
        return 'double';
      }
      return 'integer';
    }
    else if (value === null) {
      return 'null';
    }
    else if (typeof value === 'object') {
      const keys = Object.keys(value);
      return `object (${keys.length} item${keys.length === 1 ? '' : 's'})`;
    }
    else {
      return typeof value;
    }
  }

  function stringifyValue(value) {
    if (isBsonBuiltin(value)) {
      value = value.inspect?.();
      if (value.startsWith('new ')) {
        value = value.slice(4);
      }

      if (value.startsWith('Int32(')) {
        value = value.slice(6, -1);
      }
      else if (value.startsWith('Double(')) {
        value = value.slice(7, -1);
      }
      else if (value.startsWith('Binary(Buffer.from(')) {
        value = `BinData(${JSON.stringify(value.sub_type || 0)}, ${value.slice(19, -1)}`;
      }
    }
    else if (isDate(value)) {
      value = value.toString();
    }
    return value;
  }

  function dissectObject(object) {
    const entries = [ ...Array.isArray(object) ? object.entries() : Object.entries(object) ];
    return entries.map(([ key, value ]) => {
      const type = getType(value);
      key = key + '';
      const child = {
        key,
        type,
        value: stringifyValue(value),
        menu: value?.menu,
      };

      if (type.startsWith('object') || type.startsWith('array')) {
        child.children = dissectObject(value);
      }

      return child;
    });
  }
</script>

<Grid
  key="key"
  on:select
  on:trigger
  bind:activePath
  {columns}
  {items}
  {hideObjectIndicators}
  {errorTitle}
  {errorDescription}
  {busy}
/>
