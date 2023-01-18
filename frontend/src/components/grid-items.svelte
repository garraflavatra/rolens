<script>
  import { contextMenu } from '../stores';
  import { createEventDispatcher } from 'svelte';
  import Icon from './icon.svelte';

  export let items = [];
  export let columns = [];
  export let key = '';
  export let path = [];
  export let activeKey = '';
  export let activePath = [];
  export let level = 0;
  export let striped = true;

  const dispatch = createEventDispatcher();
  let childrenOpen = {};

  $: _items = objectToArray(items).map(item => {
    item.children = objectToArray(item.children);
    return item;
  });

  function objectToArray(obj) {
    if (Array.isArray(obj)) {
      return obj;
    }
    else if (typeof obj === 'object') {
      return Object.entries(obj).map(([ k, item ]) => ({ ...item, [key]: k }));
    }
    else {
      return obj;
    }
  }

  function select(itemKey) {
    activeKey = itemKey;
    if (level === 0) {
      activePath = [ itemKey ];
    }
    else {
      activePath = [ ...path, itemKey ];
    }
    dispatch('select', { level, itemKey });
  }

  function closeAll() {
    childrenOpen = {};
    dispatch('closeAll');
  }

  function toggleChildren(itemKey, shift) {
    childrenOpen[itemKey] = !childrenOpen[itemKey];
    if (shift) {
      closeAll();
    }
  }

  function doubleClick(itemKey) {
    toggleChildren(itemKey, false);
    dispatch('trigger', { level, itemKey });
  }

  function showContextMenu(evt, item) {
    select(item[key]);
    contextMenu.show(evt, item.menu);
  }

  function formatValue(value) {
    if (Array.isArray(value)) {
      return '[...]';
    }
    if (typeof value === 'object') {
      return '{...}';
    }
    if (value === undefined || value === null) {
      return '';
    }
    if (typeof value === 'number' || typeof value === 'boolean') {
      return String(value);
    }
    if (String(value).length <= 1000) {
      return value;
    }
    return String(value).slice(0, 1000) + '…';
  }
</script>

{#each _items as item}
  <tr
    on:click={() => select(item[key])}
    on:dblclick={() => doubleClick(item[key])}
    on:contextmenu|preventDefault={evt => showContextMenu(evt, item)}
    class:selected={!activePath[level + 1] && activePath.every(k => path.includes(k) || k === item[key]) && (activePath[level] === item[key])}
    class:striped
  >
    <td class="has-toggle">
      {#if item.children?.length}
        <button
          class="toggle"
          on:click={evt => toggleChildren(item[key], evt.shiftKey)}
          style:transform="translateX({level * 10}px)"
        >
          <Icon name={childrenOpen[item[key]] ? 'chev-d' : 'chev-r'} />
        </button>
      {/if}
    </td>

    {#each columns as column, columnIndex}
      {@const value = item[column.key]}
      <td class:right={column.right} title={value}>
        <div class="value" style:margin-left="{level * 10}px">
          {formatValue(value)}
        </div>
      </td>
    {/each}
  </tr>

  {#if item.children && childrenOpen[item[key]]}
    <svelte:self
      {columns}
      {key}
      {striped}
      path={[ ...path, item[key] ]}
      items={item.children}
      level={level + 1}
      bind:activePath
      on:closeAll={closeAll}
      on:select
      on:trigger
    />
  {/if}
{/each}

<style>
  tr.striped:nth-of-type(even) td {
    background-color: #eee;
  }
  tr.selected td {
    background-color: #00008b !important;
    color: #fff;
  }

  td {
    padding: 0.3rem;
    text-overflow: ellipsis;
    cursor: pointer;
  }
  td.has-toggle {
    width: calc(20px + 0.3rem);
  }

  td .value {
    height: 2.1ex;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 25em;
  }

  button.toggle {
    color: inherit;
    padding: 0;
    margin: 0;
    vertical-align: top;
  }
  button.toggle :global(svg) {
    width: 15px;
    height: 15px;
    vertical-align: top;
  }

  .right {
    text-align: right;
  }
</style>
