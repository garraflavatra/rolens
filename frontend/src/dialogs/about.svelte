<script>
  import Modal from '$components/modal.svelte';
  import alink from '$lib/actions/alink.js';
  import environment from '$lib/stores/environment.js';
</script>

<Modal width="400px" title=" " on:close>
  <div class="brand">
    <img src="/logo.png" alt="Rolens logo" />
    <div>
      <div class="title">
        Rolens
        <span class="version">{$environment.version}</span>
      </div>
      <div class="description">Intuitive MongoDB <br /> administration tool</div>
    </div>
  </div>

  <hr />

  <div class="info">
    <p class="copy">© Romein van Buren, 2022-2025.</p>
    <p>
      <a href="https://garraflavatra.github.io/rolens/" use:alink>Documentation</a> |
      <a href="https://github.com/garraflavatra/rolens" use:alink>GitHub</a> |
      <a href="https://github.com/garraflavatra/rolens/issues/new" use:alink>Report a bug</a> |
      <a href="https://github.com/garraflavatra/rolens/blob/main/LICENSE" use:alink>License</a>
    </p>
  </div>
</Modal>

<style>
  .brand, .info {
    line-height: 1.2;
  }

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
    font-size: 2.25em;
    font-weight: 600;
    margin-bottom: .25em;
  }
  .brand .title .version {
    font-size: 80%;
    font-weight: 300;
    opacity: 0.65;
  }
  .brand .description {
    font-size: 1.5rem;
  }

  hr {
    margin: 2rem 1rem;
  }

  .info {
    font-size: 1.25rem;
    margin: 0 1rem 1rem;
    text-align: center;
  }
  .info .copy {
    font-weight: 500;
  }
</style>
