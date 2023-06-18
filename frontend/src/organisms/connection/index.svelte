<script>
  import Icon from '$components/icon.svelte';
  import hostTree from '$lib/stores/hosttree';
  import sharedState from '$lib/stores/sharedstate';
  import CollectionView from './collection/index.svelte';
  import DatabaseView from './database/index.svelte';
  import HostView from './host/index.svelte';
  import HostTree from './hosttree.svelte';

  let path = [];

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
    collection={$hostTree[activeHostKey]?.databases[activeDbKey]?.collections?.[activeCollKey]}
    hostKey={activeHostKey}
    dbKey={activeDbKey}
    collKey={activeCollKey}
  />
{:else if activeDbKey}
  <DatabaseView
    database={$hostTree[activeHostKey]?.databases[activeDbKey]}
    hostKey={activeHostKey}
    dbKey={activeDbKey}
  />
{:else if activeHostKey}
  <HostView
    host={$hostTree[activeHostKey]}
    hostKey={activeHostKey}
  />
{/if}

<style>
  .tree {
    padding: 0.5rem;
    background-color: #fff;
  }
  .tree-buttons {
    margin-bottom: 1rem;
  }
</style>
