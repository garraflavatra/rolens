<script>
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/grid/objectgrid.svelte';

  export let host;

  let copySucceeded = false;

  async function copy() {
    const json = JSON.stringify(host.status, undefined, '\t');
    await navigator.clipboard.writeText(json);
    copySucceeded = true;
    setTimeout(() => copySucceeded = false, 1500);
  }
</script>

<div class="stats">
  <!-- <CodeExample code="db.stats()" /> -->

  <div class="grid">
    <ObjectGrid
      data={host.status || {}}
      showTypes={false}
      errorTitle={host.statusError ? 'Error fetching server status' : ''}
      errorDescription={host.statusError}
      busy={!host.status && !host.statusError && 'Fetching server status…'}
    />
  </div>

  <div class="buttons">
    <button class="button secondary" on:click={copy} disabled={!host.status}>
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
