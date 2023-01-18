<script>
  import GridItems from './grid-items.svelte';
  import Icon from './icon.svelte';

  export let columns = [];
  export let items = [];
  export let actions = [];
  export let key = 'id';
  export let activePath = [];
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
    <tbody>
      <GridItems {items} {columns} {key} {striped} bind:activePath on:select on:trigger />
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
</style>
