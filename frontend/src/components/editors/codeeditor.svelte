<script>
  import { indentWithTab } from '@codemirror/commands';
  import { indentOnInput } from '@codemirror/language';
  import { EditorState } from '@codemirror/state';
  import { EditorView, keymap } from '@codemirror/view';
  import { basicSetup } from 'codemirror';
  import { createEventDispatcher, onMount } from 'svelte';

  export let text = '';
  export let editor = undefined;
  export let extensions = [];

  const dispatch = createEventDispatcher();
  let editorParent;

  const editorState = EditorState.create({
    doc: '',
    extensions: [
      basicSetup,
      keymap.of([ indentWithTab, indentOnInput ]),
      EditorState.tabSize.of(4),
      EditorView.updateListener.of(e => {
        if (!e.docChanged) {
          return;
        }
        text = e.state.doc.toString();
        dispatch('updated', { text });
      }),
      ...extensions,
    ],
  });

  onMount(() => {
    editor = new EditorView({
      parent: editorParent,
      state: editorState,
    });

    dispatch('inited', { editor });
  });
</script>

<div bind:this={editorParent} class="editor"></div>

<style>
  .editor {
    width: 100%;
    background-color: #fff;
    border-radius: var(--radius);
    overflow: hidden;
  }

  .editor :global(.cm-editor) {
    overflow: auto;
    height: 100%;
  }
</style>
