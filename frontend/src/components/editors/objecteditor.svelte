<script>
  import { javascript } from '@codemirror/lang-javascript';
  import { onMount } from 'svelte';
  import CodeEditor from './codeeditor.svelte';

  export let text = '';
  export let editor = undefined;
  export let readonly = false;

  const extensions = [
    javascript(),
  ];

  onMount(() => {
    editor.dispatch({
      changes: {
        from: 0,
        to: editor.state.doc.length,
        insert: text,
      },
    });
  });
</script>

<CodeEditor
  bind:editor
  bind:text
  on:inited
  on:updated
  {extensions}
  {readonly}
/>
