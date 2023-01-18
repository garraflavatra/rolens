<script>
  import Icon from '../../../components/icon.svelte';
  import { input } from '../../../actions';
  import { UpdateItems } from '../../../../wailsjs/go/app/App';
  import CodeExample from '../../../components/code-example.svelte';

  export let collection = {};

  const atomicUpdateOperators = {
    'Fields': {
      $currentDate: 'Current Date',
      $inc: 'Increment',
      $min: 'Min',
      $max: 'Max',
      $mul: 'Multiply',
      $rename: 'Rename',
      $set: 'Set',
      $setOnInsert: 'Set on Insert',
      $unset: 'Unset',
    },
    'Array': {
      $addToSet: 'Add to Set',
      $pop: 'Pop',
      $pull: 'Pull',
      $push: 'Push',
      $pullAll: 'Push All',
    },
    'Modifiers': {
      $each: 'Each',
      $position: 'Position',
      $slice: 'Slice',
      $sort: 'Sort',
    },
    'Bitwise': {
      $bit: 'Bit',
    },
  };

  const form = { query: '{}', parameters: [ {} ] };

  $: code = `db.${collection.key}.${form.many ? 'updateMany' : 'updateOne'}()`;

  async function submitQuery() {
    // form = { query: '{}', parameters: [ {} ] };

    const result = await UpdateItems(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));
    console.log(result);
  }

  function removeParam(index) {
    if (form.parameters.length < 2) {
      return;
    }
    form.parameters.splice(index, 1);
    form.parameters = form.parameters;
  }

  function addParameter(index) {
    if (typeof index !== 'number') {
      form.parameters = [ ...form.parameters, {} ];
    }
    else {
      form.parameters = [
        ...form.parameters.slice(0, index),
        {},
        ...form.parameters.slice(index),
      ];
    }
  }
</script>

<form class="update" on:submit|preventDefault={submitQuery}>
  <CodeExample language="json" {code} />

  <div class="options">
    <label class="field">
      <span class="label">Upsert</span>
      <span class="checkbox">
        <input type="checkbox" bind:checked={form.upsert} />
      </span>
    </label>

    <label class="field">
      <span class="label">Many</span>
      <span class="checkbox">
        <input type="checkbox" bind:checked={form.many} />
      </span>
    </label>
  </div>

  <label class="field">
    <span class="label">Filter</span>
    <input type="text" class="code" bind:value={form.query} use:input={{ json: true, autofocus: true }} placeholder={'{}'} />
  </label>

  <fieldset class="parameters">
    {#each form.parameters as param, index}
      <fieldset class="parameter">
        <label class="field">
          <select bind:value={param.type}>
            {#each Object.entries(atomicUpdateOperators) as [groupName, options]}
              <optgroup label={groupName}>
                {#each Object.entries(options) as [key, label]}
                  <option value={key} disabled={form.parameters.some(p => p.type === key)}>
                    {label}
                  </option>
                {/each}
              </optgroup>
            {/each}
          </select>
        </label>

        <label class="field">
          <input type="text" class="code" bind:value={param.value} placeholder={'{}'} use:input={{ json: true }} />
        </label>

        <button class="btn" on:click={() => addParameter()} type="button">
          <Icon name="+" />
        </button>
        <button class="btn" disabled={form.parameters.length < 2} on:click={() => removeParam(index)} type="button">
          <Icon name="-" />
        </button>
      </fieldset>
    {/each}
  </fieldset>

  <div class="result">
    <div></div>
    <button class="btn" type="submit">Update</button>
  </div>
</form>

<style>
  .update {
    display: grid;
    gap: 0.5rem;
    grid-template: auto auto auto 1fr / 1fr;
  }

  .options {
    display: flex;
    gap: 0.5rem;
  }

  .parameters {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .parameter {
    display: grid;
    grid-template: 1fr / auto 1fr auto auto;
    gap: 0.5rem;
  }

  .result {
    display: flex;
    justify-content: space-between;
  }
</style>
