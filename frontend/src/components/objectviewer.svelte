<script>
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';
  import { createEventDispatcher, onDestroy } from 'svelte';
  import ObjectEditor from './objecteditor.svelte';
  import { jsonLooseParse } from '$lib/strings';

  export let data;
  export let saveable = false;

  const dispatch = createEventDispatcher();
  let copySucceeded = false;
  let timeout;
  let text = JSON.stringify(data, undefined, '\t');
  let newData;
  let invalid = false;

  $: {
    try {
      newData = jsonLooseParse(text);
    }
    catch {
      invalid = true;
    }
  }

  async function copy() {
    await navigator.clipboard.writeText(text);
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  function close() {
    data = undefined;
    text = '';
  }

  function save() {
    dispatch('save', text);
  }

  onDestroy(() => clearTimeout(timeout));
</script>

{#if data}
  <Modal bind:show={data} contentPadding={false}>
    <div class="objectviewer">
      <ObjectEditor {text} />
    </div>

    <svelte:fragment slot="footer">
      {#if saveable}
        <button class="btn" on:click={save} disabled={invalid}>
          <Icon name="save" /> Save
        </button>
      {/if}

      <button class="btn secondary" on:click={close}>
        <Icon name="x" /> Close
      </button>

      <button class="btn secondary" on:click={copy}>
        <Icon name={copySucceeded ? 'check' : 'clipboard'} /> Copy
      </button>
    </svelte:fragment>
  </Modal>
{/if}

<style>
  .objectviewer {
    display: flex;
    position: relative;
    justify-content: stretch;
    align-items: stretch;
    height: 100%;
  }
</style>
