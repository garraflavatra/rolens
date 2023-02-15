import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { dirname } from 'path';
import { fileURLToPath } from 'url';

const currentDir = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [ svelte() ],
  resolve: {
    alias: {
      '$components': currentDir + '/src/components',
      '$organisms': currentDir + '/src/organisms',
      '$wails': currentDir + '/wailsjs',
      '$lib': currentDir + '/src/lib',
    },
  },
});
