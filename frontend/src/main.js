import './styles/reset.css';
import './styles/style.css';

import { EventsOn, LogError } from '$wails/runtime/runtime.js';
import dialogs from '$lib/dialogs.js';

import App from './app.svelte';
import AboutDialog from './dialogs/about.svelte';
import SettingsDialog from './dialogs/settings/index.svelte';

window.addEventListener('unhandledrejection', event => {
  LogError('Unhandled JS rejection: ' + event.reason);
});

EventsOn('global.about', () => dialogs.new(AboutDialog));
EventsOn('global.settings', () => dialogs.new(SettingsDialog));

const app = new App({ target: document.getElementById('app') });

export default app;
