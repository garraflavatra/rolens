<script>
  export let label = 'No items';
  export let image = '/empty.svg';
</script>

<div class="blankstate">
  <div class="content">
    <img src={image} alt="">
    <p>{label}</p>
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
    height: 100px;
    width: auto;
    filter: grayscale(1);
    opacity: 0.4;
  }

  p {
    margin-top: 2rem;
    color: #777;
  }
</style>
