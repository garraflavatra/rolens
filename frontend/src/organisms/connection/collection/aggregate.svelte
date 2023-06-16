<script>
  import Details from '$components/details.svelte';
  import Icon from '$components/icon.svelte';
  import Modal from '$components/modal.svelte';
  import ObjectEditor from '$components/objecteditor.svelte';
  import { aggregationStageDocumentationURL, aggregationStages } from '$lib/mongo';
  import Collation from '$lib/mongo/collation.svelte';
  import { jsonLooseParse, looseJsonIsValid } from '$lib/strings';
  import { Aggregate } from '$wails/go/app/App';
  import { BrowserOpenURL } from '$wails/runtime/runtime';
  import { onMount } from 'svelte';

  export let collection;

  const options = {};
  let stages = [];
  let settingsModalOpen = false;
  $: invalid = !stages.length || stages.some(stage => !stage.data || !looseJsonIsValid(stage.data));

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
    const pipeline = stages.map(stage => {
      return { [stage.type]: jsonLooseParse(stage.data) };
    });
    await Aggregate(
      collection.hostKey,
      collection.dbKey,
      collection.key,
      JSON.stringify(pipeline),
      JSON.stringify(options)
    );
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
            <Icon name="?" />
          </button>
        </label>

        <!-- svelte-ignore a11y-label-has-associated-control -->
        <label class="field">
          <ObjectEditor bind:text={stage.data}
            on:inited={e => {
              e.detail.editor.dispatch({
                changes: {
                  from: 0,
                  to: e.detail.editor.state.doc.length,
                  insert: '{\n\t\n}',
                },
                selection: {
                  anchor: 3,
                },
              });
              e.detail.editor.focus();
            }} />
        </label>
      </Details>
    {/each}

    <button class="button-small" type="button" on:click={addStage}>
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

<Modal title="Advanced aggregation settings" bind:show={settingsModalOpen}>
  <div class="settinggrid">
    <label for="allowDiskUse">Allow disk use</label>
    <div class="field">
      <span class="checkbox">
        <input type="checkbox" id="allowDiskUse" bind:checked={options.allowDiskUse} />
      </span>
    </div>
  </div>

  <Details title="Set custom collation options">
    <Collation bind:collation={options.collation} />
  </Details>
</Modal>

<style>
  .aggregate {
    display: grid;
    gap: 0.5rem;
    grid-template: 1fr auto / 1fr;
  }

  .aggregate :global(details) {
    position: relative;
    z-index: 2;
  }
  .aggregate :global(details + details::before) {
    content: "";
    position: absolute;
    top: -0.6rem;
    left: 50%;
    width: 1px;
    height: 0.6rem;
    background-color: #888;
    z-index: 1;
  }

  .settinggrid {
    margin-bottom: 0.5rem;
  }

  .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .field + .field {
    margin-top: 0.5rem;
  }
</style>
