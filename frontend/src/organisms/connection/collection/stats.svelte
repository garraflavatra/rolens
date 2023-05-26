<script>
  import ObjectGrid from '$components/objectgrid.svelte';

  export let collection;
</script>

<div class="stats">
  <!-- <CodeExample code="db.stats()" /> -->

  <div class="grid">
    <ObjectGrid data={collection.stats} />
  </div>
</div>

<style>
  .stats {
    display: grid;
    gap: 0.5rem;
    grid-template: auto / 1fr;
  }

  .stats .grid {
    overflow: auto;
    min-height: 0;
    min-width: 0;
    border: 1px solid #ccc;
  }
</style>
