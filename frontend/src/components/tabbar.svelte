<script>
  import { createEventDispatcher } from 'svelte';
  import Icon from './icon.svelte';

  export let tabs = [];
  export let selectedKey = {};
  export let canAddTab = false;

  const dispatch = createEventDispatcher();

  function select(tabKey) {
    selectedKey = tabKey;
    dispatch('select', tabKey);
  }
</script>

<nav class="tabs">
  <ul>
    {#each tabs as tab (tab.key)}
      <li class:active={tab.key === selectedKey}>
        <button class="tab" on:click={() => select(tab.key)}>{tab.title}</button>
        {#if tab.closable}
          <button class="btn-sm" on:click={() => dispatch('closeTab', tab.key)}>
            <Icon name="x" />
          </button>
        {/if}
      </li>
    {/each}

    {#if canAddTab}
      <li class="tab add">
        <button class="tab" on:click={() => dispatch('addTab')}>
          <Icon name="+" />
        </button>
      </li>
    {/if}
  </ul>
</nav>

<style>
  ul {
    overflow-x: scroll;
    display: flex;
    list-style: none;
  }
  li {
    display: inline-block;
    flex-grow: 1;
    position: relative;
  }

  li.add {
    flex: 0 1;
  }

  .tabs :global(svg) {
    width: 13px;
    height: 13px;
    vertical-align: bottom;
  }
  li.active :global(svg) {
    color: #fff;
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
    color: #fff;
    background-color: #00008b;
    border-color: #00008b;
    cursor: not-allowed;
  }

  .btn-sm {
    position: absolute;
    right: 7px;
    top: 7px;
  }
</style>
