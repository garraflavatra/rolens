<script>
  import { PerformFind } from '../../../wailsjs/go/main/App';
  import CodeExample from '../../components/code-example.svelte';
  import { onMount } from 'svelte';
  import { input } from '../../actions';
  import ObjectGrid from '../../components/objectgrid.svelte';
  import Icon from '../../components/icon.svelte';

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
  let activeKey = '';
  $: code = `db.${collection.key}.find(${form.query || '{}'}${form.fields && form.fields !== '{}' ? `, ${form.fields}` : ''}).sort(${form.sort})${form.skip ? `.skip(${form.skip})` : ''}${form.limit ? `.limit(${form.limit})` : ''};`;

  $: if (collection) {
    result = [];
  }

  async function submitQuery() {
    activeKey = '';
    result = await PerformFind(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));
    queryField?.focus();
    queryField?.select();
  }

  function prev() {
    form.skip -= form.limit;
    if (form.skip < 0) {
      form.skip = 0;
    }
    submitQuery();
  }

  function next() {
    form.skip += form.limit;
    submitQuery();
  }

  function remove() {
    // eslint-disable-next-line no-alert
    alert('yet to be implemented');
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

<div class="result">
  <ObjectGrid data={result} bind:activeKey />
  <div class="controls">
    <div>
      {#if result}
        Results: {result.length}
      {/if}
    </div>
    <div>
      <button class="btn danger" on:click={remove} disabled={!activeKey}>
        <Icon name="-" />
      </button>
      <button class="btn" on:click={prev} disabled={!form.limit || (form.skip <= 0) || !result?.length}>
        <Icon name="chev-l" />
      </button>
      <button class="btn" on:click={next} disabled={!form.limit || ((result?.length || Infinity) < form.limit) || !result?.length}>
        <Icon name="chev-r" />
      </button>
    </div>
  </div>
</div>

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

  .result {
    flex: 1;
    display: flex;
    flex-flow: column;
    margin-top: 0.5rem;
    gap: 0.5rem;
  }
  .result > :global(.grid) {
    flex: 1;
  }
  .result > .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
