<script>
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/objectgrid.svelte';

  export let host;

  let copySucceeded = false;

  async function copy() {
    const json = JSON.stringify(host.systemInfo, undefined, '\t');
    await navigator.clipboard.writeText(json);
    copySucceeded = true;
    setTimeout(() => copySucceeded = false, 1500);
  }
</script>

<div class="stats">
  <!-- <CodeExample code="db.stats()" /> -->

  <div class="grid">
    <ObjectGrid
      data={host.systemInfo}
      errorTitle={host.systemInfoError ? 'Error fetching system info' : ''}
      errorDescription={host.systemInfoError}
      busy={!host.systemInfo && !host.systemInfoError && 'Fetching system info…'}
    />
  </div>

  <div class="buttons">
    <button class="btn secondary" on:click={copy} disabled={!host.systemInfo}>
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
