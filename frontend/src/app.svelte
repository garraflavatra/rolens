<script>
  import BlankState from '$components/blankstate.svelte';
  import ContextMenu from '$components/contextmenu.svelte';
  import contextMenu from '$lib/stores/contextmenu';
  import environment from '$lib/stores/environment';
  import hostTree from '$lib/stores/hosttree';
  import applicationInited from '$lib/stores/inited';
  import windowTitle from '$lib/stores/windowtitle';
  import About from '$organisms/about.svelte';
  import Connection from '$organisms/connection/index.svelte';
  import Settings from '$organisms/settings/index.svelte';
  import { EventsOn } from '$wails/runtime';
  import { tick } from 'svelte';

  const activeHostKey = '';
  let activeDbKey = '';
  let activeCollKey = '';
  let settingsModalOpen = false;
  let aboutModalOpen = false;
  let connectionManager;
  let showWelcomeScreen = undefined;

  hostTree.subscribe(h => {
    if (h && (showWelcomeScreen === undefined)) {
      showWelcomeScreen = !Object.keys(hostTree.get() || {}).length;
    }
  });

  async function createFirstHost() {
    showWelcomeScreen = false;
    await tick();
    hostTree.newHost();
  }

  EventsOn('OpenPreferences', () => settingsModalOpen = true);
  EventsOn('OpenAboutModal', () => aboutModalOpen = true);
</script>

<svelte:window on:contextmenu|preventDefault />

<div id="root" class="platform-{$environment?.platform}">
  <div class="titlebar">{$windowTitle}</div>

  {#if $applicationInited && (showWelcomeScreen !== undefined)}
    <main class:empty={showWelcomeScreen}>
      {#if showWelcomeScreen}
        <BlankState label="Welcome to Rolens!" image="/logo.png" pale={false} big={true}>
          <button class="btn" on:click={createFirstHost}>Add your first host</button>
        </BlankState>
      {:else}
        <Connection {activeHostKey} bind:activeCollKey bind:activeDbKey bind:this={connectionManager} />
      {/if}
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
    color: #fff;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  #root.platform-darwin .titlebar {
    height: var(--darwin-titlebar-height);
  }

  main {
    height: 100vh;
    display: grid;
    grid-template: 1fr / minmax(300px, 0.3fr) 1fr;
  }
  main.empty {
    grid-template: 1fr / 1fr;
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
