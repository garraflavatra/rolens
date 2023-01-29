<script>
  import { views } from '../../../stores';
  import { createEventDispatcher } from 'svelte';
  import { InsertItems } from '../../../../wailsjs/go/app/App';
  import { input } from '../../../actions';
  import Icon from '../../../components/icon.svelte';
  import Form from './components/form.svelte';
  import ObjectViewer from '../../../components/objectviewer.svelte';

  export let collection;

  const dispatch = createEventDispatcher();
  let json = '';
  let newItems = [];
  let insertedIds;
  let objectViewerData = '';
  let viewType = 'form';
  let formValid = false;
  $: viewsForCollection = views.forCollection(collection.hostKey, collection.dbKey, collection.key);
  $: oppositeViewType = viewType === 'table' ? 'form' : 'table';

  async function insert() {
    if (collection.viewKey === 'list') {
      insertedIds = await InsertItems(collection.hostKey, collection.dbKey, collection.key, json);
    }
    else {
      insertedIds = await InsertItems(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(newItems));
      if (insertedIds) {
        newItems = [];
      }
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
</script>

<form on:submit|preventDefault={insert}>
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
  {:else}
    <div class="form">
      <Form bind:item={newItems[0]} bind:valid={formValid} view={$views[collection.viewKey]} />
    </div>
  {/if}

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
          <Icon name={oppositeViewType} />
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
      <button type="submit" class="btn" disabled={$views[collection.viewKey]?.type === 'list' ? !json : !formValid}>
        <Icon name="+" /> Insert
      </button>
    </div>
  </div>
</form>

<ObjectViewer data={objectViewerData} />

<style>
  form {
    display: grid;
    grid-template-rows: 1fr auto;
    gap: 0.5rem;
  }

  .flex {
    display: flex;
    justify-content: space-between;
  }
</style>
