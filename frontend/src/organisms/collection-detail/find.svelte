<script>
  import { PerformFind } from '../../../wailsjs/go/main/App';
  import CodeExample from '../../components/code-example.svelte';
  import { onMount } from 'svelte';
  import { input } from '../../actions';
  import ObjectGrid from '../../components/objectgrid.svelte';

  export let collection;

  const defaults = {
    query: '{}',
    sort: '{ "_id": 1 }',
    fields: '{}',
    skip: 0,
    limit: 30,
  };

  const form = {
    query: '{}',
    sort: '{ "_id": 1 }',
    fields: '{}',
    skip: 0,
    limit: 30,
  };

  let result = [];
  let queryField;
  $: code = `db.${collection.key}.find(${form.query || '{}'}${form.fields && form.fields !== '{}' ? `, ${form.fields}` : ''}).sort(${form.sort})${form.skip ? `.skip(${form.skip})` : ''}${form.limit ? `.limit(${form.limit})` : ''};`;

  $: if (collection) {
    result = [];
  }

  async function submitQuery() {
    result = await PerformFind(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));
  }

  onMount(() => {
    queryField?.focus();
    queryField?.select();
  });
</script>

<form on:submit|preventDefault={submitQuery}>
  <div class="row-one">
    <label class="field">
      <span class="label">Query or id</span>
      <input type="text" class="code" bind:this={queryField} bind:value={form.query} use:input={{ json: true }} placeholder={defaults.query} />
    </label>

    <label class="field">
      <span class="label">Sort</span>
      <input type="text" class="code" bind:value={form.sort} use:input={{ json: true }} placeholder={defaults.sort} />
    </label>
  </div>

  <div class="row-two">
    <label class="field">
      <span class="label">Fields</span>
      <input type="text" class="code" bind:value={form.fields} use:input={{ json: true }} placeholder={defaults.fields} />
    </label>

    <label class="field">
      <span class="label">Skip</span>
      <input type="number" min="0" bind:value={form.skip} use:input placeholder={defaults.skip} />
    </label>

    <label class="field">
      <span class="label">Limit</span>
      <input type="number" min="0" bind:value={form.limit} use:input placeholder={defaults.limit} />
    </label>

    <button type="submit" class="btn">Run</button>
  </div>
</form>

<CodeExample {code} />
<ObjectGrid data={result} />

<style>
  .row-one {
    display: grid;
    gap: 0.5rem;
    grid-template-columns: 3fr 2fr;
    margin-bottom: 0.5rem;
  }
  .row-two {
    display: grid;
    gap: 0.5rem;
    grid-template-columns: 5fr 1fr 1fr 1fr;
    margin-bottom: 0.5rem;
  }
</style>
