<script>
  import icons from '$lib/icons.json';

  export let name = '';
  export let spin = false;

  if (name === 'loading') {
    spin = true;
  }
</script>

<style>
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

  @keyframes spinning {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  svg.spinning {
    animation: spinning 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
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
>
  {@html icons[name] || ''}
</svg>
