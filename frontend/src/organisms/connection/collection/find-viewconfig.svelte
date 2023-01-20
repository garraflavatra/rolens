<script>
  import TabBar from '../../../components/tabbar.svelte';
  import Modal from '../../../components/modal.svelte';
  import Icon from '../../../components/icon.svelte';

  export let show = false;
  export let activeView = 'list';
  export let config = {
    hideObjectIndicators: false,
    columns: [],
  };
  export let firstItem = {};

  let activeTab = activeView || 'list';

  $: activeView && (activeTab = activeView);
  $: if (!config.columns || (config.columns.length === 0)) {
    config.columns = [ { key: '_id' } ];
  }

  function addColumn(before) {
    if (typeof before === 'number') {
      config.columns = [
        ...config.columns.slice(0, before),
        {},
        ...config.columns.slice(before),
      ];
    }
    else {
      config.columns = [ ...config.columns, {} ];
    }
  }

  function addSuggestedColumns() {
    if ((typeof firstItem !== 'object') || (firstItem === null)) {
      return;
    }
    config.columns = Object.keys(firstItem).map(key => ({ key }));
  }

  function moveColumn(oldIndex, delta) {
    const column = config.columns[oldIndex];
    const newIndex = oldIndex + delta;

    config.columns.splice(oldIndex, 1);
    config.columns.splice(newIndex, 0, column);
    config.columns = config.columns;
  }

  function removeColumn(index) {
    config.columns.splice(index, 1);
    config.columns = config.columns;
  }
</script>

<Modal title="View configuration" bind:show contentPadding={false}>
  <TabBar
    tabs={[
      { key: 'list', title: 'List view' },
      { key: 'table', title: 'Table view columns' },
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
      {#each config.columns as column, columnIndex}
        <div class="column">
          <label class="field">
            <input type="text" bind:value={column.key} placeholder="Column keypath" />
          </label>
          <button class="btn" type="button" on:click={() => addColumn(columnIndex)} title="Add column before this one">
            <Icon name="+" />
          </button>
          <button class="btn" type="button" on:click={() => moveColumn(columnIndex, -1)} disabled={columnIndex === 0} title="Move column one position up">
            <Icon name="chev-u" />
          </button>
          <button class="btn" type="button" on:click={() => moveColumn(columnIndex, 1)} disabled={columnIndex === config.columns.length - 1} title="Move column one position down">
            <Icon name="chev-d" />
          </button>
          <button class="btn danger" type="button" on:click={() => removeColumn(columnIndex)} title="Remove this column">
            <Icon name="x" />
          </button>
        </div>
      {/each}
      <button class="btn" on:click={addColumn}>
        <Icon name="+" /> Add column
      </button>
      <button class="btn" on:click={addSuggestedColumns} disabled={!firstItem}>
        <Icon name="zap" /> Add suggested columns
      </button>
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

  .column {
    display: grid;
    grid-template: 1fr / 1fr repeat(4, auto);
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
</style>
