<script>
  import { looseJsonIsValid } from '$lib/strings';
  import { EJSON } from 'bson';
  import { createEventDispatcher, onDestroy } from 'svelte';
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';
  import ObjectEditor from './objecteditor.svelte';

  export let data;
  export let readonly = false;
  export let saveable = false;
  export let successMessage = '';

  const dispatch = createEventDispatcher();
  let copySucceeded = false;
  let timeout;
  let text = EJSON.stringify(data, undefined, '\t');
  $: invalid = !looseJsonIsValid(text);

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
    dispatch('save', { text, originalData: data });
  }

  onDestroy(() => clearTimeout(timeout));
</script>

{#if data}
  <Modal bind:show={data} contentPadding={false}>
    <div class="objectviewer">
      <ObjectEditor bind:text on:updated={() => successMessage = ''} {readonly} />
    </div>

    <svelte:fragment slot="footer">
      {#if saveable}
        <button class="button" on:click={save} disabled={invalid}>
          <Icon name="save" /> Save
        </button>
      {/if}

      <button class="button secondary" on:click={close}>
        <Icon name="x" /> Close
      </button>

      <button class="button secondary" on:click={copy}>
        <Icon name={copySucceeded ? 'check' : 'clipboard'} /> Copy
      </button>

      {#if successMessage}
        <span class="flash-green">{successMessage}</span>
      {/if}
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

  .flash-green {
    margin-left: 0.5rem;
  }
</style>
