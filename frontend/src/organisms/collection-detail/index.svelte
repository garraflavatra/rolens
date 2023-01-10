<script>
  import ObjectGrid from '../../components/objectgrid.svelte';
  import CodeExample from '../../components/code-example.svelte';
  import TabBar from '../../components/tabbar.svelte';
  import Find from './find.svelte';
  import Indexes from './indexes.svelte';
  import Insert from './insert.svelte';
  import Remove from './remove.svelte';

  export let collection;
  export let hostKey;
  export let dbKey;
  export let collectionKey;

  let tab = 'find';

  $: if (collection) {
    collection.hostKey = hostKey;
    collection.dbKey = dbKey;
    collection.key = collectionKey;
  }
</script>

<div class="collection" class:empty={!collection}>
  {#if collection}
    <TabBar tabs={[
      { key: 'stats', title: 'Stats' },
      { key: 'find', title: 'Find' },
      { key: 'insert', title: 'Insert' },
      { key: 'update', title: 'Update' },
      { key: 'remove', title: 'Remove' },
      { key: 'indexes', title: 'Indexes' },
    ]} bind:selectedKey={tab} />

    <div class="container">
      {#if tab === 'stats'}
        <CodeExample code="db.stats()" />
        <ObjectGrid data={collection.stats} />
      {:else if tab === 'find'}
        <Find {collection} />
      {:else if tab === 'insert'}
        <Insert {collection} />
      {:else if tab === 'remove'}
        <Remove {collection} />
      {:else if tab === 'indexes'}
        <Indexes {collection} />
      {/if}
    </div>
  {:else}
    No collection selected
  {/if}
</div>

<style>
  .collection {
    margin: 1rem 1rem 1rem 0;
    height: 100%;
    display: flex;
    flex-flow: column;
  }

  .container {
    padding: 0.5rem 0.5rem 0;
    flex: 1;
    display: flex;
    flex-flow: column;
  }
</style>
