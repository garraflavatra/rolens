<script>
  import { resolveKeypath, setValue } from '$lib/objects';
  import contextMenu from '$lib/stores/contextmenu';
  import { createEventDispatcher } from 'svelte';
  import FormInput from './forminput.svelte';
  import Icon from './icon.svelte';

  export let items = [];
  export let columns = [];
  export let key = '';
  export let path = [];
  export let activeKey = '';
  export let activePath = [];
  export let level = 0;
  export let striped = true;
  export let hideObjectIndicators = false;
  export let hideChildrenToggles = false;
  export let canSelect = true;
  export let canRemoveItems = false;
  export let inputsValid = false;

  const dispatch = createEventDispatcher();
  const keypathProxies = {};
  const validity = {};
  let childrenOpen = {};
  let _items = [];

  $: refresh(hideObjectIndicators, items);
  $: inputsValid = Object.values(validity).every(v => v !== false);

  function refresh(hideObjectIndicators, items) {
    _items = objectToArray(items).map(item => {
      item.children = objectToArray(item.children);
      return item;
    });

    for (let index = 0; index < _items.length; index++) {
      keypathProxies[index] = new Proxy(_items, {
        get: (_items, key) => resolveKeypath(_items[index], key),
        set: (_items, key, value) => {
          setValue(_items[index], key, value);
          return true;
        },
      });
    }
  }

  function objectToArray(obj) {
    if (Array.isArray(obj)) {
      return obj;
    }
    else if ((typeof obj === 'object') && (obj !== null)) {
      return Object.entries(obj).map(([ k, item ]) => ({ ...item, [key]: k }));
    }
    else {
      return obj;
    }
  }

  function select(itemKey, index) {
    if (!canSelect) {
      return false;
    }

    toggleChildren(itemKey, false);

    if (activeKey !== itemKey) {
      activeKey = itemKey;
      if (level === 0) {
        activePath = [ itemKey ];
      }
      else {
        activePath = [ ...path, itemKey ];
      }
      dispatch('select', { level, itemKey, index });
    }
  }

  function closeAll() {
    childrenOpen = {};
    dispatch('closeAll');
  }

  function toggleChildren(itemKey, shift = false) {
    childrenOpen[itemKey] = !childrenOpen[itemKey];
    if (shift) {
      closeAll();
    }
  }

  function doubleClick(itemKey, index) {
    // toggleChildren(itemKey, false);
    dispatch('trigger', { level, itemKey, index });
    childrenOpen[itemKey] = true;
  }

  function showContextMenu(evt, item) {
    select(item[key]);
    contextMenu.show(evt, item.menu);
  }

  function removeItem(index, itemKey) {
    if (Array.isArray(items)) {
      items.splice(index, 1);
      items = items;
    }
    dispatch('removeItem', itemKey);
  }

  function formatValue(value) {
    if (Array.isArray(value)) {
      return hideObjectIndicators ? '' : '[...]';
    }
    if (typeof value === 'number' || typeof value === 'boolean') {
      return String(value);
    }
    if ((value === undefined) || (value === null)) {
      return '';
    }
    // if (new Date(value).toString() !== 'Invalid Date') {
    //   return new Date(value);
    // }
    if ((typeof value === 'object') && (value !== null)) {
      return hideObjectIndicators ? '' : '{...}';
    }
    if (String(value).length <= 1000) {
      return value;
    }
    return String(value).slice(0, 1000) + '…';
  }
</script>

{#each _items as item, index}
  <tr
    on:click={() => select(item[key], index)}
    on:dblclick={() => doubleClick(item[key], index)}
    on:contextmenu|preventDefault={evt => showContextMenu(evt, item)}
    class:selectable={canSelect}
    class:selected={canSelect && !activePath[level + 1] && activePath.every(k => path.includes(k) || k === item[key]) && (activePath[level] === item[key])}
    class:striped
  >
    {#if !hideChildrenToggles}
      <td class="has-toggle">
        {#if item.children?.length}
          <button
            class="toggle"
            on:click|stopPropagation={evt => toggleChildren(item[key], evt.shiftKey)}
            style:transform="translateX({level * 10}px)"
          >
            <Icon name={childrenOpen[item[key]] ? 'chev-d' : 'chev-r'} />
          </button>
        {/if}
      </td>
    {/if}

    <td class="has-icon">
      <div style:margin-left="{level * 10}px">
        <Icon name={item.icon} />
      </div>
    </td>

    {#each columns as column, columnIndex}
      <td class:right={column.right} title={keypathProxies[index][column.key]}>
        {#if column.inputType}
          <FormInput {column} bind:value={keypathProxies[index][column.key]} bind:valid={validity[columnIndex]} />
        {:else}
          <div class="value" style:margin-left="{level * 10}px">
            {formatValue(keypathProxies[index][column.key])}
          </div>
        {/if}
      </td>
    {/each}

    {#if canRemoveItems}
      <td class="has-button">
        <button class="button-small" type="button" on:click|stopPropagation={() => removeItem(index, item[key])} on:dblclick|stopPropagation>
          <Icon name="x" />
        </button>
      </td>
    {/if}
  </tr>

  {#if item.children && childrenOpen[item[key]]}
    <svelte:self
      {columns}
      {key}
      {striped}
      {hideObjectIndicators}
      {hideChildrenToggles}
      {canSelect}
      {canRemoveItems}
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
  tr.selectable {
    cursor: pointer;
  }
  tr.selectable.selected td {
    background-color: #00008b !important;
    color: #fff;
  }

  td {
    padding: 2px;
    text-overflow: ellipsis;
  }
  td.has-toggle {
    width: 20px;
  }
  td.has-icon {
    padding: 0;
    width: 17px;
  }
  td.has-icon :global(svg) {
    width: 13px;
    height: 13px;
  }

  td .value {
    height: 15px;
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
    width: 13px;
    height: 13px;
    vertical-align: top;
  }

  .right {
    text-align: right;
  }
</style>
