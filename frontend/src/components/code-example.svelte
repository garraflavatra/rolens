<script>
  import CodeViewer from './codeviewer.svelte';

  export let code = '';

  let modalCode = '';
</script>

<div class="examplecode" on:pointerdown={() => modalCode = code}>
  <strong>CLI command</strong>
  <code>{code}</code>
</div>

<CodeViewer bind:code={modalCode} language="js" />

<style>
  .examplecode {
    border: 1px solid #ccc;
    border-radius: 10px;
    padding: 0.5rem;
    opacity: 0.6;
  }

  strong {
    display: inline-block;
    margin-right: 1rem;
    font-weight: 700;
  }

  code {
    user-select: all;
  }
</style>
