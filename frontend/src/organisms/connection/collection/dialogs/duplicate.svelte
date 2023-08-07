<script>
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import input from '$lib/actions/input.js';
  import hostTree from '$lib/stores/hosttree.js';
  import { createEventDispatcher } from 'svelte';

  export let host = {};
  export let dbKey = '';
  export let collKey = '';

  const dispatch = createEventDispatcher();
  let newHost = host.key;
  let newDb = dbKey;
  let newColl = `${collKey}-duplicate`;

  function duplicate() {
    dispatch('duplicate', { newHost, newDb, newColl });
  }
</script>

<Modal title="Duplicate collection" width="500px" on:close>
  <div class="duplicate">
    <div class="origin">
      <div class="field">
        <span class="label">{host.name || host.uri}</span>
        <!-- <input type="text" readonly value={host.name || host.uri} /> -->
      </div>
      <div class="field">
        <span class="label">{dbKey}</span>
        <!-- <input type="text" readonly value={dbKey} /> -->
      </div>
      <div class="field">
        <span class="label">{collKey}</span>
        <!-- <input type="text" readonly value={collKey} /> -->
      </div>
    </div>

    <div class="arrow">
      <Icon name="arr-r" />
    </div>

    <div class="destination">
      <label class="field">
        <span class="label">Host</span>
        <select bind:value={newHost}>
          {#each Object.values($hostTree) as { key, name }}
            <option value={key} selected={key === host.key}>{name}</option>
          {/each}
        </select>
      </label>
      <label class="field">
        <span class="label">Database</span>
        <input type="text" bind:value={newDb} use:input />
      </label>
      <label class="field">
        <span class="label">Collection</span>
        <input type="text" bind:value={newColl} use:input={{ autofocus: true }} />
      </label>
    </div>
  </div>

  <svelte:fragment slot="footer">
    <button class="button" on:click={duplicate}>
      <Icon name="play" /> Duplicate
    </button>
  </svelte:fragment>
</Modal>

<style>
  .duplicate {
    display: grid;
    grid-template: auto / 1fr auto 1fr;
    gap: 0.5rem;
  }

  .arrow {
    align-self: center;
  }

  .field:not(:last-child) {
    margin-bottom: 0.5rem;
  }

  .origin .field .label {
    width: 100%;
    display: inline-block;
    background-color: #fff;
  }
</style>
