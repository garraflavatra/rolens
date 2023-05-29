<script>
  import Icon from '$components/icon.svelte';
import Modal from '$components/modal.svelte';
  import alink from '$lib/actions/alink';

  export let show = true;
</script>

<Modal bind:show width="400px" title=" ">
  <div class="brand">
    <img src="/logo.png" alt="Rolens logo" />
    <div>
      <div class="title">Rolens</div>
      <div class="description">Intuitive MongoDB <br/> administration tool</div>
    </div>
  </div>

  <hr />

  <div class="info">
    <p class="copy">© Romein van Buren, 2023.</p>
    <p>
      <a href="" use:alink>Documentation</a> |
      <a href="https://github.com/garraflavatra/rolens" use:alink>GitHub</a> |
      <a href="https://github.com/garraflavatra/rolens/issues/new" use:alink>Report a bug</a> |
      <a href="https://github.com/garraflavatra/rolens/blob/main/LICENSE" use:alink>License</a>
    </p>
  </div>
</Modal>

<style>
  .brand {
    display: flex;
    align-items: center;
    justify-content: space-around;
  }
  .brand > img {
    width: 125px;
    flex: 0 1 125px;
  }
  .brand .title {
    font-size: 2.25rem;
    font-weight: 600;
    line-height: 2.5rem;
  }
  .brand .description {
    font-size: 1.5rem;
    line-height: 1.6rem;
  }

  hr {
    margin: 2rem 1rem;
  }

  .info {
    font-size: 1.25rem;
    line-height: 1.25rem;
    margin: 0 1rem 1rem;
    text-align: center;
  }
  .info .copy {
    font-weight: 500;
  }
</style>
