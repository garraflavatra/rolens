<script>
  import DirectoryChooser from '$components/editors/directorychooser.svelte';
  import Grid from '$components/grid/grid.svelte';
  import Modal from '$components/modal.svelte';
  import hostTree from '$lib/stores/hosttree.js';
  import applicationSettings from '$lib/stores/settings.js';
  import { OpenConnection, OpenDatabase } from '$wails/go/app/App.js';
  import { createEventDispatcher } from 'svelte';

  export let info = {};

  const dispatch = createEventDispatcher();

  $: if (info) {
    info.outdir = info.outdir || $applicationSettings.defaultExportDirectory;
    info.filename = info.filename || `Dump ${new Date().getTime()}`;
  }

  async function selectHost(hostKey) {
    info.hostKey = hostKey;
    info.dbKey = undefined;
    info.collKeys = [];

    if (hostKey) {
      const databases = await OpenConnection(hostKey);

      if (databases && !$hostTree[hostKey]) {
        $hostTree[hostKey] = { databases: {} };
        databases.sort().forEach(dbKey => {
          $hostTree[hostKey].databases[dbKey] = $hostTree[hostKey].databases[dbKey] || { collections: {} };
        });
      }
    }
  }

  async function selectDatabase(dbKey) {
    info.collKeys = [];
    info.dbKey = dbKey;

    if (dbKey) {
      const collections = await OpenDatabase(info.hostKey, dbKey);

      for (const collKey of collections?.sort() || []) {
        $hostTree[info.hostKey].databases[dbKey].collections[collKey] = {};
      }
    }
  }

  function selectCollection(collKey) {
    info.collKeys = [ collKey ];
  }

  function performDump() {
    dispatch('dump', { info });
  }
</script>

<Modal title="Perform dump" on:close>
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
            ...Object.keys($hostTree).map(id => {
              return { id, name: $hostTree[id]?.name };
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
            ...($hostTree[info.hostKey]?.databases ? Object.keys($hostTree[info.hostKey].databases).map(id => {
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
            ...($hostTree[info.hostKey]?.databases[info.dbKey]?.collections ? Object.keys($hostTree[info.hostKey].databases[info.dbKey].collections).map(id => {
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
    <button class="button" on:click={performDump}>Perform dump</button>
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
