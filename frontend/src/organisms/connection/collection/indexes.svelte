<script>
  export let collection;

  const indexes = [];

  function getIndexes() {}
</script>

<div class="buttons">
  <button class="btn" on:click={getIndexes}>Get indexes</button>
  <button class="btn danger" disabled={!indexes?.length}>Drop selected</button>
  <button class="btn">Create&hellip;</button>
</div>
