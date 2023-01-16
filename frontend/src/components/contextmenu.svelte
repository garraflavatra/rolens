<script>
  import { createEventDispatcher } from 'svelte';

  export let items = undefined;
  export let position = undefined;

  const dispatch = createEventDispatcher();

  function close() {
    dispatch('close');
  }

  function click(fn) {
    fn?.();
    close();
  }
</script>

{#if items && position}
  <div class="backdrop" on:pointerdown={close}></div>
  <ul class="contextmenu" role="" style:left="{position[0]}px" style:top="{position[1]}px">
    {#each items as item}
      {#if item.separator}
        <hr />
      {:else}
        <li>
          <button class="item" on:click={() => click(item.fn)}>
            {item.label}
          </button>
        </li>
      {/if}
    {/each}
  </ul>
{/if}

<style>
  .backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }

  .contextmenu {
    position: fixed;
    background-color: rgba(230, 230, 230, 0.7);
    -webkit-backdrop-filter: blur(30px);
    backdrop-filter: blur(30px);
    border-radius: 10px;
    padding: 5px;
    box-shadow: 0px 3px 5px rgba(0, 0, 0, 0.2);
  }

  button {
    padding: 5px;
    border-radius: 5px;
    width: 100%;
    text-align: left;
  }
  button:hover {
    background-color: #00008b;
    color: #fff;
  }
</style>
