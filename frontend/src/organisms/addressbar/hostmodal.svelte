<script>
  import Modal from '../../components/modal.svelte';
  import { createEventDispatcher } from 'svelte';
  import Icon from '../../components/icon.svelte';

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

<Modal bind:show={modalOpen} title="Hosts">
  {#if Object.keys(hosts).length}
    <ul class="hosts">
      {#each Object.entries(hosts) as [hostKey, host]}
        <li>
          <div class="host">
            <!-- <div class="name">{host.name}</div> -->
            <button class="btn" on:click={() => select(hostKey)} title="Connect to {host.name}">
              {host.name}
            </button>
            <button class="btn" title="Edit {host.name}">
              <Icon name="edit" />
            </button>
            <button class="btn" title="Remove {host.name}">
              <Icon name="x" />
            </button>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</Modal>

<style>
  .hosts {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0.5rem;
  }

  .host {
    display: grid;
    grid-template: 1fr 1fr / 1fr auto;
    align-items: center;
    width: 100%;
    height: 100%;
    border-radius: 10px;
    overflow: hidden;
    border: 1px solid #ccc;
  }
  .host button {
    border-radius: 0;
    background: #eee;
    color: inherit;
    border: none;
    border-left: 1px solid #ccc;
  }
  .host button:hover {
    background-color: #ddd;
  }
  .host button:first-child {
    grid-row: 1 / 3;
    height: 100%;
    border: none;
    background-color: #fff;
  }
  .host button:first-child:hover {
    background-color: #eee;
  }
  .host button:last-child {
    border-top: 1px solid #ccc;
  }
</style>
