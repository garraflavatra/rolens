<script>
  import { startProgress } from '$lib/progress';
  import connections from '$lib/stores/connections';
  import { Hosts, RenameCollection } from '$wails/go/app/App';
  import { EnterText } from '$wails/go/ui/UI';
  import { EventsOn } from '$wails/runtime/runtime';
  import { onMount } from 'svelte';
  import HostView from './host/index.svelte';
  import DatabaseView from './database/index.svelte';
  import CollectionView from './collection/index.svelte';
  import DumpInfo from './dump.svelte';
  import HostDetail from './hostdetail.svelte';
  import HostTree from './hosttree.svelte';
  import sharedState from '$lib/stores/sharedstate';
  import Icon from '$components/icon.svelte';
  import hosts from '$lib/stores/hosts';

  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  let hostTree;
  let showHostDetail = false;
  let hostDetailKey = '';
  let exportInfo;

  $: sharedState.currentHost.set(activeHostKey);
  $: sharedState.currentDb.set(activeDbKey);
  $: sharedState.currentColl.set(activeCollKey);

  export function createHost() {
    hostDetailKey = '';
    showHostDetail = true;
  }

  function editHost(hostKey) {
    hostDetailKey = hostKey;
    showHostDetail = true;
  }

  export async function createDatabase() {
    const name = await EnterText('Create a database', 'Enter the database name. Note: databases in MongoDB do not exist until they have a collection and an item. Your new database will not persist on the server; fill it to have it created.');
    if (name) {
      $connections[activeHostKey].databases[name] = { collections: {} };
    }
  }

  async function renameCollection(oldCollKey) {
    const newCollKey = await EnterText('Rename collection', `Enter a new name for collection ${oldCollKey}.`, oldCollKey);
    if (newCollKey && (newCollKey !== oldCollKey)) {
      const progress = startProgress(`Renaming collection "${oldCollKey}" to "${newCollKey}"…`);
      const ok = await RenameCollection(activeHostKey, activeDbKey, oldCollKey, newCollKey);
      if (ok) {
        activeCollKey = newCollKey;
        await hostTree.reload();
      }
      progress.end();
    }
  }

  export async function createCollection() {
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
</script>

<div class="tree">
  <div class="tree-buttons">
    <button class="button-small" on:click={createHost}>
      <Icon name="+" /> New host
    </button>
  </div>

  <HostTree
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

{#if activeCollKey}
  <CollectionView
    collection={$connections[activeHostKey]?.databases[activeDbKey]?.collections?.[activeCollKey]}
    hostKey={activeHostKey}
    dbKey={activeDbKey}
    collKey={activeCollKey}
  />
{:else if activeDbKey}
  <DatabaseView
    database={$connections[activeHostKey]?.databases[activeDbKey]}
    hostKey={activeHostKey}
    dbKey={activeDbKey}
  />
{:else if activeHostKey}
  <HostView
    host={$connections[activeHostKey]}
    hostKey={activeHostKey}
  />
{/if}

<HostDetail
  bind:show={showHostDetail}
  on:reload={hosts.update}
  hostKey={activeHostKey}
/>

<DumpInfo bind:info={exportInfo} />

<style>
  .tree {
    padding: 0.5rem;
    background-color: #fff;
  }
  .tree-buttons {
    margin-bottom: 1rem;
  }
</style>
