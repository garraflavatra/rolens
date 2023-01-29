<script>
  import { isDate } from '../../../../utils';
  import { input } from '../../../../actions';

  export let column = {};
  export let value = undefined;
  export let valid = true;

  const onValid = () => valid = true;
  const onInvalid = () => valid = false;
  const numericTypes = [ 'int', 'long', 'uint64', 'double', 'decimal' ];
  let dateInput;
  let timeInput;
  $: type = column.inputType;
  $: mandatory = column.mandatory;

  $: if (value === undefined) {
    dateInput && (dateInput.value = undefined);
    timeInput && (timeInput.value = undefined);
    mandatory && (valid = false);
  }

  function setDate(event) {
    if (event?.currentTarget?.value) {
      if (!isDate(value)) {
        value = new Date(event.currentTarget.value);
      }
      const date = event.currentTarget.value.split('-').map(n => parseInt(n));
      value.setFullYear(date[0], date[1], date[2]);
    }
  }

  function setTime(event) {
    if (event?.currentTarget?.value) {
      const time = event.currentTarget.value.split(':').map(n => parseInt(n));
      value.setHours?.(time[0], time[1], 0, 0);
    }
  }

  function selectChange() {
    if ((value === undefined) && mandatory) {
      valid = false;
    }
    else {
      valid = true;
    }
  }
</script>

<div class="field {type}">
  {#if type === 'string'}
    <input type="text" bind:value use:input={{ type, onValid, onInvalid, mandatory }} />
  {:else if numericTypes.includes(type)}
    <input type="number" bind:value use:input={{ type, onValid, onInvalid, mandatory }} />
  {:else if type === 'bool'}
    <select bind:value on:change={selectChange}>
      <option value={undefined} disabled={mandatory}>Unset</option>
      <option value={true}>Yes / true</option>
      <option value={false}>No / false</option>
    </select>
  {:else if type === 'date'}
    {@const isNotDate = !isDate(value)}
    <input type="date" bind:this={dateInput} on:input={setDate} use:input />
    <input type="time" bind:this={timeInput} on:input={setTime} disabled={isNotDate} title={isNotDate ? 'Enter a date first' : ''} />
  {/if}
</div>

<style>
  .field.date {
    display: grid;
    grid-template: 1fr / 3fr 1fr;
  }
</style>
