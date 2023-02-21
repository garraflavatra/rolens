<script>
  import busy from '$lib/stores/busy';
  import { connections } from '$lib/stores/connections';
  import { Hosts, RenameCollection } from '$wails/go/app/App';
  import { EnterText } from '$wails/go/ui/UI';
  import { EventsOn } from '$wails/runtime/runtime';
  import { onMount } from 'svelte';
  import CollectionDetail from './collection/index.svelte';
  import Export from './export/export.svelte';
  import HostDetail from './hostdetail.svelte';
  import HostTree from './hosttree.svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  let hostTree;
  let showHostDetail = false;
  let hostDetailKey = '';
  let exportInfo;

  async function getHosts() {
    hosts = await Hosts();
  }

  function createHost() {
    hostDetailKey = '';
    showHostDetail = true;
  }

  function editHost(hostKey) {
    hostDetailKey = hostKey;
    showHostDetail = true;
  }

  async function createDatabase() {
    const name = await EnterText('Create a database', 'Enter the database name. Note: databases in MongoDB do not exist until they have a collection and an item. Your new database will not persist on the server; fill it to have it created.');
    if (name) {
      $connections[activeHostKey].databases[name] = { collections: {} };
    }
  }

  async function renameCollection(oldCollKey) {
    const newCollKey = await EnterText('Rename collection', `Enter a new name for collection ${oldCollKey}.`, oldCollKey);
    if (newCollKey && (newCollKey !== oldCollKey)) {
      busy.start();
      const ok = await RenameCollection(activeHostKey, activeDbKey, oldCollKey, newCollKey);
      if (ok) {
        activeCollKey = newCollKey;
        await hostTree.reload();
      }
      busy.end();
    }
  }

  async function createCollection() {
    const name = await EnterText('Create a collection', 'Note: collections in MongoDB do not exist until they have at least one item. Your new collection will not persist on the server; fill it to have it created.');
    if (name) {
      $connections[activeHostKey].databases[activeDbKey].collections[name] = {};
    }
  }

  function exportCollection(collKey) {
    exportInfo = {
      type: 'export',
      filetype: 'json',
      hostKey: activeHostKey,
      dbKey: activeDbKey,
      collKeys: [ collKey ],
    };
  }

  function dumpCollection(collKey) {
    exportInfo = {
      type: 'dump',
      filetype: 'bson',
      hostKey: activeHostKey,
      dbKey: activeDbKey,
      collKeys: [ collKey ],
    };
  }

  EventsOn('CreateHost', createHost);
  EventsOn('CreateDatabase', createDatabase);
  EventsOn('CreateCollection', createCollection);
  onMount(getHosts);
</script>

<div class="tree">
  <HostTree
    {hosts}
    bind:activeHostKey
    bind:activeCollKey
    bind:activeDbKey
    bind:this={hostTree}
    on:newHost={createHost}
    on:newDatabase={createDatabase}
    on:newCollection={createCollection}
    on:editHost={e => editHost(e.detail)}
    on:renameCollection={e => renameCollection(e.detail)}
    on:exportCollection={e => exportCollection(e.detail)}
    on:dumpCollection={e => dumpCollection(e.detail)}
  />
</div>

<CollectionDetail
  collection={$connections[activeHostKey]?.databases[activeDbKey]?.collections?.[activeCollKey]}
  hostKey={activeHostKey}
  dbKey={activeDbKey}
  collectionKey={activeCollKey}
  {hosts}
/>

<HostDetail
  bind:show={showHostDetail}
  on:reload={getHosts}
  hostKey={activeHostKey}
  {hosts}
/>

<Export bind:info={exportInfo} {hosts} />

<style>
  .tree {
    padding: 0.5rem;
    background-color: #fff;
  }
</style>
