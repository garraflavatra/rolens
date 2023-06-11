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

const dialogs = { new: newDialog };

export default dialogs;
