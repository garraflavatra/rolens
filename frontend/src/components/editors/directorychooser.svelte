<script>
  import { ChooseDirectory } from '$wails/go/app/App.js';

  export let value = '';
  export let id = '';
  export let title = 'Choose a directory';

  async function selectDir() {
    value = await ChooseDirectory(title) || value;
  }
</script>

<input
  type="text"
  on:pointerdown={selectDir}
  readonly
  {id}
  {value}
/>

<style>
  input {
    cursor: pointer;
  }
</style>
