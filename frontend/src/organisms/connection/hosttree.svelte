<script>
  import Grid from '$components/grid/grid.svelte';
  import hostTree from '$lib/stores/hosttree.js';

  export let path = [];
</script>

<Grid
  striped={false}
  columns={[
    { key: 'name' },
    { key: 'count', right: true },
  ]}
  items={Object.values($hostTree || {}).map(host => {
    return {
      id: host.key,
      name: host.name,
      loading: host.loading,
      icon: 'server',

      children: Object.values(host.databases || {})
        .sort((a, b) => a.key.localeCompare(b))
        .map(database => {
          return {
            id: database.key,
            name: database.key,
            loading: database.loading,
            count: Object.keys(database.collections || {}).length || '',
            icon: 'db',

            children: Object.values(database.collections)
              .sort((a, b) => a.key.localeCompare(b))
              .map(collection => {
                return {
                  id: collection.key,
                  name: collection.key,
                  icon: 'list',
                  menu: [
                    { label: 'Export collection…', fn: collection.export },
                    { label: 'Dump collection (BSON via mongodump)…', fn: collection.dump },
                    { separator: true },
                    { label: 'Rename collection…', fn: collection.rename },
                    { label: 'Duplicate collection…', fn: collection.duplicate },
                    { separator: true },
                    { label: 'Truncate collection…', fn: collection.truncate },
                    { label: 'Drop collection…', fn: collection.drop },
                    { separator: true },
                    { label: 'New collection…', fn: database.newCollection },
                  ],
                };
              }) || [],

            menu: [
              { label: 'Dump database (BSON via mongodump)…', fn: database.dump },
              { label: 'Drop database…', fn: database.drop },
              { separator: true },
              { label: 'New database…', fn: host.newDatabase },
              { label: 'New collection…', fn: database.newCollection },
            ],
          };
        }),

      menu: [
        { label: 'New database…', fn: host.newDatabase },
        { separator: true },
        { label: `Edit host ${host.name}…`, fn: host.edit },
        { label: 'Remove host…', fn: host.remove },
      ],
    };
  })}
  on:select={e => {
    let level;
    ({ path, level } = e.detail);

    switch (level) {
      case 0: return $hostTree[path[0]].open();
      case 1: return $hostTree[path[0]].databases[path[1]].open();
      case 2: return $hostTree[path[0]].databases[path[1]].collections[path[2]].open();
    }
  }}
/>

<!--
  actions={[
    { icon: 'reload', fn: reload },
    { icon: '+', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'new')) },
    { icon: 'edit', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'edit')), disabled: !activeHostKey },
    { icon: '-', fn: evt => contextMenu.show(evt, buildMenu(activeHostKey, activeDbKey, activeCollKey, 'drop')), disabled: !activeDbKey },
  ]}
-->
