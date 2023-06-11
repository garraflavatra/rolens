<script>
  import BlankState from '$components/blankstate.svelte';
  import TabBar from '$components/tabbar.svelte';
  import { EventsOn } from '$wails/runtime/runtime';
  import { tick } from 'svelte';
  import Aggregate from './aggregate.svelte';
  import ViewConfig from './dialogs/viewconfig.svelte';
  import Find from './find.svelte';
  import Indexes from './indexes.svelte';
  import Insert from './insert.svelte';
  import Remove from './remove.svelte';
  import Stats from './stats.svelte';
  import Update from './update.svelte';

  export let collection;
  export let hostKey;
  export let dbKey;
  export let collKey;

  let tab = 'find';
  let find;
  let viewConfigModalOpen = false;
  let firstItem;

  $: if (collection) {
    collection.hostKey = hostKey;
    collection.dbKey = dbKey;
    collection.key = collKey;
  }

  $: if (hostKey || dbKey || collKey) {
    tab = 'find';
  }

  EventsOn('OpenCollectionTab', name => (tab = name || tab));

  async function catchQuery(event) {
    tab = 'find';
    await tick();
    find.performQuery(event.detail);
  }

  function openViewConfig(event) {
    firstItem = event.detail?.firstItem;
    viewConfigModalOpen = true;
  }
</script>

<div class="view" class:empty={!collection}>
  {#if collection}
    {#key collection}
      <TabBar tabs={[
        { key: 'stats', icon: 'chart', title: 'Stats' },
        { key: 'find', icon: 'db', title: 'Find' },
        { key: 'insert', icon: '+', title: 'Insert' },
        { key: 'update', icon: 'edit', title: 'Update' },
        { key: 'remove', icon: 'trash', title: 'Remove' },
        { key: 'indexes', icon: 'list', title: 'Indexes' },
        { key: 'aggregate', icon: 're', title: 'Aggregate' },
      ]}
        bind:selectedKey={tab} />

      <div class="container">
        {#if tab === 'stats'} <Stats {collection} />
        {:else if tab === 'find'} <Find {collection} bind:this={find} on:openViewConfig={openViewConfig} />
        {:else if tab === 'insert'} <Insert {collection} on:performFind={catchQuery} on:openViewConfig={openViewConfig} />
        {:else if tab === 'update'} <Update {collection} on:performFind={catchQuery} />
        {:else if tab === 'remove'} <Remove {collection} />
        {:else if tab === 'indexes'} <Indexes {collection} />
        {:else if tab === 'aggregate'} <Aggregate {collection} />
        {/if}
      </div>
    {/key}
  {:else}
    <BlankState label="Select a collection to continue" />
  {/if}
</div>

<style>
  .view {
    height: 100%;
    display: grid;
    grid-template: auto 1fr / 1fr;
  }
  .view.empty {
    grid-template: 1fr / 1fr;
  }

  .container {
    padding: 0.5rem;
    display: flex;
    align-items: stretch;
    overflow: auto;
    min-height: 0;
    min-width: 0;
  }
  .container > :global(*) {
    width: 100%;
  }
</style>
