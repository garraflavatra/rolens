<script>
  import { createEventDispatcher } from 'svelte';
  import GridItems from './grid-items.svelte';
  import Icon from './icon.svelte';

  export let columns = [];
  export let items = [];
  export let key = 'id';
  export let activePath = [];
  export let striped = true;
  export let showHeaders = false;
  export let hideObjectIndicators = false;
  export let hideChildrenToggles = false;
  export let canAddRows = false;
  export let canSelect = true;
  export let canRemoveItems = false;
  export let inputsValid = false;
  // export let actions = [];

  const dispatch = createEventDispatcher();

  function addRow() {
    dispatch('addRow');
  }
</script>

<div class="grid">
  <!-- {#if actions?.length}
    <div class="actions">
      {#each actions as action}
        <button class="btn" on:click={action.fn} disabled={action.disabled}>
          {#if action.icon}<Icon name={action.icon} />{/if}
          {action.label || ''}
        </button>
      {/each}
    </div>
  {/if} -->

  <table>
    {#if showHeaders && columns.some(col => col.title)}
      <thead>
        <tr>
          {#if !hideChildrenToggles}
            <th class="has-toggle"></th>
          {/if}

          <th class="has-icon"></th>

          {#each columns as column}
            <th scope="col">{column.title || ''}</th>
          {/each}

          {#if canRemoveItems}
            <th class="has-button"></th>
          {/if}
        </tr>
      </thead>
    {/if}

    <tbody>
      <GridItems
        {items}
        {columns}
        {key}
        {striped}
        {canSelect}
        {canRemoveItems}
        {hideObjectIndicators}
        {hideChildrenToggles}
        bind:activePath
        bind:inputsValid
        on:select
        on:trigger
      />
    </tbody>

    {#if canAddRows}
      <tfoot>
        <button class="btn-sm" type="button" on:click={addRow}>
          <Icon name="+" />
        </button>
      </tfoot>
    {/if}
  </table>
</div>

<style>
  .grid {
    width: 100%;
    height: 100%;
    background-color: #fff;
  }

  /* .actions {
    margin-bottom: 0.5rem;
    padding: 0.5rem;
    border-bottom: 1px solid #ccc;
  }
  .actions button {
    margin-right: 0.2rem;
  } */

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
    padding: 2px;
  }

  /* tfoot button {
    margin-top: 0.5rem;
  } */
</style>
