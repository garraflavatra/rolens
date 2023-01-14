<script>
  import { onMount } from 'svelte';
  import { DropDatabase, Hosts, OpenCollection, OpenConnection, OpenDatabase } from '../wailsjs/go/app/App';
  import AddressBar from './organisms/addressbar/index.svelte';
  import Grid from './components/grid.svelte';
  import CollectionDetail from './organisms/collection-detail/index.svelte';
  import { busy, contextMenu } from './stores';
  import ContextMenu from './components/contextmenu.svelte';

  const connections = {};
  let hosts = {};
  let activeHostKey = '';
  let activeDbKey = '';
  let activeCollKey = '';
  let addressBarModalOpen = true;

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

  async function dropDatabase(dbKey) {
    $busy = true;
    await DropDatabase(activeHostKey, dbKey);
    await openConnection(activeHostKey);
    $busy = false();
  }

  async function openCollection(collKey) {
    $busy = true;
    const stats = await OpenCollection(activeHostKey, activeDbKey, collKey);
    connections[activeHostKey].databases[activeDbKey].collections[collKey].stats = stats;
    $busy = false;
  }

  async function reload() {
    activeHostKey && await openConnection(activeHostKey);
    activeDbKey && await openDatabase(activeDbKey);
  }

  onMount(() => {
    Hosts().then(h => hosts = h);
  });
</script>

<main>
  <AddressBar {hosts} bind:activeHostKey on:select={e => openConnection(e.detail)} bind:modalOpen={addressBarModalOpen} />

  {#if host && connection}
    <div class="databaselist">
      <Grid
        columns={[ { key: 'id' }, { key: 'collCount', right: true } ]}
        items={Object.keys(connection.databases).map(id => ({
          id,
          collCount: Object.keys(connection.databases[id].collections || {}).length,
          children: connection.databases[id].collections || [],
          menu: [ { label: `Drop ${id}`, fn: () => dropDatabase(id) } ],
        }))}
        actions={[
          { icon: 'reload', fn: reload },
          { icon: '+' },
          { icon: '-' },
        ]}
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
</main>

{#key $contextMenu}
  <ContextMenu {...$contextMenu} on:close={contextMenu.hide} />
{/key}

<style>
  main {
    height: 100vh;
    display: grid;
    grid-template: 3rem auto / 250px 1fr;
    gap: 0.5rem;
    padding: 0.5rem;
  }
  main > :global(.addressbar) {
    grid-column: 1 / 3;
  }

  .databaselist {
    overflow: scroll;
  }
</style>
