<script>
  import { tick } from 'svelte';
  import TabBar from '../../components/tabbar.svelte';
  import Find from './find.svelte';
  import Indexes from './indexes.svelte';
  import Insert from './insert.svelte';
  import Remove from './remove.svelte';
  import Stats from './stats.svelte';

  export let collection;
  export let hostKey;
  export let dbKey;
  export let collectionKey;

  let tab = 'find';
  let find;

  $: if (collection) {
    collection.hostKey = hostKey;
    collection.dbKey = dbKey;
    collection.key = collectionKey;
  }

  async function catchQuery(event) {
    tab = 'find';
    await tick();
    find.performQuery(event.detail);
  }
</script>

<div class="collection" class:empty={!collection}>
  {#if collection}
    {#key collection}
      <TabBar tabs={[
        { key: 'stats', title: 'Stats' },
        { key: 'find', title: 'Find' },
        { key: 'insert', title: 'Insert' },
        { key: 'update', title: 'Update' },
        { key: 'remove', title: 'Remove' },
        { key: 'indexes', title: 'Indexes' },
      ]} bind:selectedKey={tab} />

      <div class="container">
        {#if tab === 'stats'} <Stats {collection} />
        {:else if tab === 'find'} <Find {collection} bind:this={find} />
        {:else if tab === 'insert'} <Insert {collection} on:performFind={catchQuery} />
        {:else if tab === 'remove'} <Remove {collection} />
        {:else if tab === 'indexes'} <Indexes {collection} />
        {/if}
      </div>
    {/key}
  {:else}
    No collection selected
  {/if}
</div>

<style>
  .collection {
    height: 100%;
    display: grid;
    grid-template: auto 1fr / 1fr;
    gap: 0.5rem;
  }

  .container {
    padding: 0 0.5rem;
    display: flex;
    align-items: stretch;
  }
  .container > :global(*) {
    width: 100%;
  }
</style>
