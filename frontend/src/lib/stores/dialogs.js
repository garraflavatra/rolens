import { writable } from "svelte/store";

const { set, subscribe, update } = writable([]);

function newDialog(dialogComponent, data) {
  const host = document.createElement('div');
  host.className = 'dialogoutlet';
  const instance = new dialogComponent({ target: host, props: data });

}

const dialogs = { new: newDialog };

export default dialogs;
