<script>
  import { createEventDispatcher } from 'svelte';
  import Icon from './icon.svelte';

  export let columns = [];
  export let items = [];
  export let key = 'id';
  export let activeKey = '';
  export let activeChildKey = '';
  export let showHeaders = true;
  export let level = 0;
  export let contained = false;

  const dispatch = createEventDispatcher();
  const childrenOpen = {};

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

  function toggleChildren(itemKey) {
    childrenOpen[itemKey] = !childrenOpen[itemKey];
  }
</script>

<div class:grid={level === 0} class:subgrid={level > 0} class:contained>
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
        <tr on:click={() => select(item[key])} class:selected={activeKey === item[key] && !activeChildKey}>
          <td class="has-toggle">
            {#if item.children}
              <button class="toggle" on:click={() => toggleChildren(item[key])}>
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
    background-color: #fff;
    height: 100%;
    width: 100%;
    overflow: scroll;
  }
  .grid.contained {
    border: 1px solid #ccc;
  }
  .subgrid {
    width: 100%;
  }

  table {
    border-collapse: collapse;
    width: 100%;
  }

  table thead {
    border-bottom: 2px solid #ccc;
  }

  table th {
    font-weight: 600;
    text-align: left;
  }

  tr {
    cursor: pointer;
  }

  td, th {
    padding: 0.3rem;
    height: 100%;
  }
  td.has-toggle {
    width: calc(20px + 0.3rem);
  }
  td.subgrid-parent {
    padding: 0;
  }

  table tbody tr.selected td {
    background-color: #00008b;
    color: #fff;
  }

  button.toggle {
    color: inherit;
    padding: 0;
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
