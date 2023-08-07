<script>
  import input from '$lib/actions/input.js';
  import { canBeObjectId, numericInputTypes } from '$lib/mongo.js';
  import { ObjectId } from 'bson';
  import { onMount } from 'svelte';
  import Datepicker from './datepicker.svelte';
  import Icon from '../icon.svelte';

  export let column = {};
  export let value = undefined;
  export let valid = true;
  export let autofocus = false;

  const onValid = () => valid = true;
  const onInvalid = () => valid = false;
  let objectIdInput;
  let showDatepicker;
  let selectInput;
  $: type = column.inputType;
  $: mandatory = column.mandatory;

  $: if ((value === undefined) && mandatory) {
    valid = false;
  }

  function markInputValid(input) {
    input.setCustomValidity('');
    input.reportValidity();
    input.classList.remove('invalid');
    valid = true;
  }

  function setObjectId(event) {
    if (canBeObjectId(valid)) {
      value = new ObjectId(event.currentTarget?.value);
    }
  }

  function generateObjectId() {
    value = new ObjectId();
    objectIdInput.value = value.toString();
    markInputValid(objectIdInput);
    objectIdInput.disabled = true;
  }

  function editObjectId() {
    if (!objectIdInput) {
      return;
    }
    objectIdInput.disabled = false;
    objectIdInput.focus();
  }

  function selectChange() {
    if ((value === undefined) && mandatory) {
      valid = false;
    }
    else {
      valid = true;
    }
  }

  onMount(() => {
    if (autofocus && selectInput) {
      selectInput.focus();
    }
  });
</script>

<div class="forminput {type}">
  <div class="field">
    {#if type === 'string'}
      <input
        type="text"
        bind:value
        use:input={{ type, onValid, onInvalid, mandatory, autofocus }}
        autocomplete="off"
        spellcheck="false" />
    {:else if type === 'objectid'}
      <input
        type="text"
        bind:this={objectIdInput}
        on:input={setObjectId}
        use:input={{ type, onValid, onInvalid, mandatory, autofocus }}
      />
    {:else if numericInputTypes.includes(type)}
      <input
        type="number"
        bind:value
        use:input={{ type, onValid, onInvalid, mandatory, autofocus }}
      />
    {:else if type === 'bool'}
      <select bind:value on:change={selectChange} bind:this={selectInput}>
        <option value={undefined} disabled={mandatory}>Unset</option>
        <option value={true}>Yes / true</option>
        <option value={false}>No / false</option>
      </select>
    {:else if type === 'date'}
      <input type="text" readonly value={value?.toString() || '...'} on:focus={() => showDatepicker = true} />
    {/if}
  </div>

  <div class="actions">
    {#if type === 'objectid'}
      {#if objectIdInput?.disabled}
        <button class="button-small" type="button" title="Edit object id" on:click={editObjectId}>
          <Icon name="edit" />
        </button>
      {/if}
      <button
        class="button-small"
        type="button"
        title="Generate random object id"
        on:click={generateObjectId}
      >
        <Icon name="reload" />
      </button>
    {:else if type === 'date'}
      <button
        class="button-small"
        type="button"
        title="Edit date"
        on:click={() => showDatepicker = true}
      >
        <Icon name="edit" />
      </button>
      <button
        class="button-small"
        type="button"
        title="Set date to now"
        on:click={() => value = new Date()}
      >
        <Icon name="o" />
      </button>
    {/if}
    <button
      class="button-small"
      type="button"
      title="Reset field to default value"
      on:click={() => value = undefined}
    >
      <Icon name="trash" />
    </button>
  </div>
</div>

{#if type === 'date'}
  <Datepicker bind:value bind:show={showDatepicker} />
{/if}

<style>
  .forminput {
    position: relative;
  }

  .forminput.date input {
    cursor: pointer;
  }

  .actions {
    display: flex;
    position: absolute;
    right: 5px;
    top: 5px;
    background-color: #fff;
  }
  .actions button:last-child {
    border-radius: 2px 6px 6px 2px;
  }
</style>
