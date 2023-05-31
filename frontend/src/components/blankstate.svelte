<script>
  export let label = 'No items';
  export let image = '/empty.svg';
  export let pale = true;
  export let big = false;
</script>

<div class="blankstate" class:pale class:big>
  <div class="content">
    <img src={image} alt="" />
    <p>{label}</p>
    <slot />
  </div>
</div>

<style>
  .blankstate {
    display: flex;
  }

  .content {
    margin: auto;
    text-align: center;
  }

  img {
    height: 150px;
    width: auto;
  }

  p {
    margin: 2.85rem 0;
    font-size: 1.25rem;
    line-height: 1.25rem;
  }

  .blankstate :global(.btn) {
    font-size: 1.35rem;
    padding: 1rem;
  }

  .big p {
    font-size: 1.9rem;
    line-height: 1.9rem;
  }

  .pale img {
    filter: grayscale(1);
    opacity: 0.4;
  }
  .pale p {
    color: #777;
  }
</style>
