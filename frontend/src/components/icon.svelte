<script>
  export let name = '';
</script>

<style>
  :global(.field) svg {
    width: 13px;
    height: 13px;
    margin-right: 2px;
  }

  :global(.btn) svg {
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
  class:spinning={name === 'loading'}
>
  {#if name === 'radio'}
    <circle cx="12" cy="12" r="2"></circle><path d="M16.24 7.76a6 6 0 0 1 0 8.49m-8.48-.01a6 6 0 0 1 0-8.49m11.31-2.82a10 10 0 0 1 0 14.14m-14.14 0a10 10 0 0 1 0-14.14"></path>
  {:else if name === 'chev-l'}
    <polyline points="15 18 9 12 15 6"></polyline>
  {:else if name === 'chev-r'}
    <polyline points="9 18 15 12 9 6"></polyline>
  {:else if name === 'chev-d'}
    <polyline points="6 9 12 15 18 9"></polyline>
  {:else if name === 'chev-u'}
    <path d="m18 15-6-6-6 6" />
  {:else if name === 'chevs-l'}
    <path d="m11 17-5-5 5-5M18 17l-5-5 5-5" />
  {:else if name === 'chevs-r'}
    <path d="m13 17 5-5-5-5M6 17l5-5-5-5" />
  {:else if name === 'arr-d'}
    <path d="M12 5v14M19 12l-7 7-7-7" />
  {:else if name === 'db'}
    <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
  {:else if name === 'x'}
    <line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line>
  {:else if name === '+'}
    <line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line>
  {:else if name === '-'}
    <line x1="5" y1="12" x2="19" y2="12"></line>
  {:else if name === 'reload'}
    <polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
  {:else if name === 'clipboard'}
    <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8"
      y="2"
      width="8"
      height="4"
      rx="1"
      ry="1"></rect>
  {:else if name === 'edit'}
    <path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
  {:else if name === 'check'}
    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline>
  {:else if name === 'list'}
    <path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01" />
  {:else if name === 'table'}
    <path d="M9 3H5a2 2 0 0 0-2 2v4m6-6h10a2 2 0 0 1 2 2v4M9 3v18m0 0h10a2 2 0 0 0 2-2V9M9 21H5a2 2 0 0 1-2-2V9m0 0h18" />
  {:else if name === 'form'}
    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" /><path d="M14 2v6h6M16 13H8M16 17H8M10 9H8" />
  {:else if name === 'cog'}
    <circle cx="12" cy="12" r="3" /><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
  {:else if name === 'zap'}
    <path d="M13 2 3 14h9l-1 8 10-12h-9l1-8z" />
  {:else if name === 'server'}
    <rect x="2"
      y="2"
      width="20"
      height="8"
      rx="2"
      ry="2" /><rect x="2"
      y="14"
      width="20"
      height="8"
      rx="2"
      ry="2" /><path d="M6 6h.01M6 18h.01" />
  {:else if name === 'text'}
    <path d="M4 7V4h16v3M9 20h6M12 4v16" />
  {:else if name === 'hash'}
    <path d="M4 9h16M4 15h16M10 3 8 21M16 3l-2 18" />
  {:else if name === 'toggle-l'}
    <rect x="1"
      y="5"
      width="22"
      height="14"
      rx="7"
      ry="7" /><circle cx="8" cy="12" r="3" />
  {:else if name === 'cal'}
    <rect x="3"
      y="4"
      width="18"
      height="18"
      rx="2"
      ry="2" /><path d="M16 2v4M8 2v4M3 10h18" />
  {:else if name === 'code'}
    <path d="m16 18 6-6-6-6M8 6l-6 6 6 6" />
  {:else if name === 'target'}
    <circle cx="12" cy="12" r="10" /><circle cx="12" cy="12" r="6" /><circle cx="12" cy="12" r="2" />
  {:else if name === 'trash'}
    <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M10 11v6M14 11v6" />
  {:else if name === 'anchor'}
    <circle cx="12" cy="5" r="3" /><path d="M12 22V8M5 12H2a10 10 0 0 0 20 0h-3" />
  {:else if name === 'o'}
    <circle cx="12" cy="12" r="10" />
  {:else if name === 'info'}
    <circle cx="12" cy="12" r="10" /><path d="M12 16v-4M12 8h.01" />
  {:else if name === 'play'}
    <path d="m5 3 14 9-14 9V3z" />
  {:else if name === 'upload'}
    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M17 8l-5-5-5 5M12 3v12" />
  {:else if name === 'save'}
    <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" /><path d="M17 21v-8H7v8M7 3v5h8" />
  {:else if name === 're'}
    <path d="m17 1 4 4-4 4" /><path d="M3 11V9a4 4 0 0 1 4-4h14M7 23l-4-4 4-4" /><path d="M21 13v2a4 4 0 0 1-4 4H3" />
  {:else if name === 'chart'}
    <path d="M18 20V10M12 20V4M6 20v-6" />
  {:else if name === '?'}
    <circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line>
  {:else if name === '!'}
    <path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0zM12 9v4M12 17h.01" />
  {:else if name === 'loading'}
    <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83" />
  {/if}
</svg>
