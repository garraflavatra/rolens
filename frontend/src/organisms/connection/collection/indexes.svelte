<script>
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/objectgrid.svelte';

  export let collection;
  export let visible = false;

  let activePath = [];
  let _indexes = [];
  let error = '';
  let busy = false;

  async function refresh() {
    if (!visible) {
      return;
    }

    busy = 'Fetching indexes…';
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

    busy = false;
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

  $: visible && refresh();
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
      {busy}
    />
  </div>

  <div class="actions">
    <button class="button" on:click={refresh}>
      <Icon name="reload" spin={busy} /> Reload
    </button>
    <button class="button" on:click={createIndex}>
      <Icon name="+" /> Create index…
    </button>
    <button class="button danger" on:click={dropIndex} disabled={!_indexes.length || !activePath[0]}>
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
