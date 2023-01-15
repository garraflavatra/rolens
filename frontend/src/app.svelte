<script>
  import { onMount, tick } from 'svelte';
  import { DropCollection, DropDatabase, Hosts, OpenCollection, OpenConnection, OpenDatabase } from '../wailsjs/go/app/App';
  import AddressBar from './organisms/addressbar/index.svelte';
  import Grid from './components/grid.svelte';
  import CollectionDetail from './organisms/collection-detail/index.svelte';
  import { busy, contextMenu } from './stores';
  import ContextMenu from './components/contextmenu.svelte';
  import Modal from './components/modal.svelte';
  import { input } from './actions';

  const connections = {};
  let hosts = {};

  let activeHostKey = '';
  let activeDbKey = '';
  let activeCollKey = '';

  let addressBarModalOpen = true;

  let newDb;
  let newDbInput;

  let newColl;
  let newCollInput;

  $: host = hosts[activeHostKey];
  $: connection = connections[activeHostKey];
  $: database = connection?.databases[activeDbKey];
  $: collection = database?.collections?.[activeCollKey];

  $: if (newDb) {
    tick().then(() => newDbInput.focus());
  }

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

  function createDatabase() {
    busy.start();
    connections[activeHostKey].databases[newDb.name] = { collections: {} };
    newDb = undefined;
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
    await reload();
    busy.end();
  }

  function createCollection() {
    busy.start();
    connections[activeHostKey].databases[activeDbKey].collections[newColl.name] = {};
    newColl = undefined;
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
    await reload();
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
          collCount: Object.keys(connection.databases[dbKey].collections || {}).length || '',
          children: Object.keys(connection.databases[dbKey].collections).map(collKey => ({
            id: collKey,
            menu: [ { label: `Drop ${collKey}`, fn: () => dropCollection(dbKey, collKey) } ],
          })).sort((a, b) => a.id.localeCompare(b)) || [],
          menu: [ { label: `Drop ${dbKey}`, fn: () => dropDatabase(dbKey) } ],
        }))}
        actions={[
          { icon: 'reload', fn: reload },
          { icon: '+', fn: evt => {
            if (activeDbKey) {
              contextMenu.show(evt, [
                { label: 'New database', fn: () => newDb = {} },
                { label: 'New collection', fn: () => newColl = {} },
              ]);
            }
            else {
              newDb = {};
            }
          } },
          { icon: '-', fn: evt => {
            if (activeCollKey) {
              contextMenu.show(evt, [
                { label: 'Drop database', fn: () => dropDatabase(activeDbKey) },
                { label: 'Drop collection', fn: () => dropCollection(activeDbKey, activeCollKey) },
              ]);
            }
            else {
              dropDatabase(activeDbKey);
            }
          }, disabled: !activeDbKey },
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

{#if newDb}
  <Modal bind:show={newDb}>
    <p><strong>Create a database</strong></p>
    <p>Note: databases in MongoDB do not exist until they have a collection and an item. Your new database will not persist on the server; fill it to have it created.</p>
    <form on:submit|preventDefault={createDatabase}>
      <label class="field">
        <input type="text" spellcheck="false" bind:value={newDb.name} use:input placeholder="New collection name" bind:this={newDbInput} />
      </label>
      <button class="btn create" type="submit">Create database</button>
    </form>
  </Modal>
{/if}

{#if newColl}
  <Modal bind:show={newColl}>
    <p><strong>Create a collections</strong></p>
    <p>Note: collections in MongoDB do not exist until they have at least one item. Your new collection will not persist on the server; fill it to have it created.</p>
    <form on:submit|preventDefault={createCollection}>
      <label class="field">
        <input type="text" spellcheck="false" bind:value={newColl.name} use:input placeholder="New collection name" bind:this={newCollInput} />
      </label>
      <button class="btn create" type="submit">Create collection</button>
    </form>
  </Modal>
{/if}

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

  .btn.create {
    margin-top: 0.5rem;
  }
</style>
