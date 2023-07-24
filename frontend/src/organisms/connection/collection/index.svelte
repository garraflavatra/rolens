<script>
  import BlankState from '$components/blankstate.svelte';
  import TabBar from '$components/tabbar.svelte';
  import { tick } from 'svelte';

  import Aggregate from './aggregate.svelte';
  import Find from './find.svelte';
  import Indexes from './indexes.svelte';
  import Insert from './insert.svelte';
  import Remove from './remove.svelte';
  import Shell from '../shell.svelte';
  import Stats from './stats.svelte';
  import Update from './update.svelte';

  export let host;
  export let database;
  export let collection;
  export let tab = 'find';

  const tabs = {
    'stats': { icon: 'chart', title: 'Stats', component: Stats },
    'find': { icon: 'db', title: 'Find', component: Find },
    'insert': { icon: '+', title: 'Insert', component: Insert },
    'update': { icon: 'edit', title: 'Update', component: Update },
    'remove': { icon: 'trash', title: 'Remove', component: Remove },
    'indexes': { icon: 'list', title: 'Indexes', component: Indexes },
    'aggregate': { icon: 're', title: 'Aggregate', component: Aggregate },
    'shell': { icon: 'shell', title: 'Shell', component: Shell },
  };

  for (const key of Object.keys(tabs)) {
    tabs[key].key = key;
  }

  async function catchQuery(event) {
    tab = 'find';
    await tick();
    tabs.find.instance.performQuery(event.detail);
  }
</script>

<div class="view" class:empty={!collection}>
  {#if collection}
    <TabBar tabs={Object.values(tabs)} bind:selectedKey={tab} />

    {#each Object.values(tabs) as view}
      <div class="container" class:hidden={tab !== view.key}>
        <svelte:component
          this={view.component}
          visible={tab === view.key}
          on:performFind={catchQuery}
          {host}
          {database}
          {collection}
        />
      </div>
    {/each}
  {:else}
    <BlankState label="Select a collection to continue" />
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
