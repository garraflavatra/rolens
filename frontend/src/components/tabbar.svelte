<script>
  import { createEventDispatcher } from 'svelte';
  import Icon from './icon.svelte';

  export let tabs = [];
  export let selectedKey = '';
  export let canAddTab = false;
  export let compact = true;

  const dispatch = createEventDispatcher();
  let navEl;

  function select(tabKey) {
    selectedKey = tabKey;
    dispatch('select', tabKey);
  }
</script>

<nav class="tabs" class:compact bind:this={navEl}>
  <ul>
    {#each tabs as tab (tab.key)}
      <li class="tab" class:active={tab.key === selectedKey} class:closable={tab.closable}>
        <button class="tab" on:click={() => select(tab.key)}>
          {#if tab.icon} <Icon name={tab.icon} /> {/if}
          <span class="label">{tab.title}</span>
        </button>

        {#if tab.closable}
          <button class="button-small" on:click={() => dispatch('closeTab', tab.key)}>
            <Icon name="x" />
          </button>
        {/if}
      </li>
    {/each}

    {#if canAddTab}
      <li>
        <button class="button-small" on:click={() => dispatch('addTab')}>
          <Icon name="+" />
        </button>
      </li>
    {/if}
  </ul>
</nav>

<style>
  ul {
    overflow-x: auto;
    display: flex;
    list-style: none;
  }

  li {
    display: inline-block;
    position: relative;
  }

  nav.tabs :global(svg) {
    vertical-align: bottom;
  }

  button.tab {
    width: 100%;
    padding: 0.7rem 1rem;
    border: 1px solid #ccc;
    border-right: none;
    cursor: pointer;
    background-color: #fff;
  }
  button.tab:hover {
    background-color: #eee;
  }
  button.tab:active {
    background-color: #ddd;
  }
  li:last-child button.tab {
    border-right: 1px solid #ccc;
  }
  li.active button.tab {
    background-color: var(--selection);
    border-color: var(--primary);
    cursor: not-allowed;
  }

  button.tab .label {
    margin-top: 5px;
    display: inline-block;
  }

  nav.tabs .button-small {
    margin: 12px 0 0 4px;
  }

  li.closable .button-small {
    display: none;
    position: absolute;
    top: 0;
    right: 7px;
  }
  li.closable.active {
    padding-right: 20px;
  }
  li.closable.active .button-small {
    display: block;
  }

  nav.tabs.compact {
    border-bottom: 1px solid #aaa;
  }
  nav.tabs.compact li {
    border-bottom: 2px solid transparent;
  }
  nav.tabs.compact li.tab:hover,
  nav.tabs.compact li.tab.active {
    color: var(--primary);
    border-color: var(--primary);
  }
  nav.tabs.compact button.tab {
    border: none;
    color: inherit;
    background-color: transparent;
    border-width: 0;
  }
  nav.tabs.compact li.active button.tab {
    border-bottom-width: 2px;
  }
</style>
