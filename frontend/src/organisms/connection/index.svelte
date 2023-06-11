<script>
  import HostView from './host/index.svelte';
  import DatabaseView from './database/index.svelte';
  import CollectionView from './collection/index.svelte';
  import DumpInfo from './database/dialogs/dump.svelte';
  import HostDetail from './host/dialogs/hostdetail.svelte';
  import HostTree from './hosttree.svelte';
  import sharedState from '$lib/stores/sharedstate';
  import Icon from '$components/icon.svelte';
  import { writable } from 'svelte/store';

  let hostTree;
  let showHostDetail = false;
  const hostDetailKey = '';
  let exportInfo;
  let path = [];

  // @todo
  const connections = writable({});

  $: activeHostKey = path[0];
  $: activeDbKey = path[1];
  $: activeCollKey = path[2];

  $: sharedState.currentHost.set(activeHostKey);
  $: sharedState.currentDb.set(activeDbKey);
  $: sharedState.currentColl.set(activeCollKey);
</script>

<div class="tree">
  <div class="tree-buttons">
    <button class="button-small" on:click={hostTree.newHost}>
      <Icon name="+" /> New host
    </button>
  </div>

  <HostTree bind:path />
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
  on:reload={hostTree.refresh}
  hostKey={hostDetailKey}
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
