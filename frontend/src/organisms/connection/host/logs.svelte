<script>
  import Grid from '$components/grid.svelte';
  import Icon from '$components/icon.svelte';
  import ObjectViewer from '$components/objectviewer.svelte';
  import input from '$lib/actions/input';
  import { BrowserOpenURL } from '$wails/runtime/runtime';
  import { onDestroy } from 'svelte';

  export let host;

  const autoReloadIntervals = [ 1, 2, 5, 10, 30, 60 ];
  let filter = 'global';
  let logs;
  let total = 0;
  let error = '';
  let copySucceeded = false;
  let autoReloadInterval = 0;
  let objectViewerData;
  let interval;
  $: filter && refresh();
  $: busy = !logs && !error && 'Requesting logs…';

  $: if (autoReloadInterval) {
    if (interval) {
      clearInterval(interval);
    }
    interval = setInterval(refresh, autoReloadInterval * 1000);
  }

  async function refresh() {
    let _logs = [];
    ({ logs: _logs, total, error } = await host.getLogs(filter));

    logs = [];
    for (let index = 0; index < _logs.length; index++) {
      const log = JSON.parse(_logs[index]);
      log._index = index;
      logs = [ ...logs, log ];
    }
  }

  function openFilterDocs() {
    BrowserOpenURL('https://www.mongodb.com/docs/manual/reference/command/getLog/#command-fields');
  }

  function openLogDetail(event) {
    objectViewerData = logs[event.detail.index];
  }

  async function copy() {
    const json = JSON.stringify(host.status, undefined, '\t');
    await navigator.clipboard.writeText(json);
    copySucceeded = true;
    setTimeout(() => copySucceeded = false, 1500);
  }

  onDestroy(() => {
    if (interval) {
      clearInterval(interval);
    }
  });
</script>

<div class="stats">
  <div class="grid">
    <Grid
      items={logs || []}
      columns={[
        { title: 'Date', key: 't.$date' },
        { title: 'Severity', key: 's' },
        { title: 'ID', key: 'id' },
        { title: 'Component', key: 'c' },
        { title: 'Context', key: 'ctx' },
        { title: 'Message', key: 'msg' },
      ]}
      key="_index"
      showHeaders
      errorTitle={error ? 'Error fetching server status' : ''}
      errorDescription={error}
      on:trigger={openLogDetail}
      {busy}
    />
  </div>

  <div class="controls">
    <div>
      <div class="field inline">
        <button class="button" on:click={refresh}>
          <Icon name="reload" spin={busy} /> Reload
        </button>

        <button class="button secondary" on:click={copy} disabled={!host.status}>
          <Icon name={copySucceeded ? 'check' : 'clipboard'} />
          Copy JSON
        </button>
      </div>

      <label class="field inline">
        <span class="label">Reload (sec)</span>
        <input type="number" class="autoreloadinput" bind:value={autoReloadInterval} list="autoreloadintervals" use:input />
      </label>

      <label class="field inline">
        <select bind:value={filter}>
          <option value="global">Global</option>
          <option value="startupWarnings">Startup warnings</option>
        </select>
        <button class="button secondary" on:click={openFilterDocs} title="Documentation">
          <Icon name="?" />
        </button>
      </label>
    </div>

    {#if total}
      <div class="total">
        Total: {total}
      </div>
    {/if}
  </div>
</div>

{#if objectViewerData}
  <ObjectViewer bind:data={objectViewerData} readonly />
{/if}

<datalist id="autoreloadintervals">
  {#each autoReloadIntervals as value}
    <option {value} />
  {/each}
</datalist>

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

  .controls {
    display: flex;
    align-items: center;
    gap: 0.1rem;
  }
  .total {
    margin-left: auto;
  }
  .autoreloadinput {
    width: 1.5rem;
  }
</style>
