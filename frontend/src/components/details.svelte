<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import Icon from './icon.svelte';

  export let title = '';
  export let initiallyOpen = false;
  export let deletable = false;

  const dispatch = createEventDispatcher();
  /** @type {HTMLDetailsElement} */ let detailsElement;

  onMount(() => {
    if (initiallyOpen && detailsElement) {
      detailsElement.open = initiallyOpen;
    }
  });
</script>

<details bind:this={detailsElement}>
  <summary>
    <Icon name="chev-d" />
    {title}

    {#if deletable}
      <button class="delete" on:click={() => dispatch('delete')}>
        <Icon name="trash" />
      </button>
    {/if}
  </summary>

  <slot />
</details>

<style>
  details {
    border: 1px solid #aaa;
    border-radius: 4px;
    padding: 0.5em 0.5em 0;
    margin-bottom: 0.5rem;
    background-color: #fff;
  }
  details[open] {
    padding: 0.5em;
  }

  summary {
    font-weight: 700;
    margin: -0.5em -0.5em 0;
    padding: 0.5em;
    display: flex;
    gap: 0.5rem;
    align-items: center;
    cursor: pointer;
  }
  summary :global(svg) {
    width: 14px;
    height: 14px;
  }
  details[open] summary {
    border-bottom: 1px solid #aaa;
    margin-bottom: 0.5em;
  }

  button.delete {
    margin-left: auto;
    color: #c00;
  }
</style>
