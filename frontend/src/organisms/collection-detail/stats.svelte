<script>
  import ObjectGrid from '../../components/objectgrid.svelte';
  import CodeExample from '../../components/code-example.svelte';

  export let collection;
</script>

<div class="stats">
  <CodeExample code="db.stats()" />
  <ObjectGrid data={collection.stats} />
</div>

<style>
  .stats {
    display: grid;
    gap: 0.5rem;
    grid-template: auto 1fr / 1fr;
  }
</style>
