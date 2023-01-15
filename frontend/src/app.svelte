<script>
  import { onMount } from 'svelte';
  import { DropCollection, DropDatabase, Hosts, OpenCollection, OpenConnection, OpenDatabase } from '../wailsjs/go/app/App';
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
    busy.start();
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

    busy.end();
  }

  async function openDatabase(dbKey) {
    busy.start();
    const collections = await OpenDatabase(activeHostKey, dbKey);

    for (const collKey of collections || []) {
      connections[activeHostKey].databases[dbKey].collections[collKey] = {};
    }

    busy.end();
  }

  async function dropDatabase(dbKey) {
    busy.start();
    await DropDatabase(activeHostKey, dbKey);
    await openConnection(activeHostKey);
    busy.end();
  }

  async function openCollection(collKey) {
    busy.start();
    const stats = await OpenCollection(activeHostKey, activeDbKey, collKey);
    connections[activeHostKey].databases[activeDbKey].collections[collKey].stats = stats;
    busy.end();
  }

  async function dropCollection(dbKey, collKey) {
    busy.start();
    await DropCollection(activeHostKey, dbKey, collKey);
    await openConnection(activeHostKey);
    await openDatabase(dbKey);
    busy.end();
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
        items={Object.keys(connection.databases).map(dbKey => ({
          id: dbKey,
          collCount: Object.keys(connection.databases[dbKey].collections || {}).length,
          children: Object.keys(connection.databases[dbKey].collections).map(collKey => ({
            id: collKey,
            menu: [ { label: `Drop ${collKey}`, fn: () => dropCollection(dbKey, collKey) } ],
          })) || [],
          menu: [ { label: `Drop ${dbKey}`, fn: () => dropDatabase(dbKey) } ],
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
