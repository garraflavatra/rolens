<script>
  import { onDestroy } from 'svelte';
  import BlankState from './blankstate.svelte';
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
  export let canSelect = true;
  export let canRemoveItems = false;
  export let inputsValid = false;
  export let errorTitle = '';
  export let errorDescription = '';
  export let busy = false;

  let copySucceeded = false;
  let timeout;

  async function copyErrorDescription() {
    await navigator.clipboard.writeText(errorDescription);
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  onDestroy(() => clearTimeout(timeout));
</script>

<div class="grid">
  <!-- {#if actions?.length}
    <div class="actions">
      {#each actions as action}
        <button class="button" on:click={action.fn} disabled={action.disabled}>
          {#if action.icon}<Icon name={action.icon} />{/if}
          {action.label || ''}
        </button>
      {/each}
    </div>
  {/if} -->

  {#if busy}
    <BlankState label={(busy === true) ? 'Loading…' : busy} icon="loading" />
  {:else if errorTitle || errorDescription}
    <BlankState title={errorTitle} label={errorDescription} icon="!">
      <button class="button-small" on:click={copyErrorDescription}>
        <Icon name={copySucceeded ? 'check' : 'clipboard'} /> Copy error message
      </button>
    </BlankState>
  {:else}
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
          on:removeItem
        />
      </tbody>
    </table>
  {/if}
</div>

<style>
  .grid {
    width: 100%;
    height: 100%;
    background-color: #fff;
  }

  table {
    border-collapse: collapse;
    width: 100%;
    background-color: #fff;
  }

  thead th {
    font-weight: 600;
    text-align: left;
    padding: 2px;
    /* border-bottom: 2px solid #ccc; */
    box-shadow: 0 2px #ccc;
    background-color: #fff;
    position: sticky;
    top: 0;
  }

  .grid :global(.blankstate) {
    height: 100%;
    padding: 1rem;
  }

  /* tfoot button {
    margin-top: 0.5rem;
  }
  .actions {
    margin-bottom: 0.5rem;
    padding: 0.5rem;
    border-bottom: 1px solid #ccc;
  }
  .actions button {
    margin-right: 0.2rem;
  } */
</style>
