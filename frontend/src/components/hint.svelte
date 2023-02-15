<script>
  import Icon from './icon.svelte';

</script>
<div class="hint">
  <Icon name="info" />
  <div class="content"><slot /></div>
</div>

<style>
  .hint {
    padding: 0.5rem;
    margin-bottom: 0.5rem;
    display: grid;
    gap: 1rem;
    grid-template: 1fr / auto 1fr;
    background-color: rgba(255, 214, 64, 0.5);
    border: 1px solid #c5a42e;
    border-radius: 10px;
    align-items: center;
  }

  .content {
    align-self: center;
    line-height: 1.3;
  }
  .content :global(strong) {
    font-weight: 700;
  }
</style>
