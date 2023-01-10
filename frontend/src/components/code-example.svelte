<script>
  export let code = '';
</script>

<div class="examplecode">
  <strong>CLI command</strong>
  <code>{code}</code>
</div>

<style>
  .examplecode {
    border: 1px solid #ccc;
    border-radius: 10px;
    padding: 0.5rem;
    opacity: 0.6;
    margin-bottom: 0.5rem;
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
