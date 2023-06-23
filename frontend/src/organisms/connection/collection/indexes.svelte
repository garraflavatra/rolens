<script>
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/objectgrid.svelte';
  import { onMount } from 'svelte';

  export let collection;

  let activePath = [];
  let _indexes = [];
  let error = '';

  async function refresh() {
    error = await collection.getIndexes();
    if (!error) {
      _indexes = collection.indexes.map(idx => {
        return {
          name: idx.name,
          background: idx.background || false,
          unique: idx.unique || false,
          sparse: idx.sparse || false,
          model: idx.model,
        };
      });
    }
  }

  async function createIndex() {
    const newIndexName = await collection.newIndex();
    if (newIndexName) {
      await refresh();
    }
  }

  async function dropIndex(indexName) {
    if (typeof indexName !== 'string') {
      indexName = activePath[0];
    }

    const success = await collection.getIndexByName(indexName).drop();

    if (success) {
      activePath[0] = '';
      await refresh();
    }
  }

  onMount(refresh);
</script>

<div class="indexes">
  <div class="grid">
    <ObjectGrid
      key="name"
      data={_indexes}
      getRootMenu={(_, idx) => [ { label: 'Drop this index', fn: () => dropIndex(idx.name) } ]}
      errorTitle={error ? 'Error while getting indexes' : ''}
      errorDescription={error}
      bind:activePath
    />
  </div>

  <div class="actions">
    <button class="btn" on:click={refresh}>
      <Icon name="reload" /> Reload
    </button>
    <button class="btn" on:click={createIndex}>
      <Icon name="+" /> Create index…
    </button>
    <button class="btn danger" on:click={dropIndex} disabled={!_indexes.length || !activePath[0]}>
      <Icon name="x" /> Drop selected
    </button>
  </div>
</div>

<style>
  .indexes {
    display: grid;
    gap: 0.5rem;
    grid-template: 1fr auto / 1fr;
  }

  .indexes .grid {
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
  }
</style>
