<script>
  import Modal from '$components/modal.svelte';
  import { createEventDispatcher, onMount, tick } from 'svelte';

  export let title = '';
  export let description = '';
  export let value = '';

  const dispatch = createEventDispatcher();
  let input;

  function submit() {
    dispatch('submit', { value });
  }

  function close() {
    dispatch('close');
  }

  onMount(() => tick().then(() => input.select()));
</script>

<Modal {title} on:close width="350px">
  {#if description}
    <p>{description}</p>
  {/if}

  <form on:submit|preventDefault={submit}>
    <label class="field">
      <input type="text" bind:value bind:this={input} spellcheck="false" />
    </label>
  </form>

  <svelte:fragment slot="footer">
    <button on:click={submit} class="btn">OK</button>
    <button on:click={close} class="btn secondary">Cancel</button>
  </svelte:fragment>
</Modal>

<style>
  p {
    line-height: 1.25;
  }
</style>
