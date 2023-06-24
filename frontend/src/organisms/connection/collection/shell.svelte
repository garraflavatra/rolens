<script>
  import BlankState from '$components/blankstate.svelte';
  import CodeEditor from '$components/codeeditor.svelte';
  import Icon from '$components/icon.svelte';
  import { javascript } from '@codemirror/lang-javascript';
  import { onDestroy } from 'svelte';

  export let collection;

  const extensions = [ javascript() ];
  let script = '';
  let result = {};
  let copySucceeded = false;
  let timeout;
  let busy = false;

  async function run() {
    busy = true;
    result = await collection.executeShellScript(script);
    busy = false;
  }

  async function copyErrorDescription() {
    await navigator.clipboard.writeText(result.errorDescription);
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  onDestroy(() => clearTimeout(timeout));
</script>

<div class="shell">
  <div class="panels">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="field">
      <CodeEditor bind:text={script} {extensions} />
    </label>

    <div class="output">
      {#if busy}
        <BlankState icon="loading" label="Executing…" />
      {:else if result.errorTitle || result.errorDescription}
        <BlankState title={result.errorTitle} label={result.errorDescription} icon="!">
          <button class="button-small" on:click={copyErrorDescription}>
            <Icon name={copySucceeded ? 'check' : 'clipboard'} /> Copy error message
          </button>
        </BlankState>
      {:else}
        <pre>{result.output || ''}</pre>
      {/if}
    </div>
  </div>

  <div class="controls">
    {#key result}
      <div class="status flash-green">
        {#if result?.status}
          Exit code: {result.status}
        {/if}
      </div>
    {/key}

    <button class="btn" on:click={run}>
      <Icon name="play" /> Run
    </button>
  </div>
</div>

<style>
  .shell {
    display: grid;
    grid-template: 1fr auto / 1fr;
  }

  .panels {
    display: flex;
    height: 100%;
  }
  .panels > * {
    flex: 1 1 0;
    width: 100%;
  }

  .field :global(.editor) {
    border-radius: 0;
  }

  .output {
    background-color: #111;
    color: #fff;
    overflow: auto;
    display: flex;
  }
  .output :global(*) {
    color: #fff;
  }
  .output pre {
    font-family: monospace;
    padding: 0.5rem;
  }
  .output :global(.blankstate) {
    margin: auto;
    padding: 0.5rem;
  }

  .controls {
    margin-top: 0.5rem;
    display: flex;
    align-items: center;
  }
  .controls .status {
    margin-right: auto;
  }
</style>
