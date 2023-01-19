<script>
  export let data;
  export let depth = 0;
  export let readonly = false;
  export let level = 0;
  export let last = true;
  export let draggable = false;
  export let kp = '';

  const collapsedSymbol = '...';
  const getType = i => {
    if (i === null) {
      return 'null';
    }
    return typeof i;
  };

  let displayOnly = true;
  let items;
  let isArray;
  let openBracket;
  let closeBracket;
  $: {
    items = getType(data) === 'object' ? Object.keys(data) : [];
    isArray = Array.isArray(data);
    openBracket = isArray ? '[' : '{';
    closeBracket = isArray ? ']' : '}';
  }

  let collapsed;
  $: collapsed = depth < level;

  const format = i => {
    switch (getType(i)) {
      case 'string':
        return `${i}`;
      case 'function':
        return 'f () {...}';
      case 'symbol':
        return i.toString();
      default:
        return i;
    }
  };
  const clicked = e => {
    if (e.shiftKey) {
      if (depth == 0) {
        depth = 999;
      }
      else {
        depth = 0;
      }
    }
    collapsed = !collapsed;
  };

  let invalid = false;
  let dbg;

  function json2data() {
    try {
      data = JSON.parse(dbg.value);
      invalid = false;
    }
    catch {
      invalid = true;
      if (dbg.value.trim == '') {
        data = {};
      }
    }
  }

  function dragstart(e, keypath, value) {
    console.log('kp:', keypath);
    const item = {};
    item[keypath] = value;
    e.dataTransfer.setData('text/plain', JSON.stringify(item));
  }

  function onKeydown(event) {
    const save = (event.key === 's') && (event.metaKey || event.ctrlKey);
    if (!save) {
      event.stopPropagation();
    }
  }
</script>

{#if displayOnly}
  {#if items.length}
    <span class:root={level == 0} class:hidden={collapsed}>
      {#if draggable && isArray}
        <span on:dragstart={e => dragstart(e, kp, data)} draggable="true" class="bracket" on:click={clicked} tabindex="0">{openBracket}</span>
      {:else}
        <span class="bracket" on:click={clicked} tabindex="0">{openBracket}</span>
      {/if}
      <ul on:dblclick={() => (readonly ? displayOnly = true : displayOnly = false)} >
        {#each items as i, idx}
          <li>
            {#if !isArray}
              {#if draggable}
                <span on:dragstart={e => dragstart(e, kp ? kp + '.' + i : i, data[i])} draggable="true" class="key">{i}:</span>
              {:else}
                <span class="key">{i}:</span>
              {/if}
            {/if}
            {#if getType(data[i]) === 'object'}
              <svelte:self {readonly} {draggable} kp={kp ? kp + '.' + i : i} data={data[i]} {depth} level={level + 1} last={idx === items.length - 1} />
            {:else}
              {#if draggable}
                <span on:dragstart={e => dragstart(e, kp ? kp + '.' + i : i, data[i])} draggable="true" class="val {getType(data[i])}">{format(data[i])}</span>{#if idx < items.length - 1}<span draggable class="comma">,</span>{/if}
              {:else}
                <span class="val {getType(data[i])}">{format(data[i])}</span>{#if idx < items.length - 1}<span class="comma">,</span>{/if}
              {/if}
            {/if}
          </li>
        {/each}
      </ul>
      <span class="bracket" on:click={clicked} tabindex="0">{closeBracket}</span>{#if !last}<span
          class="comma">,</span>
    {/if}
    </span>
    <span style="padding: {level == 0 ? 10 : 0}px;" class="bracket" class:hidden={!collapsed} on:click={clicked} tabindex="0">{openBracket}{collapsedSymbol}{closeBracket}</span>{#if !last && collapsed}<span class="comma">,</span>{/if}
  {:else}
    {@html isArray ? '[]' : '{}'}
  {/if}
{:else}
  <textarea on:keydown="{onKeydown}" class="debug" spellcheck="false" bind:this={dbg} class:invalid on:input={json2data}>{JSON.stringify(data, null, 2)}</textarea>
{/if}

<style>
  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    font-family: inherit;
    font-size: inherit;
    padding-left: var(--nodePaddingLeft, 1rem);
    border-left: var(--nodeBorderLeft, 1px dashed #d0d0f0);
    color: var(--nodeColor, #666);
  }
  li {
    white-space: nowrap;
  }
  .root {
    font-family: menlo, monospace;
    font-size: 90%;
    overflow: auto;
  }
  .hidden {
    display: none;
  }
  .bracket {
    cursor: pointer;
  }
  .bracket:hover {
    background: var(--bracketHoverBackground, #d1d5db);
  }
  .comma {
    color: var(--nodeColor, #374151);
    opacity: 0.5;
  }
  .val {
    color: var(--leafDefaultColor, #9ca3af);
    white-space: nowrap;
  }
  .val[draggable] {
    cursor: move;
  }
  .val.string {
    color: var(--leafStringColor, #000);
  }
  .val.string:before {
    content: "'";
    opacity: 0.4;
  }
  .val.string:after {
    content: "'";
    opacity: 0.4;
  }
  .val.number {
    color: var(--leafNumberColor, #d97706);
  }
  .val.boolean {
    color: var(--leafBooleanColor, #3994dd);
  }
  .key.draggable {
    cursor: move;
  }

  textarea.debug {
    font-family: menlo, monospace;
    padding: 10px;
    flex: 1 0;
    white-space: pre;
    font-size: 90%;
    border: none;
    margin: 0;
    height: 100%;
    outline: none;
    line-height: 1.5;
    resize: none;
  }

  textarea.invalid {
    background: #ffe3e3;
    color: #b30202 !important;
  }
</style>
