<script>
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';
  import ObjectTree from './objecttree.svelte';
  import { onDestroy } from 'svelte';

  export let data;

  let copySucceeded = false;
  let timeout;
  let _data;

  $: if (data) {
    _data = JSON.parse(JSON.stringify(data));
    for (const key of Object.keys(_data)) {
      if (typeof _data[key] === 'undefined') {
        delete _data[key];
      }
    }
  }

  async function copy() {
    await navigator.clipboard.writeText(JSON.stringify(_data));
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  onDestroy(() => clearTimeout(timeout));
</script>

{#if data}
  <Modal bind:show={data} title="Object viewer">
    <div class="objectviewer">
      <div class="buttons">
        <button class="btn" on:click={copy}>
          <Icon name={copySucceeded ? 'check' : 'clipboard'} />
        </button>
      </div>
      <div class="code">
        <ObjectTree data={_data} />
      </div>
    </div>
  </Modal>
{/if}

<style>
  .objectviewer {
    position: relative;
  }
  .buttons {
    position: absolute;
    top: 0;
    right: 0;
  }
  .buttons button {
    margin-left: 1rem;
  }
</style>
