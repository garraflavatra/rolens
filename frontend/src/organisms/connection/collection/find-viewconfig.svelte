<script>
  import TabBar from '../../../components/tabbar.svelte';
  import Modal from '../../../components/modal.svelte';
  import Icon from '../../../components/icon.svelte';
  import { views } from '../../../stores';
  import { randomString } from '../../../utils';
  import { input } from '../../../actions';

  export let collection;
  export let show = false;
  export let activeViewKey = 'list';
  export let firstItem = {};

  $: tabs = Object.entries($views).filter(v => (
    v[0] === 'list' || (
      v[1].host === collection.hostKey &&
      v[1].database === collection.dbKey &&
      v[1].collection === collection.key
    )
  )).sort((a, b) => sortTabKeys(a[0], b[0]))
    .map(([ key, v ]) => ({ key, title: v.name, closable: key !== 'list' }));

  function sortTabKeys(a, b) {
    if (a === 'list') {
      return -1;
    }
    if (b === 'list') {
      return 1;
    }
    else {
      return a.localeCompare(b);
    }
  }

  function createView() {
    const newViewKey = randomString();
    $views[newViewKey] = {
      name: 'Table view',
      host: collection.hostKey,
      database: collection.dbKey,
      collection: collection.key,
      type: 'table',
      columns: [ { key: '_id' } ],
    };
    activeViewKey = newViewKey;
  }

  function removeView(viewKey) {
    const keys = Object.keys($views).sort(sortTabKeys);
    const oldIndex = keys.indexOf(viewKey);
    const newKey = keys[oldIndex - 1];
    console.log(keys, oldIndex, newKey);
    activeViewKey = newKey;
    delete $views[viewKey];
    $views = $views;
  }

  function addColumn(before) {
    if (typeof before === 'number') {
      $views[activeViewKey].columns = [
        ...$views[activeViewKey].columns.slice(0, before),
        {},
        ...$views[activeViewKey].columns.slice(before),
      ];
    }
    else {
      $views[activeViewKey].columns = [ ...$views[activeViewKey].columns, {} ];
    }
  }

  function addSuggestedColumns() {
    if ((typeof firstItem !== 'object') || (firstItem === null)) {
      return;
    }
    $views[activeViewKey].columns = Object.keys(firstItem).map(key => ({ key }));
  }

  function moveColumn(oldIndex, delta) {
    const column = $views[activeViewKey].columns[oldIndex];
    const newIndex = oldIndex + delta;

    $views[activeViewKey].columns.splice(oldIndex, 1);
    $views[activeViewKey].columns.splice(newIndex, 0, column);
    $views[activeViewKey].columns = $views[activeViewKey].columns;
  }

  function removeColumn(index) {
    $views[activeViewKey].columns.splice(index, 1);
    $views[activeViewKey].columns = $views[activeViewKey].columns;
  }
</script>

<Modal title="View configuration" bind:show contentPadding={false}>
  <TabBar
    {tabs}
    canAddTab={true}
    on:addTab={createView}
    on:closeTab={e => removeView(e.detail)}
    bind:selectedKey={activeViewKey}
  />

  <div class="options">
    {#if $views[activeViewKey]}
      <div class="meta">
        {#key activeViewKey}
          <label class="field">
            <span class="label">View name</span>
            <input type="text" use:input={{ autofocus: true }} bind:value={$views[activeViewKey].name} disabled={activeViewKey === 'list'} />
          </label>
        {/key}
        <label class="field">
          <span class="label">View type</span>
          <select bind:value={$views[activeViewKey].type} disabled={activeViewKey === 'list'}>
            <option value="list">List view</option>
            <option value="table">Table view</option>
          </select>
        </label>
      </div>

      {#if $views[activeViewKey].type === 'list'}
        <div class="flex">
          <input type="checkbox" id="hideObjectIndicators" bind:checked={$views[activeViewKey].hideObjectIndicators} />
          <label for="hideObjectIndicators">
            Hide object indicators ({'{...}'} and [...]) in list view and show nothing instead
          </label>
        </div>
      {:else if $views[activeViewKey].type === 'table'}
        {#each $views[activeViewKey].columns as column, columnIndex}
          <div class="column">
            <label class="field">
              <input type="text" use:input bind:value={column.key} placeholder="Column keypath" />
            </label>
            <button class="btn" type="button" on:click={() => addColumn(columnIndex)} title="Add column before this one">
              <Icon name="+" />
            </button>
            <button class="btn" type="button" on:click={() => moveColumn(columnIndex, -1)} disabled={columnIndex === 0} title="Move column one position up">
              <Icon name="chev-u" />
            </button>
            <button class="btn" type="button" on:click={() => moveColumn(columnIndex, 1)} disabled={columnIndex === $views[activeViewKey].columns.length - 1} title="Move column one position down">
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
    {/if}
  </div>
</Modal>

<style>
  .options {
    padding: 1rem;
  }

  .meta {
    display: grid;
    grid-template: 1fr / 1fr 1fr;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .flex {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .column {
    display: grid;
    grid-template: 1fr / 1fr repeat(4, auto);
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
</style>
