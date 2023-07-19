<script>
  import BlankState from '$components/blankstate.svelte';
  import TabBar from '$components/tabbar.svelte';
  import { EventsOn } from '$wails/runtime/runtime';

  import Shell from '../shell.svelte';
  import Stats from './stats.svelte';

  export let host;
  export let database;
  export let tab = 'stats';

  const tabs = {
    'stats': { icon: 'chart', title: 'Database stats', component: Stats },
    'shell': { icon: 'shell', title: 'Shell', component: Shell },
  };

  for (const key of Object.keys(tabs)) {
    tabs[key].key = key;
  }

  EventsOn('OpenStatsTab', name => (tab = name || tab));
</script>

<div class="view" class:empty={!database}>
  {#if database}
    {#key database}
      <TabBar tabs={Object.values(tabs)} bind:selectedKey={tab} />

      {#each Object.values(tabs) as view}
        <div class="container" class:hidden={tab !== view.key}>
          <svelte:component this={view.component} visible={tab === view.key} {host} {database} />
        </div>
      {/each}
    {/key}
  {:else}
    <BlankState label="Select a database to continue" />
  {/if}
</div>

<style>
  .view {
    height: 100%;
    display: grid;
    grid-template: auto 1fr / 1fr;
  }
  .view.empty {
    grid-template: 1fr / 1fr;
  }

  .container {
    padding: 0.5rem;
    display: flex;
    align-items: stretch;
    overflow: auto;
    min-height: 0;
    min-width: 0;
  }
  .container.hidden {
    display: none;
  }
  .container > :global(*) {
    width: 100%;
  }
</style>
