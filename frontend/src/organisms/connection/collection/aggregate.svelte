<script>
  import Details from '$components/details.svelte';
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import MongoCollation from '$components/mongo-collation.svelte';
  import input from '$lib/actions/input';
  import { aggregationStageDocumentationURL, aggregationStages } from '$lib/mongo';
  import { Aggregate } from '$wails/go/app/App';
  import { BrowserOpenURL } from '$wails/runtime/runtime';
  import { onMount } from 'svelte';

  export let collection;

  const options = {};
  let stages = [];
  let settingsModalOpen = false;
  $: invalid = !stages.length || stages.some(stage => {
    try {
      JSON.parse(stage.data);
      return false;
    }
    catch {
      return true;
    }
  });

  function addStage() {
    stages = [ ...stages, { type: '$match' } ];
  }

  function deleteStage(index) {
    stages = [ ...stages.slice(0, index), ...stages.slice(index + 1) ];
  }

  function openStageDocs(type) {
    const url = aggregationStageDocumentationURL(type);
    BrowserOpenURL(url);
  }

  async function run() {
    const pipeline = stages.map(stage => ({ [stage.type]: JSON.parse(stage.data) }));
    await Aggregate(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(pipeline), JSON.stringify(options));
  }

  onMount(() => {
    if (!stages.length) {
      stages = [ { type: '$match' } ];
    }
  });
</script>

<form class="aggregate" on:submit|preventDefault={run}>
  <div>
    {#each stages as stage, index}
      <Details title="#{index + 1}: {stage.type}" deletable={true} on:delete={() => deleteStage(index)} initiallyOpen>
        <label class="field">
          <span class="label">Stage type</span>
          <select bind:value={stage.type}>
            {#each aggregationStages as value}
              <option {value}>{value}</option>
            {/each}
          </select>
          <button class="btn secondary" type="button" on:click={() => openStageDocs(stage.type)} title="Open documentation about {stage.type || 'this stage'} on mongodb.org">
            <Icon name="info" />
          </button>
        </label>
        <label class="field">
          <textarea bind:value={stage.data} class="code" use:input={{ type: 'json' }}></textarea>
        </label>
      </Details>
    {/each}

    <button class="btn-sm" type="button" on:click={addStage}>
      <Icon name="+" /> Add stage
    </button>
  </div>

  <div class="controls">
    <div>
      <button class="btn" type="submit" disabled={invalid}>
        <Icon name="play" /> Run pipeline
      </button>
      <button class="btn" type="button" on:click={() => settingsModalOpen = true}>
        <Icon name="cog" /> Settings
      </button>
    </div>

  </div>
</form>

<Modal title="Advanced settings" bind:show={settingsModalOpen}>
  <div class="settinggrid">
    <label for="allowDiskUse">Allow disk use</label>
    <div class="field">
      <span class="checkbox">
        <input type="checkbox" id="allowDiskUse" bind:checked={options.allowDiskUse} />
      </span>
    </div>
  </div>

  <Details title="Set custom collation options">
    <MongoCollation bind:collation={options.collation} />
  </Details>
</Modal>

<style>
  .aggregate {
    display: grid;
    gap: 0.5rem;
    grid-template: 1fr auto / 1fr;
  }

  .settinggrid {
    margin-bottom: 0.5rem;
  }

  .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  textarea {
    min-height: 100px;
    margin-top: 0.5rem;
  }
</style>
