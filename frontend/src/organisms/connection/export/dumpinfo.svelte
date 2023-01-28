<script>
  import { busy, connections } from '../../../stores';
  import Grid from '../../../components/grid.svelte';
  import Modal from '../../../components/modal.svelte';
  import { OpenConnection, OpenDatabase } from '../../../../wailsjs/go/app/App';

  export let info;
  export let hosts = {};

  $: console.log(info);

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

  function selectCollection(collKey) {
    info.collKeys = [ collKey ];
  }
</script>

<Modal bind:show={info} title="Dump database data">
  <div class="info">
    <div class="meta">
      <label class="field">
        <span class="label">Output filename</span>
        <input type="text">
      </label>
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
  </div>
</Modal>

<style>
  .info {
    display: grid;
    grid-template: auto / 1fr;
  }
  .location {
    display: grid;
    grid-template: 1fr / repeat(3, 1fr);
    gap: 1rem;
  }
  .location .grid {
    border: 1px solid #ccc;
    padding: 0.3rem;
    overflow-y: auto;
  }
</style>
