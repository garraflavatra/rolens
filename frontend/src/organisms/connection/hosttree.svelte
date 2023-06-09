<script>
  import Grid from '$components/grid.svelte';
  import { startProgress } from '$lib/progress';
  import connections from '$lib/stores/connections';
  import { createEventDispatcher } from 'svelte';
  import { DropCollection, DropDatabase, OpenCollection, OpenConnection, OpenDatabase, RemoveHost, TruncateCollection } from '../../../wailsjs/go/app/App';
  import hosts from '$lib/stores/hosts';
  import { tick } from 'svelte';
  import windowTitle from '$lib/stores/windowtitle';

  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  const dispatch = createEventDispatcher();
  let activeGridPath = [];
  // $: activeGridPath[0] = activeHostKey || undefined;
  // $: activeGridPath[1] = activeDbKey || undefined;
  // $: activeGridPath[2] = activeCollKey || undefined;
  $: host = $hosts[activeHostKey];
  $: connection = $connections[activeHostKey];
  $: database = connection?.databases[activeDbKey];
  $: collection = database?.collections?.[activeCollKey];

  export async function reload() {
    activeHostKey && await openConnection(activeHostKey);
    activeDbKey && await openDatabase(activeDbKey);
    activeCollKey && await openCollection(activeCollKey);
  }

  async function openConnection(hostKey) {
    const progress = startProgress(`Connecting to "${hostKey}"…`);

    activeCollKey = '';
    activeDbKey = '';
    activeHostKey = hostKey;

    const { databases, status, systemInfo } = await OpenConnection(hostKey);

    if (databases) {
      $connections[hostKey] = $connections[hostKey] || {};
      $connections[hostKey].status = status;
      $connections[hostKey].systemInfo = systemInfo;

      $connections[hostKey].databases = $connections[hostKey].databases || {};
      databases.forEach(dbKey => {
        $connections[hostKey].databases[dbKey] =
          $connections[hostKey].databases[dbKey] || { collections: {} };
      });

      activeHostKey = hostKey;
      dispatch('connected', hostKey);
    }

    progress.end();

    if (databases) {
      windowTitle.setSegments($hosts[activeHostKey].name, 'Rolens');
    }
  }

  async function removeHost(hostKey) {
    activeCollKey = '';
    activeDbKey = '';
    activeHostKey = '';

    await tick();
    await RemoveHost(hostKey);
    await reload();
    await hosts.update();
  }

  async function openDatabase(dbKey) {
    const progress = startProgress(`Opening database "${dbKey}"…`);
    const { collections, stats } = await OpenDatabase(activeHostKey, dbKey);
    activeDbKey = dbKey;
    activeCollKey = '';
    $connections[activeHostKey].databases[dbKey].stats = stats;

    for (const collKey of collections || []) {
      $connections[activeHostKey].databases[dbKey].collections[collKey] =
        $connections[activeHostKey].databases[dbKey].collections[collKey] ||{};
    }

    progress.end();
    windowTitle.setSegments(activeDbKey, $hosts[activeHostKey].name, 'Rolens');
  }

  async function dropDatabase(dbKey) {
    const progress = startProgress(`Dropping database "${dbKey}"…`);
    const success = await DropDatabase(activeHostKey, dbKey);
    if (success) {
      activeCollKey = '';
      activeDbKey = '';
      await reload();
    }
    progress.end();
  }

  async function openCollection(collKey) {
    const progress = startProgress(`Opening collection "${collKey}"…`);
    const stats = await OpenCollection(activeHostKey, activeDbKey, collKey);
    activeCollKey = collKey;
    $connections[activeHostKey].databases[activeDbKey].collections[collKey] = $connections[activeHostKey].databases[activeDbKey].collections[collKey] || {};
    $connections[activeHostKey].databases[activeDbKey].collections[collKey].stats = stats;
    progress.end();
    windowTitle.setSegments(activeDbKey + '.' + activeCollKey, $hosts[activeHostKey].name, 'Rolens');
  }

  async function truncateCollection(dbKey, collKey) {
    const progress = startProgress(`Truncating collection "${collKey}"…`);
    await TruncateCollection(activeHostKey, dbKey, collKey);
    await reload();
    progress.end();
  }

  async function dropCollection(dbKey, collKey) {
    const progress = startProgress(`Dropping collection "${collKey}"…`);
    const success = await DropCollection(activeHostKey, dbKey, collKey);
    if (success) {
      activeCollKey = '';
      await reload();
    }
    progress.end();
  }
</script>

<Grid
  striped={false}
  columns={[ { key: 'name' }, { key: 'count', right: true } ]}
  bind:activePath={activeGridPath}
  items={Object.keys($hosts).map(hostKey => ({
    id: hostKey,
    name: $hosts[hostKey].name,
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
        menu: [
          { label: 'Export collection (JSON, CSV)…', fn: () => dispatch('exportCollection', collKey) },
          { label: 'Dump collection (BSON via mongodump)…', fn: () => dispatch('dumpCollection', collKey) },
          { separator: true },
          { label: 'Rename collection…', fn: () => dispatch('renameCollection', collKey) },
          { label: 'Truncate collection…', fn: () => truncateCollection(dbKey, collKey) },
          { label: 'Drop collection…', fn: () => dropCollection(dbKey, collKey) },
          { separator: true },
          { label: 'New collection…', fn: () => dispatch('newCollection') },
        ],
      })) || [],
      menu: [
        { label: 'Drop database…', fn: () => dropDatabase(dbKey) },
        { separator: true },
        { label: 'New database…', fn: () => dispatch('newDatabase') },
        { label: 'New collection…', fn: () => dispatch('newCollection') },
      ],
    })),
    menu: [
      { label: 'New database…', fn: () => dispatch('newDatabase') },
      { separator: true },
      { label: `Edit host ${$hosts[hostKey].name}…`, fn: () => dispatch('editHost', hostKey) },
      { label: `Remove host…`, fn: () => removeHost(hostKey) },
    ],
  }))}
  on:select={e => {
    const key = e.detail.itemKey;
    switch (e.detail?.level) {
      case 0: return openConnection(key);
      case 1: return openDatabase(key);
      case 2: return openCollection(key);
    }
  }}
/>

<!--
  actions={[
    { icon: 'reload', fn: reload },
    { icon: '+', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'new')) },
    { icon: 'edit', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'edit')), disabled: !activeHostKey },
    { icon: '-', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'drop')), disabled: !activeDbKey },
  ]}
-->
