<script>
  import Icon from '$components/icon.svelte';
  import ObjectEditor from '$components/objecteditor.svelte';
  import { convertLooseJson } from '$lib/strings';
  import { RemoveItems } from '$wails/go/app/App';
  import { onMount } from 'svelte';

  export let collection;

  let json = '';
  let many = true;
  let result = undefined;
  let editor;
  $: code = `db.${collection.key}.remove(${json});`;

  async function removeItems() {
    result = await RemoveItems(
      collection.hostKey,
      collection.dbKey,
      collection.key,
      convertLooseJson(json),
      many
    );
  }

  onMount(() => {
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
    editor.focus();
  });
</script>

<form on:submit|preventDefault={removeItems}>
  <!-- svelte-ignore a11y-label-has-associated-control -->
  <label class="field">
    <ObjectEditor bind:text={json} bind:editor />
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
</style>
