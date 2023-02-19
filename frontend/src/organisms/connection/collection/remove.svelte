<script>
  // import CodeExample from '$components/code-example.svelte';
  import Icon from '$components/icon.svelte';
  import input from '$lib/actions/input';
  import { RemoveItems } from '$wails/go/app/App';

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
  <!-- <CodeExample {code} /> -->

  <label class="field">
    <textarea
      cols="30"
      rows="10"
      placeholder={'{}'}
      class="code"
      bind:value={json}
      use:input={{ type: 'json', autofocus: true }}
    ></textarea>
  </label>

  <div class="actions">
    <button type="submit" class="btn danger">
      <Icon name="-" /> Remove
    </button>

    <label class="field many">
      <span class="label">Many</span>
      <span class="checkbox">
        <input type="checkbox" bind:checked={many} />
      </span>
    </label>

    {#key result}
      {#if typeof result === 'number'}
        <span class="flash-green">Removed {result} item{result === 1 ? '' : 's'}</span>
      {/if}
    {/key}
  </div>
</form>

<style>
  form {
    display: grid;
    grid-template-rows: 1fr auto;
    gap: 0.5rem;
  }

  .many {
    display: inline-flex;
  }

  textarea {
    resize: none;
  }
</style>
