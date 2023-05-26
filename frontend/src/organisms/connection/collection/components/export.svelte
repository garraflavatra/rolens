<script>
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import { startProgress } from '$lib/progress';
  import views from '$lib/stores/views';
  import { createEventDispatcher } from 'svelte';

  export let info;
  export let collection;

  const dispatch = createEventDispatcher();
  let viewKey = collection.viewKey;
  $: viewKey = collection.viewKey;

  async function performExport() {
    const progress = startProgress('Performing export…');
    info.view = $views[viewKey];
    //...
    progress.end();
  }
</script>

<Modal bind:show={info} title="Export results" width="400px">
  <form on:submit|preventDefault={performExport}>
    <label class="field">
      <span class="label">Export</span>
      <select bind:value={info.contents}>
        <option value="all">all records</option>
        <option value="query">all records matching query</option>
        <option value="querylimitskip">all records matching query, considering limit and skip</option>
      </select>
    </label>
    <label class="field">
      <span class="label">Format</span>
      <select bind:value={info.format}>
        <option value="jsonarray">JSON array</option>
        <option value="jsonnewline">JSON: newline-separated objects</option>
        <option value="csv">CSV</option>
      </select>
    </label>
    <label class="field">
      <span class="label">View to use</span>
      <select bind:value={viewKey}>
        {#each Object.entries(views.forCollection(collection.hostKey, collection.dbKey, collection.key)) as [key, { name }]}
          <option value={key}>{name}</option>
        {/each}
      </select>
      <button class="btn" type="button" on:click={() => dispatch('openViewConfig')} title="Edit view">
        <Icon name="cog" />
      </button>
    </label>
  </form>
</Modal>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
</style>
