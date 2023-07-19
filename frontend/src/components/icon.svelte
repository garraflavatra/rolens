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
  @keyframes spinning {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  svg {
    transition: transform 0.25s;
    will-change: transform;
  }
  svg.spinning {
    animation: spinning 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
  }

  :global(.field) svg {
    width: 13px;
    height: 13px;
    margin-right: 2px;
  }
  :global(.button) svg {
    height: 13px;
    width: auto;
    vertical-align: bottom;
  }
  :global(.blankstate .button) svg {
    height: 17px;
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
