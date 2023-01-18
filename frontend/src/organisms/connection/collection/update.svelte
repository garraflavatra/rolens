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
  const allOperators = Object.values(atomicUpdateOperators).map(Object.keys).flat();

  const form = { query: '{}', parameters: [ { type: '$set' } ] };
  let updatedCount;
  $: code = buildCode(form);

  function buildCode(form) {
    const method = form.many ? 'updateMany' : 'updateOne';

    let operation = '{ ' + form.parameters.filter(p => p.type).map(p => `${p.type}: ${p.value || '{}'}`).join(', ') + ' }';
    if (operation === '{  }') {
      operation = '';
    }

    const code = `db.${collection.key}.${method}(${form.query || '{}'}${operation ? ', ' + operation : ''});`;
    return code;
  }

  async function submitQuery() {
    updatedCount = await UpdateItems(collection.hostKey, collection.dbKey, collection.key, JSON.stringify(form));
  }

  function removeParam(index) {
    if (form.parameters.length < 2) {
      return;
    }
    form.parameters.splice(index, 1);
    form.parameters = form.parameters;
  }

  function addParameter(index) {
    const usedOperators = form.parameters.map(p => p.type);
    const operator = allOperators.find(o => !usedOperators.includes(o));

    if (!operator) {
      return;
    }

    const newItem = { type: operator };
    if (typeof index !== 'number') {
      form.parameters = [ ...form.parameters, newItem ];
    }
    else {
      form.parameters = [
        ...form.parameters.slice(0, index),
        newItem,
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

    {#key updatedCount}
      {#if typeof updatedCount === 'number'}
        <div class="flash-green">
          Updated {updatedCount} item{updatedCount === 1 ? '' : 's'}
        </div>
      {/if}
    {/key}

    <button class="btn" type="submit">Update</button>
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

        <button class="btn" disabled={form.parameters.length >= allOperators.length} on:click={() => addParameter()} type="button">
          <Icon name="+" />
        </button>
        <button class="btn" disabled={form.parameters.length < 2} on:click={() => removeParam(index)} type="button">
          <Icon name="-" />
        </button>
      </fieldset>
    {/each}
  </fieldset>
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
    align-items: center;
  }
  .options button {
    margin-left: auto;
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
</style>
