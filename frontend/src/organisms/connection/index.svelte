<script>
  import Icon from '$components/icon.svelte';
  import hostTree from '$lib/stores/hosttree.js';
  import sharedState from '$lib/stores/sharedstate.js';
  import { EventsOn } from '$wails/runtime/runtime.js';
  import CollectionView from './collection/index.svelte';
  import DatabaseView from './database/index.svelte';
  import HostView from './host/index.svelte';
  import HostTree from './hosttree.svelte';

  let path = [];
  let prevPath = '';
  let hostTab = '';
  let dbTab = '';
  let collTab = '';

  $: if (path.join('.') !== prevPath) {
    hostTab = 'status';
    dbTab = 'stats';
    collTab = 'find';
    prevPath = path.join('.');
  }

  $: activeHostKey = path[0];
  $: activeDbKey = path[1];
  $: activeCollKey = path[2];

  $: sharedState.currentHost.set(activeHostKey);
  $: sharedState.currentDb.set(activeDbKey);
  $: sharedState.currentColl.set(activeCollKey);

  EventsOn('ui.host.new', () => hostTree.newHost());
  EventsOn('ui.host.edit', () => $hostTree[activeHostKey]?.edit());
  EventsOn('ui.host.remove', () => $hostTree[activeHostKey]?.remove());
  EventsOn('ui.host.tab', tab => {
    path = path.slice(0, 1);
    hostTab = tab;
  });

  EventsOn('ui.db.new', () => $hostTree[activeHostKey]?.newDatabase());
  EventsOn('ui.db.dump', () => $hostTree[activeHostKey]?.databases[activeDbKey]?.dump());
  EventsOn('ui.db.drop', () => $hostTree[activeHostKey]?.databases[activeDbKey]?.drop());
  EventsOn('ui.db.tab', tab => {
    path = path.slice(0, 2);
    dbTab = tab;
  });

  EventsOn('ui.coll.new', () => $hostTree[activeHostKey]?.databases[activeDbKey]?.newCollection());
  EventsOn('ui.coll.export', () => $hostTree[activeHostKey]?.databases[activeDbKey]?.collections[activeCollKey]?.export());
  EventsOn('ui.coll.truncate', () => $hostTree[activeHostKey]?.databases[activeDbKey]?.collections[activeCollKey]?.truncate());
  EventsOn('ui.coll.drop', () => $hostTree[activeHostKey]?.databases[activeDbKey]?.collections[activeCollKey]?.drop());
  EventsOn('ui.coll.tab', tab => collTab = tab);
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
  {#key activeCollKey}
    <CollectionView
      host={$hostTree[activeHostKey]}
      database={$hostTree[activeHostKey]?.databases[activeDbKey]}
      collection={$hostTree[activeHostKey]?.databases[activeDbKey]?.collections?.[activeCollKey]}
      bind:tab={collTab}
    />
  {/key}
{:else if activeDbKey}
  {#key activeDbKey}
    <DatabaseView
      host={$hostTree[activeHostKey]}
      database={$hostTree[activeHostKey]?.databases[activeDbKey]}
      bind:tab={dbTab}
    />
  {/key}
{:else if activeHostKey}
  {#key activeHostKey}
    <HostView host={$hostTree[activeHostKey]} bind:tab={hostTab} />
  {/key}
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
