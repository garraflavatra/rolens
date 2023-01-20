<script>
  import { FindItems, Hosts, RemoveItemById, UpdateHost } from '../../../../wailsjs/go/app/App';
  import CodeExample from '../../../components/code-example.svelte';
  import { input } from '../../../actions';
  import ObjectGrid from '../../../components/objectgrid.svelte';
  import Icon from '../../../components/icon.svelte';
  import ObjectViewer from '../../../components/objectviewer.svelte';
  import FindViewConfigModal from './find-viewconfig.svelte';
  import { onMount } from 'svelte';
  import Grid from '../../../components/grid.svelte';

  export let collection;

  const defaults = {
    query: '{}',
    sort: '{ "_id": 1 }',
    fields: '{}',
    skip: 0,
    limit: 15,
  };

  let form = { ...defaults };
  let view = 'list';
  let result = {};
  let submittedForm = {};
  let queryField;
  let activePath = [];
  let objectViewerData;
  let viewConfigModalOpen = false;
  let viewConfig = {};
  $: code = `db.${collection.key}.find(${form.query || '{}'}${form.fields && form.fields !== '{}' ? `, ${form.fields}` : ''}).sort(${form.sort})${form.skip ? `.skip(${form.skip})` : ''}${form.limit ? `.limit(${form.limit})` : ''};`;
  $: lastPage = (submittedForm.limit && result?.results?.length) ? Math.max(0, Math.ceil((result.total - submittedForm.limit) / submittedForm.limit)) : 0;
  $: activePage = (submittedForm.limit && submittedForm.skip && result?.results?.length) ? submittedForm.skip / submittedForm.limit : 0;

  $: collection && refresh();
  $: updateConfig(viewConfig);

  async function getViewConfig() {
    try {
      const hosts = await Hosts();
      viewConfig = hosts?.[collection.hostKey]?.databases?.[collection.dbKey]?.collections?.[collection.key]?.viewConfig || {};
      console.log(hosts, viewConfig);
    }
    catch (e) {
      console.error(e);
    }
  }

  async function updateConfig(viewConfig) {
    try {
      const hosts = await Hosts();
      hosts[collection.hostKey].databases = hosts[collection.hostKey].databases || {};
      hosts[collection.hostKey].databases[collection.dbKey] = hosts[collection.hostKey].databases[collection.dbKey] || {};
      hosts[collection.hostKey].databases[collection.dbKey].collections = hosts[collection.hostKey].databases[collection.dbKey].collections || {};
      hosts[collection.hostKey].databases[collection.dbKey].collections[collection.key] = hosts[collection.hostKey].databases[collection.dbKey].collections[collection.key] || {};
      hosts[collection.hostKey].databases[collection.dbKey].collections[collection.key].viewConfig = viewConfig;
      await UpdateHost(collection.hostKey, JSON.stringify(hosts[collection.hostKey]));
    }
    catch (e) {
      console.error(e);
    }
  }

  async function submitQuery() {
    activePath = [];
    result = await FindItems(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));
    if (result) {
      submittedForm = JSON.parse(JSON.stringify(form));
    }
    resetFocus();
  }

  async function refresh() {
    await getViewConfig();
    await submitQuery();
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

  function openJson(itemId) {
    const item = result?.results?.find(i => i._id == itemId);
    objectViewerData = item;
  }

  function toggleView() {
    view = view === 'table' ? 'list' : 'table';
  }

  export function performQuery(q) {
    form = { ...defaults, ...q };
    submitQuery();
  }

  onMount(refresh);
</script>

<div class="find">
  <form on:submit|preventDefault={submitQuery}>
    <div class="form-row one">
      <label class="field">
        <span class="label">Query or id</span>
        <input type="text" class="code" bind:this={queryField} bind:value={form.query} use:input={{ json: true, autofocus: true }} placeholder={defaults.query} />
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
        <input type="number" min="0" bind:value={form.skip} use:input placeholder={defaults.skip} list="skipstops" />
      </label>

      <label class="field">
        <span class="label">Limit</span>
        <input type="number" min="0" bind:value={form.limit} use:input placeholder={defaults.limit} list="limits" />
      </label>

      <button type="submit" class="btn">Run</button>
    </div>
  </form>

  <CodeExample {code} />

  <div class="result">
    <div class="grid">
      {#key result}
        {#if view === 'table'}
          <Grid
            key="_id"
            columns={viewConfig.columns?.map(c => ({ key: c.key, title: c.key })) || []}
            showHeaders={true}
            items={result.results || []}
            bind:activePath
            on:trigger={e => openJson(e.detail?.itemKey)}
          />
        {:else if view === 'list'}
          <ObjectGrid
            data={result.results}
            hideObjectIndicators={viewConfig?.hideObjectIndicators}
            bind:activePath
            on:trigger={e => openJson(e.detail?.itemKey)}
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
        <button class="btn" on:click={() => viewConfigModalOpen = true} title="Configure view">
          <Icon name="cog" />
        </button>
        <button class="btn" on:click={toggleView} title="Toggle view">
          <Icon name={view === 'table' ? 'list' : 'table'} />
        </button>
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

<ObjectViewer bind:data={objectViewerData} />
<FindViewConfigModal bind:show={viewConfigModalOpen} activeView={view} bind:config={viewConfig} firstItem={result.results?.[0]} />

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
    overflow: auto;
    min-height: 0;
    min-width: 0;
  }
  .result > .grid {
    overflow: auto;
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
  }
  .result > .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>
