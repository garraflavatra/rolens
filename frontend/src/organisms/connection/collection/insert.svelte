<script>
  import Details from '$components/details.svelte';
  import Grid from '$components/grid.svelte';
  import Icon from '$components/icon.svelte';
  import ObjectViewer from '$components/objectviewer.svelte';
  import input from '$lib/actions/input';
  import { randomString } from '$lib/math';
  import { inputTypes } from '$lib/mongo';
  import views from '$lib/stores/views';
  import { capitalise } from '$lib/strings';
  import { InsertItems } from '$wails/go/app/App';
  import { EJSON } from 'bson';
  import { createEventDispatcher } from 'svelte';
  import Form from './components/form.svelte';

  export let collection;

  const dispatch = createEventDispatcher();
  const formValidity = {};

  let json = '';
  let newItems = [];
  let insertedIds;
  let objectViewerData = '';
  let viewType = 'form';
  let allValid = false;

  $: viewsForCollection = views.forCollection(collection.hostKey, collection.dbKey, collection.key);
  $: oppositeViewType = viewType === 'table' ? 'form' : 'table';
  $: allValid = Object.values(formValidity).every(v => v !== false);

  $: {
    if (collection.viewKey === 'list') {
      try {
        newItems = EJSON.parse(json, { relaxed: false });
      }
      catch { /* ok */ }
    }
    else {
      json = EJSON.stringify(newItems, undefined, 2, { relaxed: false });
    }
  }

  $: if ((viewType === 'form') && !newItems?.length)  {
    newItems = [ {} ];
  }

  async function insert() {
    insertedIds = await InsertItems(collection.hostKey, collection.dbKey, collection.key, json);
    if ((collection.viewKey === 'list') && insertedIds) {
      newItems = [];
    }
  }

  function showDocs() {
    dispatch('performFind', {
      query: insertedIds.length === 1
        ? `{ "_id": ${JSON.stringify(insertedIds[0])} }`
        : `{ "_id": { "$in": [ ${insertedIds.map(id => JSON.stringify(id)).join(', ')} ] } }`,
    });
  }

  function switchViewType() {
    viewType = oppositeViewType;
  }

  function showJson() {
    if (viewType === 'form') {
      objectViewerData = { ...(newItems[0] || {}) };
    }
    else if (viewType === 'table') {
      objectViewerData = [ ...newItems ];
    }
  }

  function addRow(beforeIndex = -1) {
    if ((beforeIndex === -1) || (typeof beforeIndex !== 'number')) {
      newItems = [ ...newItems, {} ];
    }
    else {
      newItems = [
        ...newItems.slice(0, beforeIndex),
        {},
        ...newItems.slice(beforeIndex + 1),
      ];
    }
  }

  function deleteRow(index) {
    newItems.splice(index, 1);
    newItems = newItems;
  }
</script>

<form on:submit|preventDefault={insert}>
  <div class="items">
    {#if collection.viewKey === 'list'}
      <label class="field">
        <textarea
          cols="30"
          rows="10"
          placeholder="[]"
          class="code"
          bind:value={json}
          use:input={{ type: 'json', autofocus: true }}
        ></textarea>
      </label>
    {:else if viewType === 'form'}
      <div class="form">
        {#each newItems as item, index}
          <Details
            title="Item #{index + 1} {(item._id !== undefined) ? `(${item._id})` : ''}"
            initiallyOpen={index === 0}
            deletable={true}
            on:delete={() => deleteRow(index)}
          >
            <fieldset>
              <Form
                bind:item={newItems[index]}
                bind:valid={formValidity[index]}
                view={$views[collection.viewKey]}
                emptyHint="This form has no fields. Use the view configurator in the bottom right corner to add fields."
              />
            </fieldset>
          </Details>
        {/each}
      </div>
    {:else if viewType === 'table'}
      <div class="table">
        <Grid
          key="id"
          items={newItems}
          columns={
            $views[collection.viewKey]?.columns
              ?.filter(c => inputTypes.includes(c.inputType))
              .map(c => ({ ...c, id: randomString(8), title: c.key })) || []
          }
          showHeaders={true}
          canSelect={false}
          canRemoveItems={true}
          hideChildrenToggles={true}
          on:addRow={addRow}
          bind:inputsValid={allValid}
        />
      </div>
    {/if}

    {#if collection.viewKey !== 'list'}
      <button class="btn-sm" type="button" on:click={addRow}>
        <Icon name="+" /> Add item
      </button>
    {/if}
  </div>

  <div class="flex">
    <div>
      {#if insertedIds}
        <span class="flash-green">Success! {insertedIds.length} document{insertedIds.length > 1 ? 's' : ''} inserted</span>
      {/if}
    </div>
    <div>
      {#if insertedIds}
        <button class="btn" type="button" on:click={showDocs}>View inserted docs</button>
      {/if}
      {#if collection.viewKey !== 'list'}
        <button class="btn" type="button" on:click={showJson} title="Show JSON">
          <Icon name="code" />
        </button>
        <button class="btn" type="button" on:click={switchViewType} title="Edit as {oppositeViewType}">
          <Icon name={oppositeViewType} /> {capitalise(oppositeViewType)}
        </button>
      {/if}
      <label class="field inline">
        <select bind:value={collection.viewKey}>
          {#each Object.entries(viewsForCollection) as [key, view]}
            <option value={key}>{key === 'list' ? 'Raw JSON' : view.name}</option>
          {/each}
        </select>
        <button class="btn" type="button" on:click={() => dispatch('openViewConfig')} title="Configure view">
          <Icon name="cog" />
        </button>
      </label>
      <button type="submit" class="btn" disabled={$views[collection.viewKey]?.type === 'list' ? !json : !allValid}>
        <Icon name="+" /> Insert
      </button>
    </div>
  </div>
</form>

<ObjectViewer data={objectViewerData} />

<style>
  form {
    display: grid;
    grid-template: 1fr auto / 1fr;
    gap: 0.5rem;
  }

  .items {
    overflow: auto;
  }
  .items .table {
    background-color: #fff;
    border: 1px solid #ccc;
  }

  .flex {
    display: flex;
    justify-content: space-between;
  }
</style>
