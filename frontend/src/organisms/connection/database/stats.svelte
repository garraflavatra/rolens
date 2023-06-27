<script>
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/objectgrid.svelte';

  export let database;

  let copySucceeded = false;

  async function copy() {
    const json = JSON.stringify(database.stats, undefined, '\t');
    await navigator.clipboard.writeText(json);
    copySucceeded = true;
    setTimeout(() => copySucceeded = false, 1500);
  }
</script>

<div class="stats">
  <!-- <CodeExample code="db.stats()" /> -->

  <div class="grid">
    <ObjectGrid
      data={database.stats}
      errorTitle={database.statsError ? 'Error fetching database stats' : ''}
      errorDescription={database.statsError}
      busy={!database.stats && !database.statsError && `Fetching stats for ${database.key}…`}
    />
  </div>

  <div class="buttons">
    <button class="button secondary" on:click={copy} disabled={!database.stats}>
      <Icon name={copySucceeded ? 'check' : 'clipboard'} />
      Copy JSON
    </button>
  </div>
</div>

<style>
  .stats {
    display: grid;
    gap: 0.5rem;
    grid-template: 1fr auto / 1fr;
  }

  .stats .grid {
    overflow: auto;
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
  }
</style>
