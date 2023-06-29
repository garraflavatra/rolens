<script>
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import views from '$lib/stores/views';
  import { createEventDispatcher } from 'svelte';

  export let collection;
  export let query = {};

  const dispatch = createEventDispatcher();
  const exportInfo = { ...query, viewKey: collection.viewKey };

  function submit() {
    exportInfo.view = $views[exportInfo.viewKey];
    dispatch('export', { exportInfo });
  }
</script>

<Modal title="Export results" width="450px" on:close>
  <form on:submit|preventDefault={submit}>
    <label class="field">
      <span class="label">Export</span>
      <select bind:value={exportInfo.contents}>
        <option value="all">all records</option>
        <option value="query" disabled={!query}>all records matching query</option>
        <option value="querylimitskip" disabled={!query}>all records matching query, considering limit and skip</option>
      </select>
    </label>

    <label class="field">
      <span class="label">Format</span>
      <select bind:value={exportInfo.format}>
        <option value="jsonarray">JSON array (*.json)</option>
        <option value="ndjson">Newline delimited JSON (*.ndjson)</option>
        <option value="csv">CSV (*.csv)</option>
        <option value="excel">Excel (*.xlsx)</option>
      </select>
    </label>

    <label class="field">
      <span class="label">View to use</span>
      <select bind:value={exportInfo.viewKey}>
        {#each Object.entries(views.forCollection(collection.hostKey, collection.dbKey, collection.key)) as [ key, { name } ]}
          <option value={key}>{name}</option>
        {/each}
      </select>
      <button class="button" type="button" on:click={() => dispatch('openViewConfig')} title="Edit view">
        <Icon name="cog" />
      </button>
    </label>
  </form>

  <svelte:fragment slot="footer">
    <button class="button" on:click={submit}>
      <Icon name="play" /> Start export
    </button>
  </svelte:fragment>
</Modal>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
</style>
