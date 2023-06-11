<script>
  import DirectoryChooser from '$components/directorychooser.svelte';
  import Grid from '$components/grid.svelte';
  import Modal from '$components/modal.svelte';
  import { startProgress } from '$lib/progress';
  import connections from '$lib/stores/connections';
  import hosts from '$lib/stores/hosts';
  import applicationSettings from '$lib/stores/settings';
  import { OpenConnection, OpenDatabase, PerformDump } from '$wails/go/app/App';

  export let info;

  $: if (info) {
    info.outdir = info.outdir || $applicationSettings.defaultExportDirectory;
    info.filename = info.filename || `Dump ${new Date().getTime()}`;
  }

  async function selectHost(hostKey) {
    info.hostKey = hostKey;
    info.dbKey = undefined;
    info.collKeys = [];

    if (hostKey) {
      const progress = startProgress(`Opening connection to host "${hostKey}"`);
      const databases = await OpenConnection(hostKey);

      if (databases && !$connections[hostKey]) {
        $connections[hostKey] = { databases: {} };
        databases.sort().forEach(dbKey => {
          $connections[hostKey].databases[dbKey] = $connections[hostKey].databases[dbKey] || { collections: {} };
        });
      }

      progress.end();
    }
  }

  async function selectDatabase(dbKey) {
    info.collKeys = [];
    info.dbKey = dbKey;

    if (dbKey) {
      const progress = startProgress(`Opening database "${dbKey}"`);
      const collections = await OpenDatabase(info.hostKey, dbKey);

      for (const collKey of collections?.sort() || []) {
        $connections[info.hostKey].databases[dbKey].collections[collKey] = {};
      }

      progress.end();
    }
  }

  async function performDump() {
    const ok = await PerformDump(JSON.stringify(info));
    if (ok) {
      info = undefined;
    }
  }

  function selectCollection(collKey) {
    info.collKeys = [ collKey ];
  }
</script>

<Modal bind:show={info} title="Perform dump">
  <form on:submit|preventDefault={performDump}>
    <label class="field">
      <span class="label">Output destination:</span>
      <DirectoryChooser bind:value={info.outdir} />
      <span class="label">/</span>
      <input type="text" bind:value={info.filename} />
    </label>

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
            ...Object.keys($hosts).map(id => {
              return { id, name: $hosts[id]?.name };
            }),
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
            ...($connections[info.hostKey]?.databases ? Object.keys($connections[info.hostKey].databases).map(id => {
              return { id, name: id };
            }) : []
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
            ...($connections[info.hostKey]?.databases[info.dbKey]?.collections ? Object.keys($connections[info.hostKey].databases[info.dbKey].collections).map(id => {
              return { id, name: id };
            }) : []
            ),
          ]}
          on:select={e => selectCollection(e.detail?.itemKey)}
        />
      </div>
    </div>
  </form>

  <svelte:fragment slot="footer">
    <button class="btn" on:click={performDump}>Perform dump</button>
  </svelte:fragment>
</Modal>

<style>
  form {
    display: grid;
    grid-template: auto 1fr auto / 1fr;
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
</style>
