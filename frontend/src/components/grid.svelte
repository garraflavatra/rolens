<script>
  import GridItems from './grid-items.svelte';
  import Icon from './icon.svelte';

  export let columns = [];
  export let items = [];
  export let actions = [];
  export let key = 'id';
  export let activeKey = '';
  export let activeChildKey = '';
  export let showHeaders = true;
  export let striped = true;
</script>

<div class="grid">
  {#if actions?.length}
    <div class="actions">
      {#each actions as action}
        <button class="btn" on:click={action.fn} disabled={action.disabled}>
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
      <GridItems {items} {columns} {key} {striped} bind:activeKey bind:activeChildKey on:select on:selectChild on:trigger />
    </tbody>
  </table>
</div>

<style>
  .grid {
    width: 100%;
    height: 100%;
    background-color: #fff;
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
    background-color: #fff;
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

  th {
    padding: 0.3rem;
    text-overflow: ellipsis;
  }
</style>
