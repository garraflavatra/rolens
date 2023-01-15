<script>
  import { contextMenu } from '../stores';
  import { createEventDispatcher } from 'svelte';
  import Icon from './icon.svelte';

  export let columns = [];
  export let items = [];
  export let actions = [];
  export let key = 'id';
  export let activeKey = '';
  export let activeChildKey = '';
  export let showHeaders = true;
  export let level = 0;
  export let contained = false;

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
    activeChildKey = '';
    dispatch('select', itemKey);
  }

  function selectChild(itemKey, childKey) {
    select(itemKey);
    activeChildKey = childKey;
    dispatch('selectChild', childKey);
  }

  function toggleChildren(itemKey, closeAll) {
    childrenOpen[itemKey] = !childrenOpen[itemKey];
    if (closeAll) {
      childrenOpen = {};
      dispatch('closeAll');
    }
  }

  function doubleClick(itemKey) {
    toggleChildren(itemKey, false);
    dispatch('trigger', itemKey);
  }

  function showContextMenu(evt, item) {
    select(item[key]);
    contextMenu.show(evt, item.menu);
  }
</script>

<div class:grid={level === 0} class:subgrid={level > 0} class:contained>
  {#if actions?.length}
    <div class="actions">
      {#each actions as action}
        <button class="btn" on:click={action.fn}>
          {#if action.icon}<Icon name={action.icon} />{/if}
          {action.label || ''}
        </button>
      {/each}
    </div>
  {/if}

  <table>
    {#if showHeaders && columns.some(col => col.title)}
      <thead>
        <tr>
          <th class="has-toggle"></th>
          {#each columns as column}
            <th scope="col">{column.title || ''}</th>
          {/each}
        </tr>
      </thead>
    {/if}

    <tbody>
      {#each _items as item (item[key])}
        <tr
          on:click={() => select(item[key])}
          on:dblclick={() => doubleClick(item[key])}
          on:contextmenu|preventDefault={evt => showContextMenu(evt, item)}
          class:selected={activeKey === item[key] && !activeChildKey}
        >
          <td class="has-toggle">
            {#if item.children?.length}
              <button class="toggle" on:click={evt => toggleChildren(item[key], evt.shiftKey)}>
                <Icon name={childrenOpen[item[key]] ? 'chev-d' : 'chev-r'} />
              </button>
            {/if}
          </td>

          {#each columns as column}
            {@const value = item[column.key]}
            <td class:right={column.right}>
              {#if typeof value !== 'object'}
                {value || ''}
              {/if}
            </td>
          {/each}
        </tr>

        {#if item.children && childrenOpen[item[key]]}
          <tr>
            <td></td>
            <td colspan={columns.length + 1} class="subgrid-parent">
              <svelte:self
                {columns}
                {key}
                bind:activeKey={activeChildKey}
                showHeaders={false}
                items={item.children}
                level={level + 1}
                on:select={e => selectChild(item[key], e.detail)}
                on:closeAll={() => (childrenOpen = {})}
              />
            </td>
          </tr>
        {/if}
      {/each}
    </tbody>
  </table>
</div>

<style>
  .grid {
    width: 100%;
    height: 100%;
    background-color: #fff;
  }
  .grid.contained {
    border: 1px solid #ccc;
  }
  .subgrid {
    width: 100%;
  }

  .actions {
    margin-bottom: 0.5rem;
    padding: 0.5rem;
    border-bottom: 1px solid #ccc;
  }
  .actions button {
    margin-right: 0.2rem;
  }

  table {
    border-collapse: collapse;
    width: 100%;
  }

  table thead {
    border-bottom: 2px solid #ccc;
  }

  th {
    font-weight: 600;
    text-align: left;
  }

  tr {
    cursor: pointer;
  }

  td, th {
    padding: 0.3rem;
  }
  td.has-toggle {
    width: calc(20px + 0.3rem);
  }
  td.subgrid-parent {
    padding: 0;
  }

  tbody tr.selected td {
    background-color: #00008b;
    color: #fff;
  }

  button.toggle {
    color: inherit;
    padding: 0;
    margin: 0;
    display: contents;
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
