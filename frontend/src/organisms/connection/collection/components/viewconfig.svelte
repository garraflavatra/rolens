<script>
  import TabBar from '../../../../components/tabbar.svelte';
  import Modal from '../../../../components/modal.svelte';
  import Icon from '../../../../components/icon.svelte';
  import { views } from '../../../../stores';
  import { randomString } from '../../../../utils';
  import { input } from '../../../../actions';

  export let collection;
  export let show = false;
  export let activeViewKey = 'list';
  export let firstItem = {};

  $: tabs = Object.entries(views.forCollection(collection.hostKey, collection.dbKey, collection.key))
    .sort((a, b) => sortTabKeys(a[0], b[0]))
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
      columns: [ { key: '_id', showInTable: true, inputType: 'objectid', mandatory: true } ],
    };
    activeViewKey = newViewKey;
  }

  function removeView(viewKey) {
    const keys = Object.keys($views).sort(sortTabKeys);
    const oldIndex = keys.indexOf(viewKey);
    const newKey = keys[oldIndex - 1];
    activeViewKey = newKey;
    delete $views[viewKey];
    $views = $views;
  }

  function addColumn(before) {
    if (typeof before === 'number') {
      $views[activeViewKey].columns = [
        ...$views[activeViewKey].columns.slice(0, before),
        { showInTable: true, inputType: 'none' },
        ...$views[activeViewKey].columns.slice(before),
      ];
    }
    else {
      $views[activeViewKey].columns = [
        ...$views[activeViewKey].columns,
        { showInTable: true, inputType: 'none' },
      ];
    }
  }

  function addSuggestedColumns() {
    if ((typeof firstItem !== 'object') || (firstItem === null)) {
      return;
    }

    $views[activeViewKey].columns = Object.keys(firstItem).sort().map(key => ({
      key,
      showInTable: true,
      inputType: 'none',
    }));
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
          <select bind:value={$views[activeViewKey].type} disabled>
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
        <div class="columns">
          {#each $views[activeViewKey].columns as column, columnIndex}
            <div class="column">
              <label class="field">
                <input type="text" use:input bind:value={column.key} placeholder="Column keypath" />
              </label>

              <label class="field" title="Show column in table view">
                <span class="label">
                  <Icon name="table" />
                </span>
                <span class="checkbox">
                  <input type="checkbox" bind:checked={column.showInTable} />
                </span>
              </label>

              <label class="field" title="Input type in form view">
                <span class="label">
                  <Icon name="form" />
                </span>
                <select bind:value={column.inputType}>
                  <option value="none">Hidden in form</option>
                  <optgroup label="Strings">
                    <option value="string">String</option>
                    <option value="objectid">ObjectID</option>
                  </optgroup>
                  <optgroup label="Integers">
                    <option value="int">Integer (32-bit, signed)</option>
                    <option value="uint64">Integer (64-bit, unsigned)</option>
                    <option value="long">Long (64-bit integer, signed)</option>
                  </optgroup>
                  <optgroup label="Floats">
                    <option value="double">Double (64-bit)</option>
                    <option value="decimal">Decimal (128-bit)</option>
                  </optgroup>
                  <optgroup label="Miscellaneous">
                    <option value="bool">Boolean</option>
                    <option value="date">Date</option>
                  </optgroup>
                </select>
              </label>

              <label class="field" title="Mandatory (field must be valid in order to submit form)">
                <span class="label">
                  <Icon name="target" />
                </span>
                <span class="checkbox">
                  <input type="checkbox" bind:checked={column.mandatory} disabled={column.inputType === 'none'} />
                </span>
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
          {:else}
            <p>No columns yet</p>
          {/each}
        </div>
        <button class="btn" on:click={addColumn}>
          <Icon name="+" /> Add column
        </button>
        <button class="btn" on:click={addSuggestedColumns} disabled={!Object.keys(firstItem || {}).length}>
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

  .columns {
    border: 1px solid #ccc;
    overflow: auto;
    padding: 0.5rem 0.5rem 0;
    margin-bottom: 0.5rem;
  }
  .columns p {
    margin-bottom: 0.5rem;
  }
  .columns .column {
    display: grid;
    grid-template: 1fr / 1fr repeat(7, auto);
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
</style>
