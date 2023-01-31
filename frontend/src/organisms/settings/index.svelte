<script>
  import { input } from '../../actions';
  import Modal from '../../components/modal.svelte';
  import { applicationSettings as settings } from '../../stores';

  export let show = false;
</script>

<Modal title="Prefrences" bind:show>
  <div class="prefs">
    <label for="defaultLimit">Initial number of items to retrieve using one query (limit):</label>
    <label class="field">
      <input type="number" bind:value={$settings.defaultLimit} id="defaultLimit" use:input={{ autofocus: true }} />
      <span class="label">items</span>
    </label>

    <label for="defaultSort">Default sort query</label>
    <label class="field">
      <input type="text" class="code" bind:value={$settings.defaultSort} id="defaultSort" use:input={{ type: 'json' }} />
    </label>

    <label for="autosubmitQuery">Autosubmit query</label>
    <span>
      <input type="checkbox" id="autosubmitQuery" bind:checked={$settings.autosubmitQuery} />
      <label for="autosubmitQuery">Query items automatically after opening a collection</label>
    </span>
  </div>
</Modal>

<style>
  .prefs {
    display: grid;
    grid-template-columns: auto auto;
    gap: 0.5rem;
    align-items: center;
  }
</style>
