import { writable } from "svelte/store";

const { set, subscribe, update } = writable([]);

const dialogs = {
  subscribe,

  create: function(dialogComponent, data) {
    const host = document.createElement('div');
    host.className = 'dialogoutlet';
    const instance = new dialogComponent({ target: host, props: data });
    //...
  },
};

export default dialogs;
