<script>
  import icons from '$lib/icons.json';

  export let name = '';
  export let spin = false;
  export let rotation = 0;

  if (name === 'loading') {
    spin = true;
  }
</script>

<style>
  svg {
    transition: transform 0.25s;
    will-change: transform;
    width: 1.2em;
    height: 1.2em;
  }
  svg.spinning {
    animation: spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  }

  :global(.field) svg {
    width: 1em;
    height: 1em;
    margin-right: 2px;
  }
  :global(.button) svg {
    height: 1em;
    width: auto;
    vertical-align: bottom;
  }
  :global(.blankstate .button) svg {
    height: 1.25em;
    vertical-align: -3px;
    margin-right: 4px;
  }
</style>

<svg
  xmlns="http://www.w3.org/2000/svg"
  width="24"
  height="24"
  viewBox="0 0 24 24"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
  class:spinning={spin}
  style:transform="rotate({rotation}deg)"
>
  {@html icons[name] || ''}
</svg>
