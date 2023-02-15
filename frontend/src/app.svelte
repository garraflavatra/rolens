<script>
  import ContextMenu from '$components/contextmenu.svelte';
  import { connections } from '$lib/stores/connections';
  import contextMenu from '$lib/stores/contextmenu';
  import environment from '$lib/stores/environment';
  import applicationInited from '$lib/stores/inited';
  import About from '$organisms/about/index.svelte';
  import Connection from '$organisms/connection/index.svelte';
  import Settings from '$organisms/settings/index.svelte';
  import { EventsOn } from '$wails/runtime';

  const hosts = {};
  const activeHostKey = '';
  let activeDbKey = '';
  let activeCollKey = '';
  let settingsModalOpen = false;
  let aboutModalOpen = false;

  $: host = hosts[activeHostKey];
  $: connection = $connections[activeHostKey];
  $: database = connection?.databases[activeDbKey];
  $: collection = database?.collections?.[activeCollKey];

  EventsOn('OpenPrefrences', () => settingsModalOpen = true);
  EventsOn('OpenAboutModal', () => aboutModalOpen = true);
</script>

<svelte:window on:contextmenu|preventDefault />

<div id="root" class="platform-{$environment?.platform}">
  <div class="titlebar"></div>

  {#if $applicationInited}
    <main class:empty={!host || !connection}>
      <Connection {hosts} bind:activeCollKey bind:activeDbKey {activeHostKey} />
    </main>

    {#key $contextMenu}
      <ContextMenu {...$contextMenu} on:close={contextMenu.hide} />
    {/key}
    <Settings bind:show={settingsModalOpen} />
    <About bind:show={aboutModalOpen} />
  {/if}
</div>

<style>
  .titlebar {
    height: 0;
    background-color: #00002a;
    --wails-draggable: drag;
  }
  #root.platform-darwin .titlebar {
    height: var(--darwin-titlebar-height);
  }

  main {
    height: 100vh;
    display: grid;
    grid-template: 1fr / 250px 1fr;
  }
  #root.platform-darwin main {
    height: calc(100vh - var(--darwin-titlebar-height));
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
