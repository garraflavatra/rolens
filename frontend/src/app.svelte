<script>
  import { onMount } from 'svelte';
  import { Hosts, OpenConnection } from '../wailsjs/go/app/App';
  import BlankState from './components/blankstate.svelte';
  import ContextMenu from './components/contextmenu.svelte';
  import AddressBar from './organisms/addressbar/index.svelte';
  import Connection from './organisms/connection/index.svelte';
  import { busy, contextMenu, connections } from './stores';

  let hosts = {};
  let environment;

  let activeHostKey = '';
  let activeDbKey = '';
  let activeCollKey = '';

  let addressBarModalOpen = true;

  $: host = hosts[activeHostKey];
  $: connection = $connections[activeHostKey];
  $: database = connection?.databases[activeDbKey];
  $: collection = database?.collections?.[activeCollKey];

  async function openConnection(hostKey) {
    busy.start();
    const databases = await OpenConnection(hostKey);

    if (databases) {
      $connections[hostKey] = { databases: {} };
      databases.forEach(dbKey => {
        $connections[hostKey].databases[dbKey] = { collections: {} };
      });
      activeHostKey = hostKey;
      addressBarModalOpen = false;
      window.runtime.WindowSetTitle(`${hosts[activeHostKey].name} - Mongodup`);
    }

    busy.end();
  }

  onMount(() => {
    window.runtime.Environment().then(e => environment = e);
    Hosts().then(h => hosts = h);
  });
</script>

<svelte:window on:contextmenu|preventDefault />

<div id="app" class="platform-{environment?.platform}">
  <div class="titlebar"></div>

  {#if environment}
    <main class:empty={!host || !connection}>
      <AddressBar {hosts} bind:activeHostKey on:select={e => openConnection(e.detail)} bind:modalOpen={addressBarModalOpen} />

      {#if host && connection}
        <Connection {hosts} bind:activeCollKey bind:activeDbKey {activeHostKey} />
      {:else}
        <BlankState label="A database client is nothing without a host" image="/fish.svg" />
      {/if}
    </main>

    {#key $contextMenu}
      <ContextMenu {...$contextMenu} on:close={contextMenu.hide} />
    {/key}
  {/if}
</div>

<style>
  .titlebar {
    height: 0;
    background-color: #00002a;
    --wails-draggable: drag;
  }
  #app.platform-darwin .titlebar {
    height: var(--darwin-titlebar-height);
  }

  main {
    height: 100vh;
    display: grid;
    grid-template: 3rem auto / 250px 1fr;
    gap: 0.5rem;
    padding: 0.5rem;
  }
  #app.platform-darwin main {
    height: calc(100vh - var(--darwin-titlebar-height));
  }
  main.empty {
    grid-template: 3rem auto / 1fr;
  }

  main > :global(*) {
    overflow: auto;
    min-height: 0;
    min-width: 0;
  }
  main > :global(.addressbar) {
    grid-column: 1 / 3;
  }

  .databaselist {
    overflow: scroll;
  }

  .btn.create {
    margin-top: 0.5rem;
  }
</style>
