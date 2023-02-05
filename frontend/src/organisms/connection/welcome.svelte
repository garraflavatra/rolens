<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function createHost() {
    dispatch('createHost');
  }
</script>

<div class="welcome">
  <div class="brand">
    <img src="/logo.png" alt="" class="logo" />
    <div class="text">
      <div class="name">Welcome to Rolens!</div>
      <div class="subtitle">A modest MongoDB client</div>
    </div>
  </div>
  <button class="btn" on:click={createHost}>Create your first host</button>
</div>

<style>
  .welcome {
    /* transform: translateY(-80px); */
    margin-top: -90px;
    padding: 2rem;
  }

  .brand {
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
    font-size: 2.5rem;
    margin-bottom: 1.5rem;
    font-weight: 600;
  }
  .brand .text .subtitle {
    font-size: 1.5rem;
  }

  .logo {
    height: 250px;
  }

  /* .title {
    font-weight: 600;
    font-size: 1.5rem;
  } */

  .btn {
    margin-top: 2rem;
  }
</style>
