<script>
  import { resolveKeypath, setValue } from '../../../../utils';
  import Icon from '../../../../components/icon.svelte';
  import FormInput from './forminput.svelte';

  export let item = {};
  export let view = {};
  export let valid = false;

  const validity = {};
  $: valid = Object.values(validity).every(v => !!v);

  const iconMap = {
    string: 'text',
    int: 'hash',
    long: 'hash',
    uint64: 'hash',
    double: 'hash',
    decimal: 'hash',
    bool: 'toggle-l',
    date: 'cal',
  };

  const keypathProxy = new Proxy(item, {
    get: (item, key) => resolveKeypath(item, key),
    set: (item, key, value) => {
      setValue(item, key, value);
      return true;
    },
  });

  function reset(columnKey) {
    keypathProxy[columnKey] = undefined;
  }
</script>

{#if item && view}
  {#each view?.columns?.filter(c => c.inputType !== 'none') || [] as column}
    <!-- svelte-ignore a11y-label-has-associated-control because FormInput contains one -->
    <label class="column">
      <div class="label">
        <Icon name={iconMap[column.inputType]} />
        <span>
          {column.key}
          {#if column.mandatory}
            <span class="tag" class:invalid={!validity[column.key]}>mandatory</span>
          {/if}
        </span>
      </div>
      <div class="input">
        <FormInput {column} bind:value={keypathProxy[column.key]} bind:valid={validity[column.key]} />
        <button type="button" class="btn" title="Reset value" on:click={() => reset(column.key)} disabled={!keypathProxy[column.key]}>
          <Icon name="reload" />
        </button>
      </div>
    </label>
  {/each}
{/if}

<style>
  .column {
    display: block;
  }
  .column + .column {
    margin-top: 1rem;
  }
  .column .label {
    margin-bottom: 0.5rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .column .label :global(svg) {
    width: 13px;
    height: 13px;
  }

  .input {
    display: grid;
    grid-template: 1fr / 1fr auto;
    gap: 0.5rem;
  }

  .tag {
    display: inline-block;
    background-color: rgba(140, 140, 140, 0.1);
    color: #777;
    text-transform: uppercase;
    font-size: 10px;
    padding: 3px 5px;
    font-weight: 600;
  }
  .tag.invalid {
    background-color: rgba(255, 80, 80, 0.3);
    color: #8d2c2c;
  }
</style>
