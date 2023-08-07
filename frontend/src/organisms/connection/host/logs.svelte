<script>
  import Grid from '$components/grid/grid.svelte';
  import Icon from '$components/icon.svelte';
  import ObjectViewer from '$components/objectviewer.svelte';
  import input from '$lib/actions/input.js';
  import { logComponents, logLevels } from '$lib/mongo/index.js';
  import { BrowserOpenURL } from '$wails/runtime/runtime.js';
  import { onDestroy } from 'svelte';

  export let host;
  export let visible = false;

  const autoReloadIntervals = [ 1, 2, 5, 10, 30, 60 ];
  let filter = 'global';
  let severityFilter = '';
  let componentFilter = '';
  let logs;
  let total = 0;
  let error = '';
  let copySucceeded = false;
  let autoReloadInterval = 0;
  let objectViewerData;
  let interval;
  $: (filter || severityFilter || componentFilter) && refresh();
  $: busy = !logs && !error && 'Requesting logs…';

  $: if (autoReloadInterval) {
    if (interval) {
      clearInterval(interval);
    }
    interval = setInterval(refresh, autoReloadInterval * 1000);
  }

  async function refresh() {
    if (!visible) {
      return;
    }

    let _logs = [];
    ({ logs: _logs, total, error } = await host.getLogs(filter));
    logs = [];

    for (let index = 0; index < _logs.length; index++) {
      const log = JSON.parse(_logs[index]);

      const matchesLevel = severityFilter ? log.s?.startsWith(severityFilter) : true;
      const matchesComponent = componentFilter ? (componentFilter === log.c?.toUpperCase()) : true;

      if (matchesLevel && matchesComponent) {
        log._index = index;
        log.s = logLevels[log.s] || log.s;
        logs = [ ...logs, log ];
      }
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

  $: visible && !logs && refresh();

  onDestroy(() => {
    if (interval) {
      clearInterval(interval);
    }
  });
</script>

<div class="stats">
  <div class="formrow">
    <label class="field">
      <span class="label">Auto reload (seconds)</span>
      <input
        type="number"
        class="autoreloadinput"
        bind:value={autoReloadInterval}
        list="autoreloadintervals"
        use:input
      />
    </label>

    <label class="field">
      <span class="label">Log type</span>
      <select bind:value={filter}>
        <option value="global">Global</option>
        <option value="startupWarnings">Startup warnings</option>
      </select>
      <button class="button secondary" on:click={openFilterDocs} title="Documentation">
        <Icon name="?" />
      </button>
    </label>
  </div>

  <div class="formrow">
    <label class="field">
      <span class="label">Severity</span>
      <select bind:value={severityFilter}>
        <option value="">All</option>
        {#each Object.entries(logLevels) as [ value, name ]}
          <option {value}>{value} ({name})</option>
        {/each}
      </select>
    </label>

    <label class="field">
      <span class="label">Component</span>
      <select bind:value={componentFilter}>
        <option value="">All</option>
        {#each logComponents as value}
          <option {value}>{value}</option>
        {/each}
      </select>
    </label>
  </div>

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
    <button class="button" on:click={refresh}>
      <Icon name="reload" spin={busy} /> Reload
    </button>

    <button class="button secondary" on:click={copy} disabled={!host.status}>
      <Icon name={copySucceeded ? 'check' : 'clipboard'} />
      Copy JSON
    </button>

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
    grid-template: auto auto 1fr auto / 1fr;
  }

  .formrow {
    display: grid;
    gap: 0.5rem;
    grid-template: 1fr / 1fr 1fr;
  }

  .grid {
    overflow: auto;
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
  }

  .controls {
    display: flex;
    align-items: center;
    gap: 0.2rem;
  }
  .total {
    margin-left: auto;
  }
  .autoreloadinput {
    width: 1.5rem;
  }
</style>
