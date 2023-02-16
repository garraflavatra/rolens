<script>
  import DirectoryChooser from '$components/directorychooser.svelte';
  import Grid from '$components/grid.svelte';
  import Modal from '$components/modal.svelte';
  import busy from '$lib/stores/busy';
  import { connections } from '$lib/stores/connections';
  import applicationSettings from '$lib/stores/settings';
  import { OpenConnection, OpenDatabase, PerformExport } from '$wails/go/app/App';

  export let info;
  export let hosts = {};

  const actionLabel = {
    export: 'Perform export',
    dump: 'Perform dump',
  };

  $: if (info) {
    info.outdir = info.outdir || $applicationSettings.defaultExportDirectory;
    info.filename = info.filename || `Export ${new Date().getTime()}`;

    if (info.filetype === 'bson') {
      info.type = 'dump';
    }
    else {
      info.type = 'export';
    }
  }

  async function selectHost(hostKey) {
    info.hostKey = hostKey;
    info.dbKey = undefined;
    info.collKeys = [];

    if (hostKey) {
      busy.start();
      const databases = await OpenConnection(hostKey);

      if (databases && !$connections[hostKey]) {
        $connections[hostKey] = { databases: {} };
        databases.sort().forEach(dbKey => {
          $connections[hostKey].databases[dbKey] = $connections[hostKey].databases[dbKey] || { collections: {} };
        });
      }

      busy.end();
    }
  }

  async function selectDatabase(dbKey) {
    info.collKeys = [];
    info.dbKey = dbKey;

    if (dbKey) {
      busy.start();
      const collections = await OpenDatabase(info.hostKey, dbKey);

      for (const collKey of collections?.sort() || []) {
        $connections[info.hostKey].databases[dbKey].collections[collKey] = {};
      }

      busy.end();
    }
  }

  async function performExport() {
    await PerformExport(JSON.stringify(info));
  }

  function selectCollection(collKey) {
    info.collKeys = [ collKey ];
  }
</script>

<Modal bind:show={info} title={actionLabel[info?.type]}>
  <form on:submit|preventDefault={performExport}>
    <!-- svelte-ignore a11y-label-has-associated-control - input is in DirectoryChooser -->
    <label class="field">
      <span class="label">Output destination:</span>
      <DirectoryChooser bind:value={info.outdir} />
      <span class="label">/</span>
      <input type="text" bind:value={info.filename} />
      <span class="label">.</span>
      <select bind:value={info.filetype} class="filetype">
        <optgroup label="Dump (via mongodump)">
          <option value="bson">bson</option>
        </optgroup>
        <optgroup label="Export">
          <option value="csv">csv</option>
          <option value="json">json</option>
        </optgroup>
      </select>
    </label>

    <div class="options">
      {#if info.filetype === 'json'}
        <label class="field">
          <span class="label">Separate items using:</span>
          <select bind:value={info.jsonType}>
            <option value="newline">Newline</option>
            <option value="array">JSON array</option>
          </select>
        </label>
      {/if}
    </div>

    <div class="location">
      <div class="grid">
        <Grid
          key="id"
          columns={[ { title: 'Host', key: 'name' } ]}
          activePath={info ? [ info.hostKey ] : []}
          showHeaders
          hideChildrenToggles
          items={[
            { id: undefined, name: '(localhost)' },
            ...Object.keys(hosts).map(id => ({ id, name: hosts[id]?.name })),
          ]}
          on:select={e => selectHost(e.detail?.itemKey)}
        />
      </div>
      <div class="grid">
        <Grid
          key="id"
          columns={[ { title: 'Database', key: 'name' } ]}
          activePath={info ? [ info.dbKey ] : []}
          showHeaders
          hideChildrenToggles
          items={[
            { id: undefined, name: '(all databases)' },
            ...($connections[info.hostKey]?.databases
              ? Object.keys($connections[info.hostKey].databases).map(id => ({ id, name: id }))
              : []
            ),
          ]}
          on:select={e => selectDatabase(e.detail?.itemKey)}
        />
      </div>
      <div class="grid">
        <Grid
          key="id"
          columns={[ { title: 'Collection', key: 'name' } ]}
          activePath={info?.collKeys ? [ info.collKeys[0] ] : []}
          showHeaders
          hideChildrenToggles
          items={[
            { id: undefined, name: '(all collections)' },
            ...($connections[info.hostKey]?.databases[info.dbKey]?.collections
              ? Object.keys($connections[info.hostKey].databases[info.dbKey].collections).map(id => ({ id, name: id }))
              : []
            ),
          ]}
          on:select={e => selectCollection(e.detail?.itemKey)}
        />
      </div>
    </div>

    <div>
      <button type="submit" class="btn">{actionLabel[info.type]}</button>
    </div>
  </form>
</Modal>

<style>
  form {
    display: grid;
    grid-template: auto auto 1fr auto / 1fr;
    gap: 0.5rem;
  }
  .location {
    display: grid;
    grid-template: 1fr / repeat(3, 1fr);
    gap: 1rem;
  }
  .location .grid {
    border: 1px solid #ccc;
    overflow-y: auto;
  }

  select.filetype {
    flex: 0 1;
  }
</style>
