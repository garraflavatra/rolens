<script>
  import { createEventDispatcher } from 'svelte';

  export let items = undefined;
  export let position = undefined;

  const dispatch = createEventDispatcher();
  const buttons = [];
  let selected = -1;

  $: if (items && position) {
    selected = 0;
  }

  function close() {
    dispatch('close');
    selected = -1;
  }

  function click(fn) {
    fn?.();
    close();
  }

  function keydown(evt) {
    if (evt.key === 'Escape') {
      close();
      return;
    }
    else if (!items?.length) {
      return;
    }

    let delta = 0;
    (evt.key === 'ArrowDown') && delta++;
    (evt.key === 'ArrowUp') && delta--;

    selected += delta;
    if (selected >= items.length) {
      selected = 0;
    }
    else if (items[selected].separator) {
      selected += delta;
    }

    buttons[selected]?.focus?.();
  }
</script>

<svelte:window on:keydown={keydown} />

{#if items && position}
  <div class="backdrop" on:pointerdown={close}></div>
  <nav>
    <ul class="contextmenu" role="" style:left="{position[0]}px" style:top="{position[1]}px">
      {#each items as item, index}
        {#if item.separator}
          <hr />
        {:else}
          <li>
            <button
              class="item"
              class:selected={selected === index}
              bind:this={buttons[index]}
              on:mouseenter={() => selected = index}
              on:mouseleave={() => selected = -1}
              on:click={() => click(item.fn)}
            >
              {item.label}
            </button>
          </li>
        {/if}
      {/each}
    </ul>
  </nav>
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
  button.selected {
    background-color: #00008b;
    color: #fff;
  }

  hr {
    border: none;
    border-top: 1px solid #aaa;
    margin: 5px;
  }
</style>
