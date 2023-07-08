<script>
  import BlankState from '$components/blankstate.svelte';
  import CodeEditor from '$components/codeeditor.svelte';
  import Icon from '$components/icon.svelte';
  import { javascript } from '@codemirror/lang-javascript';
  import { onDestroy, onMount } from 'svelte';

  export let host = undefined;
  export let database = undefined;
  export let collection = undefined;
  export let visible = false;

  const placeholder = '// Write your script here...';
  const extensions = [ javascript() ];
  let script = '';
  let result = {};
  let copySucceeded = false;
  let timeout;
  let busy = false;
  let editor;

  async function run() {
    busy = true;

    if (collection) {
      result = await collection.executeShellScript(script);
    }
    else if (database) {
      result = await database.executeShellScript(script);
    }
    else if (host) {
      result = await host.executeShellScript(script);
    }

    busy = false;
  }

  async function copyErrorDescription() {
    await navigator.clipboard.writeText(result.errorDescription);
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  $: visible && editor.focus();

  onMount(() => {
    editor.dispatch({
      changes: {
        from: 0,
        to: editor.state.doc.length,
        insert: placeholder,
      },
      selection: {
        from: 0,
        anchor: 0,
        to: placeholder.length,
        head: placeholder.length,
      },
    });
    editor.focus();
  });

  onDestroy(() => clearTimeout(timeout));
</script>

<div class="shell">
  <div class="overflow">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="field">
      <CodeEditor bind:editor bind:text={script} {extensions} />
    </label>
  </div>

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
      <pre>{result.output || ''}{#if result.stderr}<div class="error">{result.stderr}</div>{/if}</pre>
    {/if}
  </div>

  <div class="controls">
    <button class="button" on:click={run}>
      <Icon name="play" /> Run
    </button>

    {#key result}
      <div class="status flash-green">
        {#if result?.status}
          Exit code: {result.status}
        {/if}
      </div>
    {/key}
  </div>
</div>

<style>
  .shell {
    display: grid;
    grid-template: 1fr auto / 1fr 1fr;
  }

  .overflow {
    overflow: auto;
  }

  .field {
    height: 100%;
  }
  .field :global(.editor) {
    border-radius: 0;
  }

  .output {
    background-color: #2e3027;
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
    user-select: text;
    -webkit-user-select: text;
    cursor: text;
  }
  .output pre .error {
    color: #ff8989;
    margin-top: 2px;
  }
  .output :global(.blankstate) {
    margin: auto;
    padding: 0.5rem;
  }

  .controls {
    margin-top: 0.5rem;
    display: flex;
    align-items: center;
    grid-column: 1 / 3;
  }
  .controls .status {
    margin-left: auto;
  }
</style>
