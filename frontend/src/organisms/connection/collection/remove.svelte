<script>
  import { RemoveItems } from '../../../../wailsjs/go/app/App';
  import CodeExample from '../../../components/code-example.svelte';

  export let collection;

  let json = '';
  let many = true;
  let result = undefined;
  $: code = `db.${collection.key}.remove(${json});`;

  async function removeItems() {
    result = await RemoveItems(collection.hostKey, collection.dbKey, collection.key, json, many);
  }
</script>

<form on:submit|preventDefault={removeItems}>
  <div class="options">
    <CodeExample {code} />
    <label class="field">
      <span class="label">Many</span>
      <span class="checkbox">
        <input type="checkbox" bind:checked={many} />
      </span>
    </label>
  </div>

  <label class="field">
    <textarea cols="30" rows="10" bind:value={json} placeholder={'{}'} class="code"></textarea>
  </label>

  <div class="flex">
    <div>
      {#key result}
        {#if typeof result === 'number'}
          <span class="flash-green">Removed {result} item{result === 1 ? '' : 's'}</span>
        {/if}
      {/key}
    </div>
    <button type="submit" class="btn">Remove</button>
  </div>
</form>

<style>
  form {
    display: grid;
    grid-template-rows: auto 1fr auto;
    gap: 0.5rem;
  }

  .options {
    display: grid;
    gap: 0.5rem;
    grid-template: 1fr / 1fr auto;
  }

  .flex {
    display: flex;
    justify-content: space-between;
  }
</style>
