import './reset.css';
import './style.css';
import './loading.css';
import App from './app.svelte';
import { LogError } from '../wailsjs/runtime/runtime';

window.addEventListener('unhandledrejection', event => {
  LogError('Unhandled Rejection in JS! Reason:');
  LogError(String(event.reason));
});

const app = new App({ target: document.getElementById('app') });

export default app;
