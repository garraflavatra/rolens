import './styles/loading.css';
import './styles/reset.css';
import './styles/style.css';

import { LogError } from '$wails/runtime';
import App from './app.svelte';

window.addEventListener('unhandledrejection', event => {
  LogError('Unhandled Rejection in JS! Reason:');
  LogError(String(event.reason));
});

// @ts-ignore Argument IS correct.
const app = new App({ target: document.getElementById('app') });

export default app;
