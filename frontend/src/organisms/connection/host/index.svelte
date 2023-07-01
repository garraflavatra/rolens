<script>
  import BlankState from '$components/blankstate.svelte';
  import TabBar from '$components/tabbar.svelte';
  import { EventsOn } from '$wails/runtime/runtime';

  import Logs from './logs.svelte';
  import Shell from '../shell.svelte';
  import Status from './status.svelte';
  import SystemInfo from './systeminfo.svelte';

  export let host;
  export let hostKey;
  export let tab = 'status';

  $: if (host) {
    host.hostKey = hostKey;
  }

  $: if (hostKey) {
    tab = 'status';
  }

  EventsOn('OpenStatusTab', name => (tab = name || tab));
</script>

<div class="view" class:empty={!host}>
  {#if host}
    {#key host}
      <TabBar
        tabs={[
          { key: 'status', icon: 'chart', title: 'Host status' },
          { key: 'shell', icon: 'shell', title: 'Shell' },
          { key: 'logs', icon: 'doc', title: 'Logs' },
          { key: 'systemInfo', icon: 'server', title: 'System info' },
        ]}
        bind:selectedKey={tab}
      />

      <div class="container">
        {#if tab === 'status'} <Status {host} />
        {:else if tab === 'logs'} <Logs {host} />
        {:else if tab === 'systemInfo'} <SystemInfo {host} />
        {:else if tab === 'shell'} <Shell {host} />
        {/if}
      </div>
    {/key}
  {:else}
    <BlankState label="Select a host to continue" />
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
