<script>
  import { indentWithTab } from '@codemirror/commands';
  import { javascript } from '@codemirror/lang-javascript';
  import { indentOnInput } from '@codemirror/language';
  import { EditorState } from '@codemirror/state';
  import { EditorView, keymap } from '@codemirror/view';
  import { basicSetup } from 'codemirror';
  import { onMount } from 'svelte';

  export let text = '';

  const editorState = EditorState.create({
    doc: '',
    extensions: [
      basicSetup,
      keymap.of([ indentWithTab, indentOnInput ]),
      javascript(),
      EditorState.tabSize.of(4),
      EditorView.updateListener.of(e => {
        // if (!e.docChanged) {
        //   return;
        // }
        text = e.state.doc.toString();
      }),
    ],
  });

  let editorParent;
  let editor;

  onMount(() => {
    editor = new EditorView({
      parent: editorParent,
      state: editorState,
    });

    editor.dispatch({
      changes: {
        from: 0,
        to: editorState.doc.length,
        insert: text,
      },
    });
  });
</script>

<div bind:this={editorParent} class="editor"></div>

<style>
  .editor {
    width: 100%;
  }

  .editor :global(.cm-editor) {
    overflow: auto;
  }
</style>
