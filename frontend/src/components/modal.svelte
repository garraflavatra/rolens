<script>
  import { fade, fly } from 'svelte/transition';
  import Icon from './icon.svelte';

  export let show = false;
  export let title = undefined;
  export let contentPadding = true;
</script>

{#if show}
  <div class="modal outer" on:mousedown|self={() => show = false} transition:fade>
    <div class="inner" transition:fly={{ y: 100 }}>
      <header>
        {#if title}
          <div class="title">{title}</div>
        {/if}
        <button class="btn close" on:click={() => show = false} title="close">
          <Icon name="x" />
        </button>
      </header>
      <div class="slot content" class:padded={contentPadding}> <slot /> </div>
      {#if $$slots.footerLeft || $$slots.footerRight}
        <footer>
          <slot name="footer" />
        </footer>
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
    padding-top: 1rem;
    cursor: pointer;
  }

  .inner {
    max-width: 80vw;
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
