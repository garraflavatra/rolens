<script>
  import { onMount } from 'svelte';
  import { Hosts, OpenCollection, OpenConnection, OpenDatabase } from '../wailsjs/go/main/App';
  import AddressBar from './organisms/addressbar/index.svelte';
  import Grid from './components/grid.svelte';
  import CollectionDetail from './organisms/collection-detail/index.svelte';
  import { busy } from './stores';

  const connections = {};
  let hosts = {};
  let activeHostKey = '';
  let activeDbKey = '';
  let activeCollKey = '';
  let addressBarModalOpen = false;

  $: host = hosts[activeHostKey];
  $: connection = connections[activeHostKey];
  $: database = connection?.databases[activeDbKey];
  $: collection = database?.collections?.[activeCollKey];

  async function openConnection(hostKey) {
    $busy = true;
    const databases = await OpenConnection(hostKey);

    if (databases) {
      connections[hostKey] = { databases: {} };
      databases.forEach(dbKey => {
        connections[hostKey].databases[dbKey] = { collections: {} };
      });
      activeHostKey = hostKey;
      addressBarModalOpen = false;
      window.runtime.WindowSetTitle(`${host.name} - Mongodup`);
    }

    $busy = false;
  }

  async function openDatabase(dbKey) {
    $busy = true;
    const collections = await OpenDatabase(activeHostKey, dbKey);

    for (const collKey of collections || []) {
      connections[activeHostKey].databases[dbKey].collections[collKey] = {};
    }

    $busy = false;
  }

  async function openCollection(collKey) {
    $busy = true;
    const stats = await OpenCollection(activeHostKey, activeDbKey, collKey);
    connections[activeHostKey].databases[activeDbKey].collections[collKey].stats = stats;
    $busy = false;
  }

  onMount(() => {
    Hosts().then(h => hosts = h);
  });
</script>

<main>
  <AddressBar {hosts} bind:activeHostKey on:select={e => openConnection(e.detail)} bind:modalOpen={addressBarModalOpen} />

  <div class="columns">
    {#if host && connection}
      <div class="hostlist">
        <Grid
          columns={[ { key: 'id' }, { key: 'collCount', right: true } ]}
          items={Object.keys(connection.databases).map(id => ({
            id,
            collCount: Object.keys(connection.databases[id].collections || {}).length,
            children: connection.databases[id].collections || [],
          }))}
          bind:activeKey={activeDbKey}
          bind:activeChildKey={activeCollKey}
          on:select={e => openDatabase(e.detail)}
          on:selectChild={e => openCollection(e.detail)}
        />
      </div>

      <div class="collection">
        <CollectionDetail
          {collection}
          hostKey={activeHostKey}
          dbKey={activeDbKey}
          collectionKey={activeCollKey}
        />
      </div>
    {/if}
  </div>
</main>

<style>
  main {
    height: 100vh;
    display: flex;
    flex-flow: column;
  }

  .columns {
    display: flex;
    gap: 1rem;
    flex: 1;
    height: 100%;
  }
  .columns > :global(*) {
    height: 100%;
    display: flex;
    flex-flow: column;
  }

  .hostlist {
    flex: 0 0 250px;
    overflow: scroll;
  }
  .collection {
    flex: 1;
    width: auto;
  }
</style>
