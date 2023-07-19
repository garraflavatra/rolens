import { AskConfirmation } from '$wails/go/app/App';
import InputDialog from '../dialogs/input.svelte';

function newDialog(dialogComponent, data = {}) {
  const outlet = document.createElement('div');
  outlet.className = 'dialogoutlet';
  document.getElementById('dialogoutlets').appendChild(outlet);

  const instance = new dialogComponent({ target: outlet, props: data });

  instance.$close = function() {
    instance.$destroy();
    outlet.remove();
  };

  instance.$on('close', instance.$close);

  return instance;
}

function enterText(title = '', description = '', value = '') {
  const instance = newDialog(InputDialog, { title, description, value });

  return new Promise(resolve => {
    instance.$on('submit', event => {
      instance.$close();
      resolve(event.detail.value);
    });
    instance.$on('close', () => resolve(undefined));
  });
}

function confirm(message = '') {
  return AskConfirmation(message);
}

const dialogs = { new: newDialog, enterText, confirm };

export default dialogs;
