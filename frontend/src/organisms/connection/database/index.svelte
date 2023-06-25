<script>
  import BlankState from '$components/blankstate.svelte';
  import TabBar from '$components/tabbar.svelte';
  import { EventsOn } from '$wails/runtime/runtime';

  import Shell from '../shell.svelte';
  import Stats from './stats.svelte';

  export let database;
  export let hostKey;
  export let dbKey;

  let tab = 'stats';

  $: if (database) {
    database.hostKey = hostKey;
    database.dbKey = dbKey;
  }

  $: if (hostKey || dbKey) {
    tab = 'stats';
  }

  EventsOn('OpenStatsTab', name => (tab = name || tab));
</script>

<div class="view" class:empty={!database}>
  {#if database}
    {#key database}
      <TabBar
        tabs={[
          { key: 'stats', icon: 'chart', title: 'Database stats' },
          { key: 'shell', icon: 'shell', title: 'Shell' },
        ]}
        bind:selectedKey={tab} />
      <div class="container">
        {#if tab === 'stats'} <Stats {database} />
        {:else if tab === 'shell'} <Shell {database} />
        {/if}
      </div>
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
  .container > :global(*) {
    width: 100%;
  }
</style>
