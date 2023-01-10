<script>
  import { PerformInsert } from '../../../wailsjs/go/main/App';

  export let collection;

  let input = '';
  let insertedIds;

  $: if (collection) {
    insertedIds = undefined;
  }

  async function insert() {
    insertedIds = await PerformInsert(collection.hostKey, collection.dbKey, collection.key, input);
  }
</script>

<form on:submit|preventDefault={insert}>
  <label class="field">
    <textarea cols="30" rows="10" bind:value={input} placeholder="[]" class="code"></textarea>
  </label>

  <div class="flex">
    <div>
      {#if insertedIds}
        Success! {insertedIds.length} document{insertedIds.length > 1 ? 's' : ''} inserted
      {/if}
    </div>
    <button type="submit" class="btn">Insert</button>
  </div>
</form>

<style>
  .flex {
    display: flex;
    justify-content: space-between;
    margin-top: 0.5rem;
  }
</style>
