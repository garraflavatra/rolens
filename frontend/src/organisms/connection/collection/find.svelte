<script>
  import Grid from '$components/grid.svelte';
  import Icon from '$components/icon.svelte';
  import ObjectGrid from '$components/objectgrid.svelte';
  import ObjectViewer from '$components/objectviewer.svelte';
  import input from '$lib/actions/input';
  import { deepClone } from '$lib/objects';
  import { startProgress } from '$lib/progress';
  import applicationSettings from '$lib/stores/settings';
  import views from '$lib/stores/views';
  import { FindItems, RemoveItemById, UpdateFoundDocument } from '$wails/go/app/App';
  import { EJSON } from 'bson';
  import { createEventDispatcher, onMount } from 'svelte';

  export let collection;

  const dispatch = createEventDispatcher();
  const defaults = {
    query: '{}',
    sort: $applicationSettings.defaultSort || '{ "_id": 1 }',
    fields: '{}',
    skip: 0,
    limit: $applicationSettings.defaultLimit || 15,
  };

  let form = { ...defaults };
  let result = {};
  let submittedForm = {};
  let queryField;
  let activePath = [];
  let objectViewerData;
  let querying = false;
  let objectViewerSuccessMessage = '';

  // $: code = `db.${collection.key}.find(${form.query || '{}'}${form.fields && form.fields !== '{}' ? `, ${form.fields}` : ''}).sort(${form.sort})${form.skip ? `.skip(${form.skip})` : ''}${form.limit ? `.limit(${form.limit})` : ''};`;
  $: viewsForCollection = views.forCollection(collection.hostKey, collection.dbKey, collection.key);
  $: lastPage = (submittedForm.limit && result?.results?.length) ? Math.max(0, Math.ceil((result.total - submittedForm.limit) / submittedForm.limit)) : 0;
  $: activePage = (submittedForm.limit && submittedForm.skip && result?.results?.length) ? submittedForm.skip / submittedForm.limit : 0;

  async function submitQuery() {
    if (querying) {
      return;
    }

    querying = true;
    const progress = startProgress('Performing query…');
    activePath = [];
    const newResult = await FindItems(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));

    if (newResult) {
      newResult.results = newResult.results?.map(s => EJSON.parse(s, { relaxed: false }));
      result = newResult;
      submittedForm = deepClone(form);
    }

    progress.end();
    resetFocus();
    querying = false;
  }

  async function refresh() {
    if ($applicationSettings.autosubmitQuery) {
      await submitQuery();
    }
  }

  async function loadQuery() {
    const query = await collection.openQueryChooser();
    if (query) {
      form = { ...query };
      submitQuery();
    }
  }

  async function saveQuery() {
    const query = await collection.openQueryChooser(form);
    if (query) {
      form = { ...query };
      submitQuery();
    }
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

  function first() {
    form.skip = 0;
    submitQuery();
  }

  function last() {
    form.skip = lastPage * submittedForm.limit;
    submitQuery();
  }

  async function removeActive() {
    if (!activePath[0]) {
      return;
    }
    const ok = await RemoveItemById(collection.hostKey, collection.dbKey, collection.key, activePath[0]);
    if (ok) {
      await submitQuery();
    }
  }

  function resetFocus() {
    queryField?.focus();
    queryField?.select();
  }

  function openJson(index) {
    const item = result?.results?.[index];
    objectViewerData = item;
  }

  export function performQuery(q) {
    form = { ...defaults, ...q };
    submitQuery();
  }

  async function saveDocument(event) {
    const progress = startProgress('Performing update…');
    const success = await UpdateFoundDocument(
      collection.hostKey,
      collection.dbKey,
      collection.key,
      EJSON.stringify({ _id: event.detail.originalData._id }),
      event.detail.text
    );

    if (success) {
      objectViewerSuccessMessage = 'Document has been saved!';
      submitQuery();
    }

    progress.end();
  }

  $: collection && refresh();
  onMount(refresh);
</script>

<div class="find">
  <form on:submit|preventDefault={submitQuery}>
    <div class="form-row one">
      <label class="field">
        <span class="label">Query or id</span>
        <input type="text"
          class="code"
          bind:this={queryField}
          bind:value={form.query}
          use:input={{ type: 'json', autofocus: true }}
          placeholder={defaults.query} />
      </label>

      <label class="field">
        <span class="label">Sort</span>
        <input type="text" class="code" bind:value={form.sort} use:input={{ type: 'json' }} placeholder={defaults.sort} />
      </label>
    </div>

    <div class="form-row two">
      <label class="field">
        <span class="label">Fields</span>
        <input type="text" class="code" bind:value={form.fields} use:input={{ type: 'json' }} placeholder={defaults.fields} />
      </label>

      <label class="field">
        <span class="label">Skip</span>
        <input type="number"
          min="0"
          bind:value={form.skip}
          use:input
          placeholder={defaults.skip}
          list="skipstops" />
      </label>

      <label class="field">
        <span class="label">Limit</span>
        <input type="number"
          min="0"
          bind:value={form.limit}
          use:input
          placeholder={defaults.limit}
          list="limits" />
      </label>
    </div>

    <div class="form-row actions">
      <button type="submit" class="btn" title="Run query">
        <Icon name="play" /> Run
      </button>
      <button class="btn secondary" type="button" on:click={() => collection.export(form)}>
        <Icon name="save" /> Export results…
      </button>
      <div class="field">
        <button class="btn secondary" type="button" on:click={loadQuery}>
          <Icon name="upload" /> Load query…
        </button>
        <button class="btn secondary" type="button" on:click={saveQuery}>
          <Icon name="save" /> Save as…
        </button>
      </div>
    </div>
  </form>

  <div class="result">
    <div class="grid">
      {#key result}
        {#if collection.viewKey === 'list'}
          <ObjectGrid
            data={result.results}
            hideObjectIndicators={$views[collection.viewKey]?.hideObjectIndicators}
            bind:activePath
            on:trigger={e => openJson(e.detail?.index)}
          />
        {:else}
          <Grid
            key="_id"
            columns={$views[collection.viewKey]?.columns?.map(c => {
              return { key: c.key, title: c.key };
            }) || []}
            showHeaders={true}
            items={result.results ? result.results.map(r => EJSON.deserialize(r)) : []}
            bind:activePath
            on:trigger={e => openJson(e.detail?.index)}
          />
        {/if}
      {/key}
    </div>

    <div class="controls">
      <div>
        {#key result}
          <span class="flash-green">Results: {result.total || 0}</span>
        {/key}
      </div>
      <div>
        <label class="field inline">
          <select bind:value={collection.viewKey}>
            {#each Object.entries(viewsForCollection) as [ key, view ]}
              <option value={key}>{view.name}</option>
            {/each}
          </select>
          <button class="btn" on:click={() => dispatch('openViewConfig', { firstItem: result.results?.[0] })} title="Configure view">
            <Icon name="cog" />
          </button>
        </label>
        <button class="btn danger" on:click={removeActive} disabled={!activePath?.length} title="Drop selected item">
          <Icon name="-" />
        </button>
        <button class="btn" on:click={first} disabled={!submittedForm.limit || (submittedForm.skip <= 0) || !result?.results || (activePage === 0)} title="First page">
          <Icon name="chevs-l" />
        </button>
        <button class="btn" on:click={prev} disabled={!submittedForm.limit || (submittedForm.skip <= 0) || !result?.results || (activePage === 0)} title="Previous {submittedForm.limit} items">
          <Icon name="chev-l" />
        </button>
        <button class="btn" on:click={next} disabled={!submittedForm.limit || ((result?.results?.length || 0) < submittedForm.limit) || !result?.results || !lastPage || (activePage >= lastPage)} title="Next {submittedForm.limit} items">
          <Icon name="chev-r" />
        </button>
        <button class="btn" on:click={last} disabled={!submittedForm.limit || ((result?.results?.length || 0) < submittedForm.limit) || !result?.results || !lastPage || (activePage >= lastPage)} title="Last page">
          <Icon name="chevs-r" />
        </button>
      </div>
    </div>
  </div>
</div>


{#if objectViewerData}
  <!-- @todo Implement save -->
  <ObjectViewer bind:data={objectViewerData} saveable on:save={saveDocument} bind:successMessage={objectViewerSuccessMessage} />
{/if}

<datalist id="limits">
  {#each [ 1, 5, 10, 25, 50, 100, 200 ] as value}
    <option {value} />
  {/each}
</datalist>

{#if submittedForm?.limit}
  <datalist id="skipstops">
    {#each Array(lastPage).fill('').map((_, i) => i * submittedForm.limit) as value}
      <option {value} />
    {/each}
  </datalist>
{/if}

<style>
  .find {
    display: grid;
    gap: 0.5rem;
    grid-template: auto 1fr / 1fr;
  }

  .form-row {
    display: grid;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
  .form-row.one {
    grid-template: 1fr / 3fr 2fr;
  }
  .form-row.two {
    grid-template: 1fr / 5fr 1fr 1fr;
  }
  .form-row.actions {
    margin-bottom: 0rem;
    grid-template: 1fr / repeat(4, auto);
    justify-content: start;
  }

  .result {
    display: grid;
    grid-template: 1fr auto / 1fr;
    gap: 0.5rem;
    overflow: auto;
    min-height: 0;
    min-width: 0;
  }
  .result > .grid {
    overflow: auto;
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
    background-color: #fff;
  }
  .result > .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
