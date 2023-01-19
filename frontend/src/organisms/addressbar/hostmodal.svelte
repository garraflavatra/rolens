<script>
  import Modal from '../../components/modal.svelte';
  import { createEventDispatcher, onMount } from 'svelte';
  import Icon from '../../components/icon.svelte';
  import { Hosts, RemoveHost } from '../../../wailsjs/go/app/App';
  import Welcome from './welcome.svelte';
  import CreateHostModal from './createhostmodal.svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let modalOpen = false;

  const dispatch = createEventDispatcher();
  let error = '';
  let createHostModalOpen = false;
  $: host = hosts?.[activeHostKey];
  $: hostCount = Object.keys(hosts).length;

  $: if (!modalOpen) {
    error = '';
  }

  function select(hostKey) {
    activeHostKey = hostKey;
    dispatch('select', hostKey);
  }

  async function getHosts() {
    try {
      const h = await Hosts();
      hosts = h || {};
    }
    catch (e) {
      error = e;
    }
  }

  async function removeHost(hostKey) {
    try {
      await RemoveHost(hostKey);
      await getHosts();
    }
    catch (e) {
      error = e;
    }
  }

  onMount(getHosts);
</script>

<Modal bind:show={modalOpen} title={hostCount && 'Hosts'} width="60vw">
  <div class="status">
    <p class:error>
      {#if error}
        <strong>Oops!</strong> {error}
      {:else}
        {hostCount} host{hostCount === 1 ? '' : 's'}
      {/if}
    </p>
    <button class="btn" on:click={() => createHostModalOpen = true}>
      Create new host
    </button>
  </div>
  {#if hostCount}
    <ul class="hosts">
      {#each Object.entries(hosts) as [hostKey, host]}
        <li>
          <div class="host">
            <button class="btn secondary" on:click={() => select(hostKey)} title="Connect to {host.name}">
              {host.name}
            </button>
            <button class="btn secondary" title="Edit {host.name}">
              <Icon name="edit" />
            </button>
            <button class="btn secondary" title="Remove {host.name}" on:click={() => removeHost(hostKey)}>
              <Icon name="x" />
            </button>
          </div>
        </li>
      {/each}
    </ul>
  {:else}
    <Welcome on:createHost={() => createHostModalOpen = true} />
  {/if}
</Modal>

<CreateHostModal bind:show={createHostModalOpen} on:reload={getHosts} />

<style>
  .hosts {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0.5rem;
  }

  .status {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .error {
    color: #c00;
  }

  .host {
    display: grid;
    grid-template: 1fr 1fr / 1fr auto;
    align-items: center;
    width: 100%;
    height: 100%;
  }
  .host button {
    border-radius: 0;
    border-left: 1px solid #ccc;
  }
  .host button:nth-child(1) {
    border-right: none;
    grid-row: 1 / 3;
    height: 100%;
    border-radius: 10px 0 0 10px;
  }
  .host button:nth-child(2) {
    border-top-right-radius: 10px;
    border-bottom: none;
  }
  .host button:nth-child(3) {
    border-bottom-right-radius: 10px;
  }
</style>
