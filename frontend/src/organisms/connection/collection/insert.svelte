<script>
  import Details from '$components/details.svelte';
  import Grid from '$components/grid/grid.svelte';
  import Icon from '$components/icon.svelte';
  import ObjectEditor from '$components/editors/objecteditor.svelte';
  import ObjectViewer from '$components/objectviewer.svelte';
  import { randomString } from '$lib/math.js';
  import { inputTypes } from '$lib/mongo/index.js';
  import views from '$lib/stores/views.js';
  import { capitalise, convertLooseJson, jsonLooseParse } from '$lib/strings.js';
  import { InsertItems } from '$wails/go/app/App.js';
  import { EJSON } from 'bson';
  import { createEventDispatcher, onMount } from 'svelte';
  import Form from './components/form.svelte';

  export let collection;
  export let visible = false;

  const dispatch = createEventDispatcher();
  const formValidity = {};

  let editor;
  let json = '';
  let newItems = [];
  let insertedIds;
  let objectViewerData = '';
  let viewType = 'form';
  let allValid = false;
  let viewsForCollection = {};

  $: oppositeViewType = viewType === 'table' ? 'form' : 'table';
  $: allValid = Object.values(formValidity).every(v => v !== false);

  $: {
    if (collection.viewKey === 'list') {
      try {
        newItems = EJSON.deserialize(jsonLooseParse(json), { relaxed: false });
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

  $: if ($views) {
    viewsForCollection = views.forCollection(collection.hostKey, collection.dbKey, collection.key);
  }

  async function insert() {
    insertedIds = await InsertItems(
      collection.hostKey,
      collection.dbKey,
      collection.key,
      convertLooseJson(json)
    );
    if ((collection.viewKey === 'list') && insertedIds) {
      newItems = [];
    }
  }

  function showDocs() {
    dispatch('performFind', {
      query: insertedIds.length === 1 ? `{ "_id": ${JSON.stringify(insertedIds[0])} }` : `{ "_id": { "$in": [ ${insertedIds.map(id => JSON.stringify(id)).join(', ')} ] } }`,
    });
  }

  function switchViewType() {
    viewType = oppositeViewType;
  }

  function showJson() {
    objectViewerData = [ ...newItems ];
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

  function openViewConfig() {
    views.openConfig(collection);
  }

  $: visible && editor.focus();

  onMount(() => {
    if (collection.viewKey === 'list') {
      editor.dispatch({
        changes: {
          from: 0,
          to: editor.state.doc.length,
          insert: '{\n\t\n}',
        },
        selection: {
          anchor: 3,
        },
      });
    }
  });
</script>

<form on:submit|preventDefault={insert}>
  <div class="items">
    {#if collection.viewKey === 'list'}
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <label class="field json">
        <ObjectEditor bind:text={json} bind:editor />
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
          columns={$views[collection.viewKey]?.columns
            ?.filter(c => inputTypes.includes(c.inputType))
            .map(c => {
              return { ...c, id: randomString(8), title: c.key };
            }) || []}
          showHeaders={true}
          canSelect={false}
          canRemoveItems={true}
          hideChildrenToggles={true}
          on:addRow={addRow}
          on:removeItem={() => deleteRow()}
          bind:inputsValid={allValid}
          bind:items={newItems}
        />
      </div>
    {/if}

    {#if collection.viewKey !== 'list'}
      <button class="button-small" type="button" on:click={addRow}>
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
        <button class="button" type="button" on:click={showDocs}>View inserted docs</button>
      {/if}
      {#if collection.viewKey !== 'list'}
        <button class="button" type="button" on:click={showJson} title="Show JSON">
          <Icon name="code" />
        </button>
        <button class="button" type="button" on:click={switchViewType} title="Edit as {oppositeViewType}">
          <Icon name={oppositeViewType} /> {capitalise(oppositeViewType)}
        </button>
      {/if}
      <label class="field inline">
        <select bind:value={collection.viewKey}>
          {#each Object.entries(viewsForCollection) as [ key, view ]}
            <option value={key}>{key === 'list' ? 'Raw JSON' : view.name}</option>
          {/each}
        </select>
        <button class="button" type="button" on:click={openViewConfig} title="Configure view">
          <Icon name="cog" />
        </button>
      </label>
      <button type="submit" class="button" disabled={$views[collection.viewKey]?.type === 'list' ? !json : !allValid}>
        <Icon name="+" /> Insert
      </button>
    </div>
  </div>
</form>

{#if objectViewerData}
  <ObjectViewer bind:data={objectViewerData} />
{/if}

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

  .field.json {
    height: 100%;
  }
</style>
