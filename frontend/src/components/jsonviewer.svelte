<script>
  import hljs from 'highlight.js/lib/core';
  import hljsJson from 'highlight.js/lib/languages/json';
  import 'highlight.js/styles/atom-one-dark.css';
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';

  export let json = '';

  hljs.registerLanguage('json', hljsJson);
  $: highlighted = json ? hljs.highlight('json', json).value : '';

  function copy() {
    navigator.clipboard.writeText(json);
  }
</script>

{#if json}
  <Modal bind:show={json} title="JSON viewer" contentPadding={false}>
    <div class="codeblock">
      <div class="buttons">
        <button class="btn" on:click={copy}>
          <Icon name="clipboard" />
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
