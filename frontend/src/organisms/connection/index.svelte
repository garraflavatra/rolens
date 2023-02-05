<script>
  import { onMount } from 'svelte';
  import { Hosts, RenameCollection } from '../../../wailsjs/go/app/App';
  import { input } from '../../actions';
  import Modal from '../../components/modal.svelte';
  import HostTree from './hosttree.svelte';
  import { busy, connections } from '../../stores';
  import CollectionDetail from './collection/index.svelte';
  import HostDetail from './hostdetail.svelte';
  import Icon from '../../components/icon.svelte';
  import { EventsOn } from '../../../wailsjs/runtime/runtime';
  import ExportInfo from './export/exportinfo.svelte';
  import DumpInfo from './export/dumpinfo.svelte';
  import Hint from '../../components/hint.svelte';

  export let hosts = {};
  export let activeHostKey = '';
  export let activeDbKey = '';
  export let activeCollKey = '';

  let hostTree;
  let newDb;
  let newColl;

  let showHostDetail = false;
  let hostDetailKey = '';

  let collToRename = '';
  let newCollKey = '';

  let exportInfo;
  let dumpInfo;

  async function getHosts() {
    hosts = await Hosts();
  }

  function createHost() {
    hostDetailKey = '';
    showHostDetail = true;
  }

  function editHost(hostKey) {
    hostDetailKey = hostKey;
    showHostDetail = true;
  }

  function createDatabase() {
    $connections[activeHostKey].databases[newDb.name] = { collections: {} };
    newDb = undefined;
  }

  function openEditCollModal(collKey) {
    newCollKey = collKey;
    collToRename = collKey;
  }

  async function renameCollection() {
    busy.start();
    const ok = await RenameCollection(activeHostKey, activeDbKey, collToRename, newCollKey);
    if (ok) {
      activeCollKey = newCollKey;
      collToRename = '';
      newCollKey = '';
      await hostTree.reload();
    }
    busy.end();
  }

  function createCollection() {
    $connections[activeHostKey].databases[activeDbKey].collections[newColl.name] = {};
    newColl = undefined;
  }

  function exportCollection(collKey) {
    exportInfo = {
      hostKey: activeHostKey,
      dbKey: activeDbKey,
      collKeys: [ collKey ],
    };
  }

  function dumpCollection(collKey) {
    dumpInfo = {
      hostKey: activeHostKey,
      dbKey: activeDbKey,
      collKeys: [ collKey ],
    };
  }

  EventsOn('CreateHost', createHost);
  EventsOn('CreateDatabase', () => newDb = {});
  EventsOn('CreateCollection', () => newColl = {});
  onMount(getHosts);
</script>

<div class="tree">
  <HostTree
    {hosts}
    bind:activeHostKey
    bind:activeCollKey
    bind:activeDbKey
    bind:this={hostTree}
    on:newHost={createHost}
    on:newDatabase={() => newDb = {}}
    on:newCollection={() => newColl = {}}
    on:editHost={e => editHost(e.detail)}
    on:renameCollection={e => openEditCollModal(e.detail)}
    on:exportCollection={e => exportCollection(e.detail)}
    on:dumpCollection={e => dumpCollection(e.detail)}
  />
</div>

<CollectionDetail
  collection={$connections[activeHostKey]?.databases[activeDbKey]?.collections?.[activeCollKey]}
  hostKey={activeHostKey}
  dbKey={activeDbKey}
  collectionKey={activeCollKey}
/>

<HostDetail
  bind:show={showHostDetail}
  on:reload={getHosts}
  hostKey={activeHostKey}
  {hosts}
/>

<ExportInfo bind:info={exportInfo} {hosts} />
<DumpInfo bind:info={dumpInfo} {hosts} />

{#if newDb}
  <Modal bind:show={newDb}>
    <p><strong>Create a database</strong></p>
    <Hint>
      Note: databases in MongoDB do not exist until they have a collection and an item. Your new database will not persist on the server; fill it to have it created.
    </Hint>
    <form on:submit|preventDefault={createDatabase}>
      <label class="field">
        <input type="text" spellcheck="false" bind:value={newDb.name} use:input={{ autofocus: true }} placeholder="New collection name" />
      </label>
      <p class="modal-actions">
        <button class="btn create" type="submit" disabled={!newDb.name?.trim()}>Create database</button>
        <button class="btn secondary" type="button" on:click={() => newDb = undefined}>Cancel</button>
      </p>
    </form>
  </Modal>
{/if}

{#if newColl}
  <Modal bind:show={newColl}>
    <p><strong>Create a collection</strong></p>
    <Hint>
      Note: collections in MongoDB do not exist until they have at least one item. Your new collection will not persist on the server; fill it to have it created.
    </Hint>
    <form on:submit|preventDefault={createCollection}>
      <label class="field">
        <input type="text" spellcheck="false" bind:value={newColl.name} use:input={{ autofocus: true }} placeholder="New collection name" />
      </label>
      <p class="modal-actions">
        <button class="btn create" type="submit" disabled={!newColl.name?.trim()}>Create collection</button>
        <button class="btn secondary" type="button" on:click={() => newColl = undefined}>Cancel</button>
      </p>
    </form>
  </Modal>
{/if}

{#if collToRename}
  <Modal bind:show={collToRename} width="400px">
    <form class="rename" on:submit|preventDefault={renameCollection}>
      <div>Renaming collection <strong>{collToRename}</strong></div>
      <Icon name="arr-d" />
      <label class="field">
        <input type="text" bind:value={newCollKey} use:input={{ autofocus: true }} spellcheck="false" />
      </label>
      <div class="cancelorsave">
        <button class="btn secondary" type="button" on:click={() => collToRename = ''}>Cancel</button>
        <button class="btn" type="submit">Save</button>
      </div>
    </form>
  </Modal>
{/if}

<style>
  .modal-actions {
    display: flex;
    justify-content: space-between;
  }

  .tree {
    padding: 0.5rem;
    background-color: #fff;
  }

  .rename {
    text-align: center;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    align-items: center;
  }
  .rename .field {
    width: 100%;
  }
  .rename input {
    text-align: center;
    width: 100%;
  }
  .rename strong {
    font-weight: 700;
  }
  .rename .cancelorsave {
    display: flex;
    gap: 0.5rem;
    justify-content: space-between;
    width: 100%;
  }
</style>
