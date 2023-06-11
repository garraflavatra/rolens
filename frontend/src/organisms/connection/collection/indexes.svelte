<script>
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/objectgrid.svelte';
  import { DropIndex, GetIndexes } from '$wails/go/app/App';
  import IndexDetail from './dialogs/indexdetail.svelte';

  export let collection;

  let indexes = [];
  let activePath = [];
  let creatingNewIndex = false;

  $: collection && getIndexes();

  async function getIndexes() {
    const result = await GetIndexes(collection.hostKey, collection.dbKey, collection.key);
    if (result) {
      indexes = result;
    }
  }

  function createIndex() {
    creatingNewIndex = true;
  }

  async function drop(key) {
    if (typeof key !== 'string') {
      key = activePath[0];
    }
    const success = await DropIndex(collection.hostKey, collection.dbKey, collection.key, key);
    if (success) {
      await getIndexes();
      activePath[0] = '';
    }
  }
</script>

<div class="indexes">
  <div class="grid">
    <ObjectGrid
      key="name"
      data={indexes}
      getRootMenu={(_, idx) => [ { label: 'Drop this index', fn: () => drop(idx.name) } ]}
      bind:activePath
    />
  </div>

  <div class="actions">
    <button class="btn" on:click={getIndexes}>
      <Icon name="reload" /> Reload
    </button>
    <button class="btn" on:click={createIndex}>
      <Icon name="+" /> Create index…
    </button>
    <button class="btn danger" on:click={drop} disabled={!indexes?.length || !activePath[0]}>
      <Icon name="x" /> Drop selected
    </button>
  </div>
</div>

<IndexDetail bind:creatingNewIndex {collection} on:reload={getIndexes} />

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
