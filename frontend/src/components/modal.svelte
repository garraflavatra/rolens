<script context="module">
  let numberOfModalsOpen = 0;
</script>

<script>
  import { createEventDispatcher } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import { Beep } from '$wails/go/ui/UI.js';
  import Icon from './icon.svelte';

  export let show = true;
  export let title = undefined;
  export let contentPadding = true;
  export let width = '80vw';
  export let overflow = true;

  const dispatch = createEventDispatcher();
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
      event.preventDefault();
      close();
    }
  }

  function close() {
    show = false;
    dispatch('close');
  }
</script>

<svelte:window on:keydown={keydown} />

{#if show}
  <div class="modal outer" on:pointerdown|self={Beep} transition:fade={{ duration: 200 }}>
    <div class="inner" style:max-width={width || '80vw'} transition:fly={{ y: 20, duration: 200 }}>
      {#if title}
        <header>
          <div class="title">{title}</div>
          <button class="button close" on:click={close} title="close" type="button">
            <Icon name="x" />
          </button>
        </header>
      {/if}

      <div class="slot content" class:padded={contentPadding} class:overflow> <slot /> </div>

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
    padding: 2rem;
    --wails-draggable: drag;
  }

  .inner {
    max-height: 80vh;
    background-color: #fff;
    margin: auto;
    width: 100%;
    border-radius: var(--radius);
    display: flex;
    flex-flow: column;
    cursor: auto;
    overflow: hidden;
    --wails-draggable: nodrag;
  }
  .inner > :global(*:first-child) {
    margin-top: 0;
  }

  .inner > :global(*:last-child) {
    margin-bottom: 0;
  }

  header {
    display: flex;
    align-items: center;
    padding: 0.75rem;
    background-color: #eee;
  }
  header .title {
    font-size: 1.5rem;
  }

  .close {
    margin-left: auto;
  }

  .content {
    max-height: 100%;
  }
  .content.overflow {
    overflow-y: auto;
  }
  .content.padded {
    padding: 0.75rem;
  }
  header + .content.padded {
    border-top: 1px solid #ccc;
  }

  footer {
    padding: 0.75rem;
    border-top: 1px solid #ccc;
  }
</style>
