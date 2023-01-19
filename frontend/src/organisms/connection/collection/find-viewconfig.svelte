<script>
  import TabBar from '../../../components/tabbar.svelte';
  import Modal from '../../../components/modal.svelte';

  export let show = false;
  export let activeView = 'list';
  export let config = {
    hideObjectIndicators: false,
    columns: [],
  };

  let activeTab = activeView || 'list';
</script>

<Modal title="View configuration" bind:show contentPadding={false}>
  <TabBar
    tabs={[
      { key: 'list', title: 'List view' },
      { key: 'table', title: 'Table view' },
    ]}
    bind:selectedKey={activeTab}
  />

  <div class="options">
    {#if activeTab === 'list'}
      <div class="flex">
        <input type="checkbox" id="hideObjectIndicators" bind:checked={config.hideObjectIndicators} />
        <label for="hideObjectIndicators">
          Hide object indicators ({'{...}'} and [...]) in list view and show nothing instead
        </label>
      </div>
    {:else if activeTab === 'table'}
      <input
        type="text"
        value={config.columns?.map(c => c.key).join(', ') || ''}
        on:input={e => config.columns = e.currentTarget.value?.split(',').map(k => ({ key: k.trim() })) || ''}
      />
    {/if}
  </div>
</Modal>

<style>
  .options {
    padding: 1rem;
  }

  .flex {
    display: flex;
    gap: 0.5rem;
  }
  .flex + .flex {
    margin-top: 1rem;
  }
</style>
