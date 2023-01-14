<script>
  import { PerformFind } from '../../../wailsjs/go/app/App';
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
    limit: 15,
  };

  let form = { ...defaults };
  let result = {};
  let submittedForm = {};
  let queryField;
  let activeKey = '';
  $: code = `db.${collection.key}.find(${form.query || '{}'}${form.fields && form.fields !== '{}' ? `, ${form.fields}` : ''}).sort(${form.sort})${form.skip ? `.skip(${form.skip})` : ''}${form.limit ? `.limit(${form.limit})` : ''};`;

  async function submitQuery() {
    activeKey = '';
    result = await PerformFind(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));
    if (result) {
      submittedForm = JSON.parse(JSON.stringify(form));
    }
    resetFocus();
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

  function resetFocus() {
    queryField?.focus();
    queryField?.select();
  }

  export function performQuery(q) {
    form = { ...defaults, ...q };
    console.log(form);
    submitQuery();
  }

  onMount(resetFocus);
</script>

<div class="find">
  <form on:submit|preventDefault={submitQuery}>
    <div class="form-row one">
      <label class="field">
        <span class="label">Query or id</span>
        <input type="text" class="code" bind:this={queryField} bind:value={form.query} use:input={{ json: true }} placeholder={defaults.query} />
      </label>

      <label class="field">
        <span class="label">Sort</span>
        <input type="text" class="code" bind:value={form.sort} use:input={{ json: true }} placeholder={defaults.sort} />
      </label>
    </div>

    <div class="form-row two">
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
    <div class="grid">
      {#key result}
        <ObjectGrid data={result.results} bind:activeKey />
      {/key}
    </div>

    <div class="controls">
      <div>
        {#key result}
          <span class="flash-green">Results: {result.total || 0}</span>
        {/key}
      </div>
      <div>
        <button class="btn danger" on:click={remove} disabled={!activeKey}>
          <Icon name="-" />
        </button>
        <button class="btn" on:click={prev} disabled={!submittedForm.limit || (submittedForm.skip <= 0) || !result?.results?.length}>
          <Icon name="chev-l" />
        </button>
        <button class="btn" on:click={next} disabled={!submittedForm.limit || ((result?.results?.length || 0) < submittedForm.limit) || !result?.results?.length}>
          <Icon name="chev-r" />
        </button>
      </div>
    </div>
  </div>
</div>

<style>
  .find {
    display: grid;
    gap: 0.5rem;
    grid-template: auto auto 1fr / 1fr;
  }

  .form-row {
    display: grid;
    gap: 0.5rem;
  }
  .form-row.one {
    grid-template: 1fr / 3fr 2fr;
    margin-bottom: 0.5rem;
  }
  .form-row.two {
    grid-template: 1fr / 5fr 1fr 1fr 1fr;
  }

  .result {
    display: grid;
    grid-template: 1fr auto / 1fr;
    gap: 0.5rem;
  }
  .result > .grid {
    overflow: auto;
  }
  .result > .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
