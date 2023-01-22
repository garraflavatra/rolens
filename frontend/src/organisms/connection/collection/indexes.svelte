<script>
  import ObjectViewer from '../../../components/objectviewer.svelte';
  import ObjectGrid from '../../../components/objectgrid.svelte';
  import { DropIndex, GetIndexes } from '../../../../wailsjs/go/app/App';
  import Icon from '../../../components/icon.svelte';
  import IndexDetail from './indexes-detail.svelte';

  export let collection;

  let indexes = [];
  let activePath = [];
  let objectViewerData = '';
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

  function openJson(indexId) {
    const item = indexes?.find(i => i.name == indexId);
    objectViewerData = item;
  }
</script>

<div class="indexes">
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

  <div class="grid">
    <ObjectGrid
      key="name"
      data={indexes}
      getRootMenu={(_, idx) => [ { label: 'Drop this index', fn: () => drop(idx.name) } ]}
      bind:activePath
      on:trigger={e => openJson(e.detail.itemKey)}
    />
  </div>
</div>

<ObjectViewer bind:data={objectViewerData} />
<IndexDetail bind:creatingNewIndex {collection} on:reload={getIndexes} />

<style>
  .indexes {
    display: grid;
    gap: 0.5rem;
    grid-template: auto 1fr / 1fr;
  }

  .indexes .grid {
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
  }
</style>
