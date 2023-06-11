<script>
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import input from '$lib/actions/input';
  import { CreateIndex } from '$wails/go/app/App';
  import { createEventDispatcher } from 'svelte';

  export let collection;

  const dispatch = createEventDispatcher();
  const index = { model: [] };

  function addRule() {
    index.model  = [ ...index.model, {} ];
  }

  function removeRule(ruleIndex) {
    index.model.splice(ruleIndex, 1);
    index.model = index.model;
  }

  async function create() {
    const newIndexName = await CreateIndex(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(index));
    if (newIndexName) {
      dispatch('reload');
      dispatch('close');
    }
  }
</script>

<Modal title="Create new index {collection ? `on collection ${collection.key}` : ''}">
  <form on:submit|preventDefault={create}>
    <label class="field name">
      <span class="label">Name</span>
      <input type="text" placeholder="Optional" bind:value={index.name} use:input={{ autofocus: true }} />
    </label>

    <div class="toggles">
      <label class="field">
        <span class="label">Background (legacy)</span>
        <span class="checkbox">
          <input type="checkbox" bind:checked={index.background} />
        </span>
      </label>
      <label class="field">
        <span class="label">Unique</span>
        <span class="checkbox">
          <input type="checkbox" bind:checked={index.unique} />
        </span>
      </label>
      <!-- <label class="field">
        <span class="label">Drop duplicates</span>
        <span class="checkbox">
          <input type="checkbox" bind:checked={index.dropDuplicates} />
        </span>
      </label> -->
      <label class="field">
        <span class="label">Sparse</span>
        <span class="checkbox">
          <input type="checkbox" bind:checked={index.sparse} />
        </span>
      </label>
    </div>

    <div class="model">
      {#each index.model as rule, ruleIndex}
        <div class="row">
          <label class="field">
            <span class="label">Key</span>
            <input type="text" placeholder="_id" bind:value={rule.key}>
          </label>
          <label class="field">
            <select bind:value={rule.sort}>
              <option value={1}>Ascending</option>
              <option value={-1}>Decending</option>
              <option value="hashed" disabled={index.model.length > 1}>Hashed</option>
            </select>
          </label>
          <button type="button" class="btn danger" on:click={() => removeRule(ruleIndex)} disabled={index.model.length < 2}>
            <Icon name="-" />
          </button>
        </div>
      {:else}
        No rules
      {/each}
    </div>
  </form>

  <div class="buttons" slot="footer">
    <button class="btn" on:click={addRule} disabled={index.model.some(r => r.sort === 'hashed')}>
      <Icon name="+" /> Add rule
    </button>
    <button class="btn" on:click={create} disabled={!index.model.length || index.model.some(r => !r.key)}>
      <Icon name="+" /> Create index
    </button>
  </div>
</Modal>

<style>
  .field.name {
    margin-bottom: 0.5rem;
  }

  .toggles {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .model {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
    padding: 0.5rem;
    border: 1px solid #ccc;
  }
  .model .row {
    display: grid;
    grid-template: 1fr / 1fr auto auto;
    gap: 0.5rem;
  }

  .buttons {
    display: flex;
    gap: 0.5rem;
  }
  .buttons:nth-child(2) {
    margin-left: auto;
  }
</style>
