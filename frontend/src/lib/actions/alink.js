import { BrowserOpenURL } from '$wails/runtime/runtime.js';

export default function alink(node) {
  node.addEventListener('click', e => {
    e.preventDefault();
    BrowserOpenURL(node.href);
  });
}
