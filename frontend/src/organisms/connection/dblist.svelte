<script>
  import { busy, contextMenu, connections } from '../../stores';
  import { createEventDispatcher } from 'svelte';
  import { DropCollection, DropDatabase, OpenCollection, OpenConnection, OpenDatabase } from '../../../wailsjs/go/app/App';
  import Grid from '../../components/grid.svelte';
  import { WindowSetTitle } from '../../../wailsjs/runtime';

  export let hosts = {};
  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  const dispatch = createEventDispatcher();
  $: host = hosts[activeHostKey];
  $: connection = $connections[activeHostKey];
  $: database = connection?.databases[activeDbKey];
  $: collection = database?.collections?.[activeCollKey];

  export async function reload() {
    activeHostKey && await openConnection(activeHostKey);
    activeDbKey && await openDatabase(activeDbKey);
    activeCollKey && await openCollection(activeCollKey);
  }

  async function openConnection(hostKey) {
    busy.start();
    const databases = await OpenConnection(hostKey);

    if (databases) {
      $connections[hostKey] = { databases: {} };
      databases.forEach(dbKey => {
        $connections[hostKey].databases[dbKey] = { collections: {} };
      });
      activeHostKey = hostKey;
      dispatch('connected', hostKey);
      WindowSetTitle(`${hosts[activeHostKey].name} - Mongodup`);
    }

    busy.end();
  }

  async function openDatabase(dbKey) {
    busy.start();
    const collections = await OpenDatabase(activeHostKey, dbKey);

    for (const collKey of collections || []) {
      $connections[activeHostKey].databases[dbKey].collections[collKey] = {};
    }

    busy.end();
  }

  async function dropDatabase(dbKey) {
    busy.start();
    await DropDatabase(activeHostKey, dbKey);
    await reload();
    busy.end();
  }

  async function openCollection(collKey) {
    busy.start();
    const stats = await OpenCollection(activeHostKey, activeDbKey, collKey);
    $connections[activeHostKey].databases[activeDbKey].collections[collKey].stats = stats;
    busy.end();
  }

  async function dropCollection(dbKey, collKey) {
    busy.start();
    await DropCollection(activeHostKey, dbKey, collKey);
    await reload();
    busy.end();
  }
</script>

{#if host && connection}
  <Grid
    striped={false}
    columns={[ { key: 'id' }, { key: 'collCount', right: true } ]}
    items={Object.keys(connection.databases).map(dbKey => ({
      id: dbKey,
      collCount: Object.keys(connection.databases[dbKey].collections || {}).length || '',
      children: Object.keys(connection.databases[dbKey].collections).map(collKey => ({
        id: collKey,
        menu: [
          { label: `Drop ${collKey}…`, fn: () => dropCollection(dbKey, collKey) },
          { label: `Drop ${dbKey}…`, fn: () => dropDatabase(dbKey) },
          { separator: true },
          { label: 'New database…', fn: () => dispatch('newDatabase') },
          { label: 'New collection…', fn: () => dispatch('newCollection') },
        ],
      })).sort((a, b) => a.id.localeCompare(b)) || [],
      menu: [
        { label: `Drop ${dbKey}…`, fn: () => dropDatabase(dbKey) },
        { separator: true },
        { label: 'New database…', fn: () => dispatch('newDatabase') },
        { label: 'New collection…', fn: () => dispatch('newCollection') },
      ],
    }))}
    actions={[
      { icon: 'reload', fn: reload },
      { icon: '+', fn: evt => {
        if (activeDbKey) {
          contextMenu.show(evt, [
            { label: 'New database…', fn: () => dispatch('newDatabase') },
            { label: 'New collection…', fn: () => dispatch('newCollection') },
          ]);
        }
        else {
          dispatch('newDatabase');
        }
      } },
      { icon: '-', fn: evt => {
        if (activeCollKey) {
          contextMenu.show(evt, [
            { label: 'Drop database…', fn: () => dropDatabase(activeDbKey) },
            { label: 'Drop collection…', fn: () => dropCollection(activeDbKey, activeCollKey) },
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
{/if}
