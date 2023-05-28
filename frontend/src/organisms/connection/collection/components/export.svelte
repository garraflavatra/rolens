<script>
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import views from '$lib/stores/views';
  import { PerformFindExport } from '$wails/go/app/App';
  import { createEventDispatcher } from 'svelte';

  export let info;
  export let collection;

  const dispatch = createEventDispatcher();
  let viewKey = collection.viewKey;
  $: viewKey = collection.viewKey;
  $: if (info) {
    info.viewKey = viewKey;
  }

  async function performExport() {
    info.view = $views[viewKey];
    const success = await PerformFindExport(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(info));

    if (success) {
      info = undefined;
    }
  }
</script>

<Modal bind:show={info} title="Export results" width="450px">
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
        <option value="ndjson">Newline delimited JSON</option>
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

    <button class="btn" type="submit">
      <Icon name="play" />
      Start export
    </button>
  </form>
</Modal>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
</style>
