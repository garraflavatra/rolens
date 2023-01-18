<script>
  import Icon from '../../components/icon.svelte';
  import HostModal from './hostmodal.svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let modalOpen = false;

  $: host = hosts?.[activeHostKey];
</script>

<div class="addressbar">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="address" on:click={() => modalOpen = true}>
    {#if host?.uri}
      {@const split = host.uri.split('://').map(s => s.split('/')).flat()}
      <span class="protocol">{split[0]}://</span><span class="hostname">{split[1]}</span><span class="path">{split.slice(2).join('/')}</span>
    {:else}
      <span class="empty">no host selected</span>
    {/if}
  </div>

  <div class="actions">
    <button class="btn" on:click={() => modalOpen = true}>
      <Icon name="db" />
    </button>
  </div>
</div>

<HostModal bind:modalOpen bind:hosts on:select />

<style>
  .addressbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0.5rem 0.5rem 1rem;
    height: 3rem;
    border: 1px solid #ccc;
    border-radius: 10px;
    background-color: #fff;
  }

  .address .protocol,
  .address .path,
  .address .empty {
    opacity: 0.6;
  }
</style>
