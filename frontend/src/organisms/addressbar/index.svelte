<script>
  import Modal from '../../components/modal.svelte';
  import Icon from '../../components/icon.svelte';
  import { createEventDispatcher } from 'svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let modalOpen = false;

  const dispatch = createEventDispatcher();
  $: host = hosts?.[activeHostKey];

  function select(hostKey) {
    activeHostKey = hostKey;
    dispatch('select', hostKey);
  }
</script>

<div class="addressbar">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="address" class:empty={!host?.uri} on:click={() => modalOpen = true}>
    {host?.uri || 'No host selected'}
  </div>

  <div class="actions">
    <button class="btn" on:click={() => modalOpen = true}>
      <Icon name="db" />
    </button>
  </div>
</div>

<Modal bind:show={modalOpen} title="Hosts">
  {#if Object.keys(hosts).length}
    <ul class="hosts">
      {#each Object.entries(hosts) as [hostKey, host]}
        <li>
          <button on:click={() => select(hostKey)}>
            {host.name}
          </button>
        </li>
      {/each}
    </ul>
  {/if}
</Modal>

<style>
  .addressbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 1rem;
    padding: 0.5rem 0.5rem 0.5rem 1rem;
    border: 1px solid #ccc;
    border-radius: 10px;
  }

  .address.empty {
    font-style: italic;
    opacity: 0.6;
  }

  .hosts {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0.5rem;
  }
  .hosts li button {
    display: block;
    width: 100%;
    height: 100%;
    padding: 1rem;
    background-color: #ddd;
    border-radius: 10px;
  }
  .hosts li button:hover {
    background-color: #ccc;
  }
</style>
