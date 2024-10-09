<template>
  <div id="app">
    <D3NetworkGraph :nodes="nodes" :links="links" :options="options" />
  </div>
</template>
<script>

export default {
  name: 'NetworkChart',
  data() {
    return {
      nodes: [
        { id: 1, name: "192.168.2.1" },
        { id: 2, name: "192.168.2.128" },
        { id: 3, name: "192.168.2.2", class: "small" },
        { id: 4, name: "192.168.2.3" },
        { id: 5, name: "192.168.2.4" },
        { id: 6, name: "192.168.2.5" },
        { id: 7, name: "192.168.2.6" },
        { id: 8, name: "192.168.2.7" },
        { id: 9, name: "192.168.2.254" },
      ],
      links: [
        { source: 1, target: 2, name: "red" },
        { source: 2, target: 8, _color: "f0f" },
        { source: 3, target: 4, class: "red" },
        { source: 4, target: 5 },
        { source: 5, target: 6 },
        { source: 7, target: 8 },
        { source: 5, target: 8 },
        { source: 3, target: 8 },
        { source: 7, target: 9 },
      ],
      options: {
        layout: {
          draggable: false,
          static: false,
        },
        force: {
          collide: 0,
          y: 0.9,
          x: 0.9,
        },
      },
    };
  },
};
</script>
