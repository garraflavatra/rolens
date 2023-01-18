<script>
  import hljs from 'highlight.js/lib/core';
  import hljsJSON from 'highlight.js/lib/languages/json';
  import hljsJavaScript from 'highlight.js/lib/languages/javascript';
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';
  import 'highlight.js/styles/atom-one-dark.css';
  import { onDestroy } from 'svelte';

  export let code = '';
  export let language = 'json';

  const languageNames = {
    json: 'JSON',
    js: 'JavaScript',
  };

  hljs.registerLanguage('json', hljsJSON);
  hljs.registerLanguage('js', hljsJavaScript);

  let copySucceeded = false;
  let timeout;
  $: highlighted = code ? hljs.highlight(code, { language }).value : '';

  async function copy() {
    await navigator.clipboard.writeText(code);
    copySucceeded = true;
    timeout = setTimeout(() => copySucceeded = false, 1500);
  }

  onDestroy(() => clearTimeout(timeout));
</script>

{#if code}
  <Modal bind:show={code} title="{languageNames[language]} viewer" contentPadding={false}>
    <div class="codeblock">
      <div class="buttons">
        <button class="btn" on:click={copy}>
          <Icon name={copySucceeded ? 'check' : 'clipboard'} />
        </button>
      </div>
      <pre><code class="hljs">{@html highlighted}</code></pre>
    </div>
  </Modal>
{/if}

<style>
  .codeblock {
    position: relative;
  }
  .buttons {
    position: absolute;
    top: 0;
    right: 0;
    margin: 1rem;
  }
  .buttons button {
    margin-left: 1rem;
    background: none;
    border: 1px solid rgba(255, 255, 255, 0.3);
  }
  .buttons button:hover {
    background-color: rgba(0, 0, 0, 0.3);
  }

  .hljs {
    user-select: text;
    line-height: 1.2;
  }

  pre :global(::selection) {
    background: rgba(5, 5, 5, 0.8);
  }
</style>
