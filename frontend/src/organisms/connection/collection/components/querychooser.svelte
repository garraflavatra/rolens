<script>
  import Grid from '$components/grid.svelte';
  import Hint from '$components/hint.svelte';
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import input from '$lib/actions/input';
  import queries from '$lib/stores/queries';
  import { createEventDispatcher } from 'svelte';

  export let queryToSave = undefined;
  export let collection = {};
  export let show = false;
  export let hosts = {};

  const dispatch = createEventDispatcher();
  let gridSelectedPath = [];
  let selectedKey = '';

  function submit() {
    if (queryToSave) {
      queryToSave.hostKey = collection.hostKey;
      queryToSave.dbKey = collection.dbKey;
      queryToSave.collKey = collection.key;

      const newId = queries.create(queryToSave);

      if (newId) {
        dispatch('created', newId);
        queryToSave = undefined;
        selectedKey = newId;
        select();
      }
    }
    else {
      select();
    }
  }

  function select() {
    dispatch('select', selectedKey);
    show = false;
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
    select();
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

<Modal bind:show title={queryToSave ? 'Save query' : 'Load query'} width="500px">
  <form on:submit|preventDefault={submit}>
    {#if queryToSave}
      <label class="field queryname">
        <span class="label">Query name</span>
        <input type="text" bind:value={queryToSave.name} use:input={{ autofocus: true }} />
      </label>
      <label class="field">
        <textarea bind:value={queryToSave.remarks} placeholder="Remarks…" use:input></textarea>
      </label>
    {/if}

    <div class="querylist">
      <Grid
        columns={[ { key: 'n', title: 'Query name' }, { key: 'h', title: 'Host' }, { key: 'ns', title: 'Namespace' } ]}
        key="n"
        items={Object.entries($queries).reduce((object, [ name, query ]) => {
          object[query.name] = { n: name, h: hosts[query.hostKey]?.name || '?', ns: `${query.dbKey}.${query.collKey}` };
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

    {#if queryToSave}
      <button class="btn" type="submit">
        <Icon name="save" /> Save query
      </button>

      {#if Object.keys($queries).includes(queryToSave.name)}
        <Hint>
          You are about to <strong>overwrite</strong> a saved query. Give it
          another name if you do not want to overwrite.
        </Hint>
      {/if}
    {:else}
      <button class="btn" type="submit" disabled={!selectedKey}>
        <Icon name="upload" /> Load query
      </button>
    {/if}
  </form>
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

  .btn + :global(.hint) {
    margin-top: 0.5rem;
  }
</style>
