<script>
  import { onMount, tick } from 'svelte';
  import { Hosts } from '../../../wailsjs/go/app/App';
  import { input } from '../../actions';
  import Modal from '../../components/modal.svelte';
  import DatabaseList from './dblist.svelte';
  import { busy, connections } from '../../stores';
  import CollectionDetail from './collection/index.svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  let environment;
  let addressBarModalOpen = true;
  let dbList;

  let newDb;
  let newDbInput;
  let newColl;
  let newCollInput;

  $: if (newDb) {
    tick().then(() => newDbInput.focus());
  }

  async function createDatabase() {
    busy.start();
    $connections[activeHostKey].databases[newDb.name] = { collections: {} };
    newDb = undefined;
    await dbList.reload();
    busy.end();
  }

  async function createCollection() {
    busy.start();
    $connections[activeHostKey].databases[activeDbKey].collections[newColl.name] = {};
    newColl = undefined;
    await dbList.reload();
    busy.end();
  }

  onMount(() => {
    window.runtime.Environment().then(e => environment = e);
    Hosts().then(h => hosts = h);
  });
</script>

<DatabaseList
  {hosts}
  bind:activeHostKey
  bind:activeCollKey
  bind:activeDbKey
  bind:this={dbList}
  on:connected={() => addressBarModalOpen = false}
  on:newDatabase={() => newDb = {}}
  on:newCollection={() => newColl = {}}
/>

<CollectionDetail
  collection={$connections[activeHostKey]?.databases[activeDbKey]?.collections?.[activeCollKey]}
  hostKey={activeHostKey}
  dbKey={activeDbKey}
  collectionKey={activeCollKey}
/>

{#if newDb}
  <Modal bind:show={newDb}>
    <p><strong>Create a database</strong></p>
    <p>Note: databases in MongoDB do not exist until they have a collection and an item. Your new database will not persist on the server; fill it to have it created.</p>
    <form on:submit|preventDefault={createDatabase}>
      <label class="field">
        <input type="text" spellcheck="false" bind:value={newDb.name} use:input placeholder="New collection name" bind:this={newDbInput} />
      </label>
      <button class="btn create" type="submit">Create database</button>
    </form>
  </Modal>
{/if}

{#if newColl}
  <Modal bind:show={newColl}>
    <p><strong>Create a collections</strong></p>
    <p>Note: collections in MongoDB do not exist until they have at least one item. Your new collection will not persist on the server; fill it to have it created.</p>
    <form on:submit|preventDefault={createCollection}>
      <label class="field">
        <input type="text" spellcheck="false" bind:value={newColl.name} use:input placeholder="New collection name" bind:this={newCollInput} />
      </label>
      <button class="btn create" type="submit">Create collection</button>
    </form>
  </Modal>
{/if}
