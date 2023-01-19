<script>
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';
  import ObjectTree from './objecttree.svelte';
  import { onDestroy } from 'svelte';

  export let data;

  let copySucceeded = false;
  let timeout;

  async function copy() {
    await navigator.clipboard.writeText(JSON.stringify(data));
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  onDestroy(() => clearTimeout(timeout));
</script>

{#if data}
  <Modal bind:show={data} title="Object viewer" contentPadding={false}>
    <div class="objectviewer">
      <div class="buttons">
        <button class="btn" on:click={copy}>
          <Icon name={copySucceeded ? 'check' : 'clipboard'} />
        </button>
      </div>
      <div class="code">
        <ObjectTree {data} />
      </div>
    </div>
  </Modal>
{/if}

<style>
  .objectviewer {
    position: relative;
  }
  .code {
    padding: 1rem;
  }
  .buttons {
    position: absolute;
    top: 0;
    right: 0;
    margin: 1rem;
  }
  .buttons button {
    margin-left: 1rem;
  }
</style>
