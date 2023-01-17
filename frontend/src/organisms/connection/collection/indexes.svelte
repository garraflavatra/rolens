<script>
  import CodeViewer from '../../../components/codeviewer.svelte';
  import ObjectGrid from '../../../components/objectgrid.svelte';
  import { GetIndexes } from '../../../../wailsjs/go/app/App';

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

  function openJson(itemId) {
    const item = indexes?.filter(i => i.name == itemId);
    json = JSON.stringify(item, undefined, 2);
  }
</script>

<div class="indexes">
  <div class="actions">
    <button class="btn" on:click={getIndexes}>Get indexes</button>
    <button class="btn danger" disabled={!indexes?.length || !activeKey}>Drop selected</button>
    <button class="btn">Create&hellip;</button>
  </div>

  <ObjectGrid key="name" data={indexes} bind:activeKey on:trigger={e => openJson(e.detail)} />
</div>

<CodeViewer bind:code={json} language="json" />

<style>
  .indexes {
    display: grid;
    gap: 0.5rem;
    grid-template: auto 1fr / 1fr;
  }
</style>