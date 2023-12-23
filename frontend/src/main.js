import './styles/reset.css';
import './styles/style.css';

import { LogError } from '$wails/runtime/runtime.js';
import App from './app.svelte';

window.addEventListener('unhandledrejection', event => {
  LogError('Unhandled JS rejection: ' + event.reason);
});

const app = new App({ target: document.getElementById('app') });

export default app;
