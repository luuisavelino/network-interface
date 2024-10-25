<template>
  <div class="flex m-8 h-screen">

    <div class="w-2/5">
      <HeaderTemplate @navigate="loadComponent" />
      <div class="container mx-auto mt-5">
        <component :is="currentComponent" @get-route="getRoute" @update-devices="getDevices"/>
      </div>
    </div>

    <div class="w-3/5 flex flex-row-reverse">
      <BubbleChart :routesData="routesData" :linesData="linesData" :showLines="showLines" />
    </div>

  </div>
</template>

<script>
import servicesEnvironment from '../services/api/environment';

import HeaderTemplate from '../components/HeaderTemplate.vue';
import BubbleChart from '../components/BubbleChart.vue';
import MessageTemplate from '../components/MessageTemplate.vue';
import FindBestRoute from '../components/FindBestRoute.vue';
import AddDevice from '../components/AddDevice.vue';

export default {
  name: 'MainPage',
  components: {
    HeaderTemplate,
    BubbleChart,
    AddDevice,
    FindBestRoute,
    MessageTemplate,
  },
  data() {
    return {
      routesData: [],
      linesData: [],
      showLines: false,
      intervalId: null,
      currentComponent: 'AddDevice'
    };
  },
  methods: {
    loadComponent(component) {
      this.currentComponent = component;
    },
    getRoute(route) {
      this.linesData = route;
    },
    drawLines() {
      this.showLines = !this.showLines;
    },
    getDevices() {
      servicesEnvironment.getEnvironment()
        .then(response => {
          this.routesData = response.data.devices;
        })
        .catch(error => {
          console.error('Erro ao buscar os dados:', error);
        });
    },
    loadDevices() {
      if (this.intervalId) return;

      this.intervalId = setInterval(() => {
        console.log('Updating routes data');
        this.getDevices()
      }, 5000);
    },
    stopUpdates() {
      if (this.intervalId) {
        clearInterval(this.intervalId);
        this.intervalId = null;
      }
    }
  },
  beforeMount() {
    this.getDevices()
    this.loadDevices()
  },
  beforeUnmount() {
    this.stopUpdates();
  }
};
</script>
