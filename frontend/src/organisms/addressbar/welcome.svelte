<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function createHost() {
    dispatch('createHost');
  }
</script>

<div class="welcome">
  <p class="title">Welcome to Rolens!</p>
  <button class="btn" on:click={createHost}>Create your first host</button>
</div>

<style>
  .welcome {
    text-align: center;
    padding: 2rem;
  }

  .title {
    font-weight: 600;
    font-size: 1.5rem;
  }

  .btn {
    margin-top: 2rem;
  }
</style>
