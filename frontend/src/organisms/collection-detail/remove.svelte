<script>
  import CodeExample from '../../components/code-example.svelte';

  export let collection;

  let remove = '';
  $: code = `db.${collection.key}.remove(${remove});`;

  function performRemove() {}
</script>

<form on:submit|preventDefault={performRemove}>
  <CodeExample {code} />

  <label class="field">
    <textarea cols="30" rows="10" bind:value={remove} placeholder={'{}'} class="code"></textarea>
  </label>

  <div class="flex">
    <div></div>
    <button type="submit" class="btn" disabled={!remove}>Remove</button>
  </div>
</form>

<style>
  .flex {
    display: flex;
    justify-content: space-between;
    margin-top: 0.5rem;
  }
</style>
