<script>
  import Modal from '$components/modal.svelte';

  export let show = false;
</script>

<Modal bind:show overflow={false}>
  <div class="about">
    <div class="brand">
      <img src="/logo.png" alt="" class="logo" />
      <div class="text">
        <div class="name">Rolens</div>
        <div class="subtitle">A modest MongoDB client</div>
      </div>
    </div>
    <div class="cols">
      <div>
        d
      </div>
      <div>
        d
      </div>
    </div>
  </div>
</Modal>

<style>
  .about {
    padding: 2rem;
  }

  .brand {
    margin-top: -90px;
    display: flex;
  }
  .brand .logo {
    height: 200px;
  }
  .brand .text {
    align-self: flex-end;
    margin: 0 0 4rem 1rem;
  }
  .brand .text .name {
    font-size: 4rem;
    margin-bottom: 2rem;
    font-weight: 600;
  }
  .brand .text .subtitle {
    font-size: 2rem;
    font-weight: 500;
  }

  .cols {
    display: grid;
    margin-top: 2rem;
    grid-template: 1fr / 1fr 1fr;
    gap: 1rem;
  }
</style>
