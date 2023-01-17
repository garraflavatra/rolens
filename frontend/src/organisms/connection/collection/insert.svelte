<script>
  import { input } from '../../../actions';
  import { createEventDispatcher } from 'svelte';
  import { InsertItems } from '../../../../wailsjs/go/app/App';

  export let collection;

  const dispatch = createEventDispatcher();
  let json = '';
  let insertedIds;

  async function insert() {
    insertedIds = await InsertItems(collection.hostKey, collection.dbKey, collection.key, json);
  }

  function showDocs() {
    dispatch('performFind', {
      query: insertedIds.length === 1
        ? `{ "_id": ${JSON.stringify(insertedIds[0])} }`
        : `{ "_id": { "$in": [ ${insertedIds.map(id => JSON.stringify(id)).join(', ')} ] } }`,
    });
  }
</script>

<form on:submit|preventDefault={insert}>
  <label class="field">
    <textarea
      cols="30"
      rows="10"
      placeholder="[]"
      class="code"
      bind:value={json}
      use:input={{ json: true, autofocus: true }}
    ></textarea>
  </label>

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
      <button type="submit" class="btn" disabled={!json}>Insert</button>
    </div>
  </div>
</form>

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
