<script>
  import CodeViewer from '../../../components/codeviewer.svelte';
  import ObjectGrid from '../../../components/objectgrid.svelte';
  import { DropIndex, GetIndexes } from '../../../../wailsjs/go/app/App';

  export let collection;

  let indexes = [];
  let activeKey = '';
  let json = '';

  async function getIndexes() {
    const result = await GetIndexes(collection.hostKey, collection.dbKey, collection.key);
    if (result) {
      indexes = result;
    }
  }

  async function dropActive() {
    if (!activeKey) {
      return;
    }
    const success = await DropIndex(collection.hostKey, collection.dbKey, collection.key, activeKey);
    if (success) {
      await getIndexes();
    }
  }

  function openJson(indexId) {
    const item = indexes?.filter(i => i.name == indexId);
    json = JSON.stringify(item, undefined, 2);
  }
</script>

<div class="indexes">
  <div class="actions">
    <button class="btn" on:click={getIndexes}>Get indexes</button>
    <button class="btn danger" on:click={dropActive} disabled={!indexes?.length || !activeKey}>
      Drop selected
    </button>
    <button class="btn">Create&hellip;</button>
  </div>

  <div class="grid">
    <ObjectGrid key="name" data={indexes} bind:activeKey on:trigger={e => openJson(e.detail)} />
  </div>
</div>

<CodeViewer bind:code={json} language="json" />

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
