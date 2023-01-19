<script context="module">
  let numberOfModalsOpen = 0;
</script>

<script>
  import { fade, fly } from 'svelte/transition';
  import Icon from './icon.svelte';

  export let show = false;
  export let title = undefined;
  export let contentPadding = true;
  export let width = '80vw';

  const level = numberOfModalsOpen + 1;
  let isNew = true;

  $: if (show) {
    numberOfModalsOpen++;
  }
  else if (!isNew) {
    numberOfModalsOpen--;
  }
  else {
    isNew = false;
  }

  function keydown(event) {
    if ((event.key === 'Escape') && (level === numberOfModalsOpen)) {
      show = false;
    }
  }
</script>

<svelte:window on:keydown|preventDefault={keydown} />

{#if show}
  <div class="modal outer" on:mousedown|self={() => show = false} transition:fade>
    <div class="inner" style:max-width={width || '80vw'} transition:fly={{ y: 100 }}>
      {#if title}
        <header>
          <div class="title">{title}</div>
          <button class="btn close" on:click={() => show = false} title="close" type="button">
            <Icon name="x" />
          </button>
        </header>
      {/if}

      <div class="slot content" class:padded={contentPadding}> <slot /> </div>

      {#if $$slots.footer}
        <footer> <slot name="footer" /> </footer>
      {/if}
    </div>
  </div>
{/if}

<style>
  .outer {
    position: fixed;
    display: flex;
    z-index: 100;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.5);
    margin: 0;
    padding-top: 50px;
    cursor: pointer;
  }
  :global(#root.platform-darwin) .outer {
    margin-top: var(--darwin-titlebar-height);
  }

  .inner {
    max-height: 80vh;
    background-color: #fff;
    margin-left: auto;
    margin-right: auto;
    margin-bottom: auto;
    width: 100%;
    border-radius: 10px;
    overflow: hidden;
    display: flex;
    flex-flow: column;
    cursor: auto;
  }
  .inner > :global(*:first-child) {
    margin-top: 0;
  }

  .inner > :global(*:last-child) {
    margin-bottom: 0;
  }

  header {
    border-bottom: 1px solid #ccc;
    display: flex;
    align-items: center;
    padding: 1rem;
  }
  header .title {
    font-size: 1.5rem;
  }

  .close {
    margin-left: auto;
  }

  .content {
    overflow-y: auto;
    max-height: 100%;
  }
  .content.padded {
    padding: 1rem;
  }

  footer {
    padding: 1rem;
    border-top: 1px solid #ccc;
  }

  @media (max-width: 600px) {
    .outer {
      padding: 0;
    }

    .inner {
      max-width: 100%;
      max-height: 100%;
      width: 100%;
      margin-top: auto;
      margin-bottom: 0;
    }
  }
</style>
