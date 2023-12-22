<script>
  import { pathsAreEqual, resolveKeypath, setValue } from '$lib/objects.js';
  import contextMenu from '$lib/stores/contextmenu.js';
  import { createEventDispatcher } from 'svelte';
  import FormInput from '$components/editors/forminput.svelte';
  import Icon from '$components/icon.svelte';

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
      if (item.children) {
        item.children = objectToArray(item.children);
      }
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
      return Object.entries(obj).map(([
        k,
        item,
      ]) => {
        return { ...item, [key]: k };
      });
    }
    else {
      return obj;
    }
  }

  function select(itemKey, index) {
    if (!canSelect) {
      return false;
    }

    activeKey = itemKey;
    activePath = [
      ...path.slice(0, level),
      itemKey,
    ];
    dispatch('select', { level, itemKey, index, path: activePath });
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
    toggleChildren(itemKey, false);
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
  {@const selected = canSelect && pathsAreEqual(activePath, [
    ...path,
    item[key],
  ])}

  <tr
    on:click={() => select(item[key], index)}
    on:dblclick={() => doubleClick(item[key], index)}
    on:contextmenu|preventDefault={evt => showContextMenu(evt, item)}
    class:selectable={canSelect}
    class:selected
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

        {#if item.loading}
          <span class="spinner" style:margin-left="{level * 10}px" />
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
          <FormInput
            {column}
            bind:value={keypathProxies[index][column.key]}
            bind:valid={validity[columnIndex]}
          />
        {:else}
          <div class="value" style:margin-left="{level * 10}px">
            {formatValue(keypathProxies[index][column.key])}
          </div>
        {/if}
      </td>
    {/each}

    {#if canRemoveItems}
      <td class="has-button">
        <button
          class="button-small"
          type="button"
          on:click|stopPropagation={() => removeItem(index, item[key])}
          on:dblclick|stopPropagation
        >
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
      path={[
        ...path,
        item[key],
      ]}
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
    background-color: var(--selection) !important;
  }

  td {
    padding: 4px 2px;
    text-overflow: ellipsis;
  }
  td.has-toggle {
    position: relative;
    width: 1.5em;
  }
  td.has-toggle .spinner {
    position: absolute;
    top: 0.3em;
    left: 0.22em;
    width: 1em;
    height: 1em;
    border: 1px solid var(--primary);
    border-radius: 50px;
    border-top-color: transparent;
    border-left-color: transparent;
    animation: .6s linear 0 spin;
    animation-iteration-count: infinite;
  }
  td.has-icon {
    padding: 0;
    width: 1.5em;
  }
  td.has-icon :global(svg) {
    width: 1em;
    height: 1em;
  }

  td .value {
    height: 1.2em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 25em;
  }

  button.toggle {
    margin: 2px 0 0 3px;
    padding: 0;
    color: inherit;
  }
  button.toggle :global(svg) {
    width: 0.9em;
    height: 0.9em;
    vertical-align: top;
  }

  .right {
    text-align: right;
  }
</style>
