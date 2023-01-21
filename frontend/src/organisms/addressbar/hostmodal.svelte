<script>
  import Modal from '../../components/modal.svelte';
  import { createEventDispatcher, onMount } from 'svelte';
  import Icon from '../../components/icon.svelte';
  import { Hosts, RemoveHost } from '../../../wailsjs/go/app/App';
  import Welcome from './welcome.svelte';
  import HostDetail from './hostdetail.svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let modalOpen = false;

  const dispatch = createEventDispatcher();
  let error = '';
  let hostDetailModalOpen = false;
  let hostDetailModalHost;
  let hostDetailModalKey = '';
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

  function createHost() {
    hostDetailModalHost = undefined;
    hostDetailModalKey = '';
    hostDetailModalOpen = true;
  }

  function editHost(hostKey) {
    hostDetailModalHost = hosts[hostKey];
    hostDetailModalKey = hostKey;
    hostDetailModalOpen = true;
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

<Modal bind:show={modalOpen} title={hostCount && 'Hosts'} width="60vw" overflow={false}>
  {#if hostCount}
    <div class="status">
      <p class:error>
        {#if error}
          <strong>Oops!</strong> {error}
        {:else}
          {hostCount} host{hostCount === 1 ? '' : 's'}
        {/if}
      </p>
      <button class="btn" on:click={createHost}>
        Create new host
      </button>
    </div>
    <ul class="hosts">
      {#each Object.entries(hosts) as [hostKey, host]}
        <li>
          <div class="host">
            <button class="btn secondary" title="Connect to {host.name}" on:click={() => select(hostKey)}>
              {host.name}
            </button>
            <button class="btn secondary" title="Edit {host.name}" on:click={() => editHost(hostKey)}>
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
    <Welcome on:createHost={createHost} />
  {/if}
</Modal>

<HostDetail
  host={hostDetailModalHost}
  hostKey={hostDetailModalKey}
  on:reload={getHosts}
  bind:show={hostDetailModalOpen}
/>

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
