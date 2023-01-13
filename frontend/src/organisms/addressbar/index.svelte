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
  <div class="address" on:click={() => modalOpen = true}>
    {#if host?.uri}
      {@const split = host.uri.split('://').map(s => s.split('/')).flat()}
      <span class="protocol">{split[0]}://</span><span class="hostname">{split[1]}</span><span class="path">{split.slice(2).join('/')}</span>
    {:else}
      <span class="empty">no host selected</span>
    {/if}
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
    padding: 0.5rem 0.5rem 0.5rem 1rem;
    height: 3rem;
    border: 1px solid #ccc;
    border-radius: 10px;
    background-color: #fff;
  }

  .address .protocol,
  .address .path,
  .address .empty {
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
