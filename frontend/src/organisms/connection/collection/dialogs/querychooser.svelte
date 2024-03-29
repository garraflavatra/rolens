<script>
  import Grid from '$components/grid/grid.svelte';
  import Hint from '$components/hint.svelte';
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import input from '$lib/actions/input.js';
  import hostTree from '$lib/stores/hosttree.js';
  import queries from '$lib/stores/queries.js';
  import { createEventDispatcher } from 'svelte';

  export let queryToSave = undefined;
  export let collection = {};

  const dispatch = createEventDispatcher();
  let gridSelectedPath = [];
  let selectedKey = '';

  function submit() {
    if (queryToSave) {
      queryToSave.hostKey = collection.hostKey;
      queryToSave.dbKey = collection.dbKey;
      queryToSave.collKey = collection.key;

      dispatch('create', { query: queryToSave });
      selectedKey = queryToSave.name;
    }
    else {
      selectActive();
    }
  }

  function selectActive() {
    dispatch('select', { query: $queries[selectedKey] });
  }

  function gridSelect(event) {
    if (event?.detail?.level === 0) {
      selectedKey = event.detail.itemKey;

      if (queryToSave) {
        queryToSave.name = event.detail.itemKey;
      }
    }
  }

  function gridTrigger(event) {
    gridSelect(event);
    selectActive();
  }

  async function gridRemove(event) {
    await queries.remove(event.detail);
  }

  $: if (queryToSave && !queryToSave.name) {
    queryToSave.name = 'New query';
  }
  $: if (queryToSave?.name) {
    gridSelectedPath = [ queryToSave.name ];
  }
  $: if (selectedKey) {
    gridSelectedPath = [ selectedKey ];
  }
</script>

<Modal title={queryToSave ? 'Save query' : 'Load query'} width="500px" on:close>
  <form on:submit|preventDefault={submit}>
    {#if queryToSave}
      <label class="field queryname">
        <span class="label">Query name</span>
        <input type="text" bind:value={queryToSave.name} use:input={{ autofocus: true }} />
      </label>
      <label class="field">
        <textarea bind:value={queryToSave.remarks} placeholder="Remarks…" use:input />
      </label>
    {/if}

    <div class="querylist">
      <Grid
        columns={[ { key: 'n', title: 'Query name' }, { key: 'h', title: 'Host' }, { key: 'ns', title: 'Namespace' } ]}
        key="n"
        items={Object.entries($queries).reduce((object, [ name, query ]) => {
          object[query.name] = {
            n: name,
            h: $hostTree[query.hostKey]?.name || '?',
            ns: `${query.dbKey}.${query.collKey}`,
          };
          return object;
        }, {})}
        showHeaders={true}
        canRemoveItems={true}
        bind:activePath={gridSelectedPath}
        on:select={gridSelect}
        on:trigger={gridTrigger}
        on:removeItem={gridRemove}
      />
    </div>

    {#if queryToSave && Object.keys($queries).includes(queryToSave.name)}
      <Hint>
        You are about to <strong>overwrite</strong> a saved query. Give it
        another name if you do not want to overwrite.
      </Hint>
    {/if}
  </form>

  <svelte:fragment slot="footer">
    {#if queryToSave}
      <button class="button" on:click={submit}>
        <Icon name="save" /> Save query
      </button>
    {:else}
      <button class="button" on:click={submit} disabled={!selectedKey}>
        <Icon name="upload" /> Load query
      </button>
    {/if}
  </svelte:fragment>
</Modal>

<style>
  .field, .querylist {
    margin-bottom: 0.5rem;
  }

  textarea {
    min-height: 75px;
  }

  .querylist {
    border: 1px solid #ccc;
    min-height: 200px;
  }

  .button + :global(.hint) {
    margin-top: 0.5rem;
  }
</style>
