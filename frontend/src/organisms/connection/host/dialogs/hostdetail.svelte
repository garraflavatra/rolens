<script>
  import Modal from '$components/modal.svelte';
  import input from '$lib/actions/input.js';
  import hostTree from '$lib/stores/hosttree.js';
  import { AddHost, UpdateHost } from '$wails/go/app/App.js';
  import { createEventDispatcher, onMount } from 'svelte';

  export let hostKey = '';

  const dispatch = createEventDispatcher();
  let form = {};
  let error = '';
  $: valid = validate(form);
  $: host = $hostTree[hostKey];

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
        const newHostKey = await AddHost(JSON.stringify(form));
        if (newHostKey) {
          hostKey = newHostKey;
        }
      }
      dispatch('updated', form);
      dispatch('close');
    }
    catch (e) {
      error = e;
    }
  }

  onMount(() => {
    form = { ...(host || {}) };
  });
</script>

<Modal title={host ? `Edit ${host.name}` : 'Create a new host'} on:close>
  <form on:submit|preventDefault={submit}>
    <label class="field">
      <span class="label">Label</span>
      <input type="text" placeholder="mywebsite.com MongoDB" bind:value={form.name} use:input={{ autofocus: true }} />
    </label>

    <label class="field">
      <span class="label">Connection string</span>
      <input
        type="text"
        placeholder="mongodb://..."
        bind:value={form.uri}
        spellcheck="false"
        use:input
      />
    </label>
  </form>

  <div class="result" slot="footer">
    <div>
      {#if error}
        <div class="error">{error}</div>
      {/if}
    </div>
    <button class="button" disabled={!valid} on:click={submit}>
      {host ? 'Save' : 'Create'}
    </button>
  </div>
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
