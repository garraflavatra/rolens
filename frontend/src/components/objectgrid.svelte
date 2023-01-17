<script>
  import Grid from './grid.svelte';

  export let data = [];
  export let key = '_id';
  export let showHeaders = false;
  export let activeKey = '';

  const columns = [
    { key: 'key', label: 'Key' },
    { key: 'value', label: 'Value' },
    { key: 'type', label: 'Type' },
  ];

  let items = [];

  $: if (data) {
    items = [];

    if (Array.isArray(data)) {
      for (const item of data) {
        const _item = {};
        _item.key = item[key];
        _item.children = dissectObject(item);
        items = [ ...items, _item ];
      }
    }
    else {
      items = dissectObject(data);
    }
  }

  function getType(value) {
    if (Array.isArray(value)) {
      return `array (${value.length} item${value.length === 1 ? '' : 's'})`;
    }
    else if (typeof value === 'number') {
      if (value.toString().includes('.')) {
        return 'double';
      }
      return 'integer';
    }
    else if (new Date(value).toString() !== 'Invalid Date') {
      return 'date';
    }
    else if (typeof value === 'object') {
      const keys = Object.keys(value);
      return `object (${keys.length} item${keys.length === 1 ? '' : 's'})`;
    }
    else {
      return typeof value;
    }
  }

  function dissectObject(object) {
    const entries = [ ...Array.isArray(object) ? object.entries() : Object.entries(object) ];
    return entries.map(([ key, value ]) => {
      key = key + '';
      const type = getType(value);
      const child = { key, value, type };

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
  on:selectChild
  on:trigger
  bind:activeKey
  {columns}
  {items}
  {showHeaders}
/>
