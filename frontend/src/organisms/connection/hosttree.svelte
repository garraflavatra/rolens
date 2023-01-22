<script>
  import { busy, contextMenu, connections } from '../../stores';
  import { createEventDispatcher } from 'svelte';
  import { DropCollection, DropDatabase, OpenCollection, OpenConnection, OpenDatabase } from '../../../wailsjs/go/app/App';
  import Grid from '../../components/grid.svelte';
  import { WindowSetTitle } from '../../../wailsjs/runtime/runtime';

  export let hosts = {};
  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  const dispatch = createEventDispatcher();
  let activeGridPath = [];
  $: activeHostKey = activeGridPath[0] || activeHostKey;
  $: activeDbKey = activeGridPath[1];
  $: activeCollKey = activeGridPath[2];
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
      WindowSetTitle(`${hosts[activeHostKey].name} - Rolens`);
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
    $connections[activeHostKey].databases[activeDbKey].collections[collKey] = $connections[activeHostKey].databases[activeDbKey].collections[collKey] || {};
    $connections[activeHostKey].databases[activeDbKey].collections[collKey].stats = stats;
    busy.end();
  }

  async function dropCollection(dbKey, collKey) {
    busy.start();
    await DropCollection(activeHostKey, dbKey, collKey);
    await reload();
    busy.end();
  }

  function buildMenu(hostKey, dbKey, collKey, context = [ 'new', 'edit', 'drop' ]) {
    // context: n = new, d = drop
    const menu = [];

    if (context.includes('edit')) {
      hostKey     && menu.push({ label: `Edit host ${hosts[hostKey].name}…`, fn: () => dispatch('editHost', hostKey) });
      collKey     && menu.push({ label: `Rename collection ${collKey}…`, fn: () => dispatch('editCollection', collKey) });
    }
    if (context.includes('drop')) {
      menu.length && menu.push({ separator: true });
      dbKey       && menu.push({ label: `Drop database ${dbKey}…`, fn: () => dropDatabase(dbKey) });
      collKey     && menu.push({ label: `Drop collection ${collKey}…`, fn: () => dropCollection(dbKey, collKey) });
    }
    if (context.includes('new')) {
      menu.length && menu.push({ separator: true });
      hostKey     && menu.push({ label: 'New database…', fn: () => dispatch('newDatabase') });
      dbKey       && menu.push({ label: 'New collection…', fn: () => dispatch('newCollection') });
    }

    return menu;
  }
</script>

<Grid
  striped={false}
  columns={[ { key: 'name' }, { key: 'count', right: true } ]}
  items={Object.keys(hosts).map(hostKey => ({
    id: hostKey,
    name: hosts[hostKey].name,
    icon: 'server',
    children: Object.keys(connection?.databases || {}).sort().map(dbKey => ({
      id: dbKey,
      name: dbKey,
      icon: 'db',
      count: Object.keys(connection.databases[dbKey].collections || {}).length || '',
      children: Object.keys(connection.databases[dbKey].collections).sort().map(collKey => ({
        id: collKey,
        name: collKey,
        icon: 'list',
        menu: buildMenu(hostKey, dbKey, collKey),
      })) || [],
      menu: buildMenu(hostKey, dbKey),
    })),
  }))}
  actions={[
    { icon: 'reload', fn: reload },
    { icon: '+', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'new')) },
    { icon: 'edit', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'edit')), disabled: !activeHostKey },
    { icon: '-', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'drop')), disabled: !activeDbKey },
  ]}
  bind:activePath={activeGridPath}
  on:select={e => {
    const key = e.detail.itemKey;
    switch (e.detail?.level) {
      case 0: return openConnection(key);
      case 1: return openDatabase(key);
      case 2: return openCollection(key);
    }
  }}
/>
