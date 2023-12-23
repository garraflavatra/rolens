<script>
  import { tick } from 'svelte';
  import { fly } from 'svelte/transition';

  import BlankState from '$components/blankstate.svelte';
  import Connection from '$organisms/connection/index.svelte';
  import ContextMenu from '$components/contextmenu.svelte';

  import contextMenu from '$lib/stores/contextmenu.js';
  import hostTree from '$lib/stores/hosttree.js';
  import applicationInited from '$lib/stores/inited.js';
  import windowTitle from '$lib/stores/windowtitle.js';

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
</script>

<svelte:window on:contextmenu|preventDefault />

<div id="root">
  <div class="titlebar">{$windowTitle}</div>

  {#if $applicationInited && (showWelcomeScreen !== undefined)}
    <main class:empty={showWelcomeScreen} in:fly={{ y: 10 }}>
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
    height: var(--titlebar-height);
    --wails-draggable: drag;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    background-color: #ddd;
  }

  main {
    height: calc(100vh - var(--titlebar-height));
    display: grid;
    grid-template: 1fr / minmax(300px, 0.3fr) 1fr;
  }
  main.empty {
    grid-template: 1fr / 1fr;
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
