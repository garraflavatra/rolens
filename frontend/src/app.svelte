<script>
  import BlankState from '$components/blankstate.svelte';
  import ContextMenu from '$components/contextmenu.svelte';
  import dialogs from '$lib/dialogs.js';
  import contextMenu from '$lib/stores/contextmenu.js';
  import hostTree from '$lib/stores/hosttree.js';
  import applicationInited from '$lib/stores/inited.js';
  import windowTitle from '$lib/stores/windowtitle.js';
  import Connection from '$organisms/connection/index.svelte';
  import { EventsOn } from '$wails/runtime/runtime.js';
  import { tick } from 'svelte';
  import AboutDialog from './dialogs/about.svelte';
  import SettingsDialog from './dialogs/settings/index.svelte';

  let showWelcomeScreen = undefined;

  applicationInited.defer(() => {
    hostTree.subscribe(hosts => {
      if (hostTree.hasBeenInited() && (showWelcomeScreen === undefined)) {
        showWelcomeScreen = !Object.keys(hosts || {}).length;
      }
    });
  });

  async function createFirstHost() {
    showWelcomeScreen = false;
    await tick();
    hostTree.newHost();
  }

  function showAboutDialog() {
    dialogs.new(AboutDialog);
  }

  function showSettings() {
    dialogs.new(SettingsDialog);
  }

  EventsOn('OpenPreferences', showSettings);
  EventsOn('OpenAboutModal', showAboutDialog);
</script>

<svelte:window on:contextmenu|preventDefault />

<div id="root">
  <div class="titlebar">{$windowTitle}</div>

  {#if $applicationInited && (showWelcomeScreen !== undefined)}
    <main class:empty={showWelcomeScreen}>
      {#if showWelcomeScreen}
        <BlankState label="Welcome to Rolens!" image="/logo.png" pale={false} big={true}>
          <button class="button" on:click={createFirstHost}>Add your first host</button>
        </BlankState>
      {:else}
        <Connection />
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
    color: #fff;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
  }
  :global([data-platform="darwin"]) .titlebar {
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
  :global([data-platform="darwin"]) main {
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
</style>
