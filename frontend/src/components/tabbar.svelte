<script>
  import { createEventDispatcher } from 'svelte';

  export let tabs = [];
  export let selectedKey = {};

  const dispatch = createEventDispatcher();

  function select(tabKey) {
    selectedKey = tabKey;
    dispatch('select', tabKey);
  }
</script>

<nav class="tabs">
  <ul>
    {#each tabs as tab (tab.key)}
      <li class="tab" class:active={tab.key === selectedKey}>
        <button on:click={() => select(tab.key)}>{tab.title}</button>
      </li>
    {/each}
  </ul>
</nav>

<style>
  .tabs {
    border-bottom: 1px solid #ccc;
    padding: 0 0.5rem;
  }
  .tabs ul {
    overflow-x: scroll;
    display: flex;
    gap: 0.5rem;
    list-style: none;
  }
  .tabs li {
    display: inline-block;
    flex-grow: 1;
  }
  .tabs li button {
    width: 100%;
    padding: 0.7rem 1rem;
    border: 1px solid #ccc;
    border-bottom: none;
    border-radius: 5px 5px 0 0;
    cursor: pointer;
    background-color: #fff;
  }
  .tabs li.active button {
    color: #fff;
    background-color: #00008b;
    border-color: #00008b;
    cursor: not-allowed;
  }
</style>
