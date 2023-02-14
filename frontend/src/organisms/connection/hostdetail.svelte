<script>
  import { input } from '../../lib/actions';
  import { createEventDispatcher } from 'svelte';
  import { AddHost, UpdateHost } from '../../../wailsjs/go/app/App';
  import Modal from '../../components/modal.svelte';

  export let show = false;
  export let hostKey = '';
  export let hosts = {};

  const dispatch = createEventDispatcher();
  let form = {};
  let error = '';
  $: valid = validate(form);
  $: host = hosts[hostKey];

  $: if (show || !show) {
    init();
  }

  function init() {
    form = { ...(host || {}) };
  }

  function validate(form) {
    return form.name && form.uri && true;
  }

  async function submit() {
    if (!valid) {
      return;
    }

    try {
      if (host && hostKey) {
        await UpdateHost(hostKey, JSON.stringify(form));
      }
      else {
        await AddHost(JSON.stringify(form));
      }
      show = false;
      dispatch('reload');
    }
    catch (e) {
      error = e;
    }
  }
</script>

<Modal bind:show title={host ? `Edit ${host.name}` : 'Create a new host'}>
  <form on:submit|preventDefault={submit}>
    <label class="field">
      <span class="label">Label</span>
      <input type="text" placeholder="mywebsite.com MongoDB" bind:value={form.name} use:input={{ autofocus: true }} />
    </label>

    <label class="field">
      <span class="label">Connection string</span>
      <input type="text" placeholder="mongodb://..." bind:value={form.uri} spellcheck="false" use:input />
    </label>

    <div class="result">
      <div>
        {#if error}
          <div class="error">{error}</div>
        {/if}
      </div>
      <button class="btn" disabled={!valid} type="submit">
        {host ? 'Save' : 'Create'}
      </button>
    </div>
  </form>
</Modal>

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .result {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
