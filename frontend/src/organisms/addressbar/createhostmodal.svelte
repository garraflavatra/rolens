<script>
  import { createEventDispatcher } from 'svelte';
  import { AddHost } from '../../../wailsjs/go/app/App';
  import Modal from '../../components/modal.svelte';

  export let show = false;

  const dispatch = createEventDispatcher('reload');
  let form = {};
  let error = '';
  $: valid = validate(form);

  $: if (show || !show) {
    form = {};
  }

  function validate(form) {
    return form.name && form.uri && true;
  }

  async function submit() {
    if (!valid) {
      return;
    }

    try {
      await AddHost(JSON.stringify(form));
      show = false;
      dispatch('reload');
    }
    catch (e) {
      error = e;
    }
  }
</script>

<Modal bind:show title="Create a new host">
  <form on:submit|preventDefault={submit}>
    <label class="field">
      <span class="label">Label</span>
      <input type="text" placeholder="mywebsite.com MongoDB" bind:value={form.name} />
    </label>

    <label class="field">
      <span class="label">Connection string</span>
      <input type="text" placeholder="mongodb://..." bind:value={form.uri} spellcheck="false" />
    </label>

    <div class="result">
      <div>
        {#if error}
          <div class="error">{error}</div>
        {/if}
      </div>
      <button class="btn" disabled={!valid} type="submit">Create</button>
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
