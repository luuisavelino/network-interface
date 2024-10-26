<template>
  <div class="flex m-8">

    <div class="w-1/2">
      <HeaderTemplate @navigate="loadComponent" />
      <div class="container mx-auto mt-5">
        <component :is="currentComponent" @get-route="getRoute" @update-devices="getDevices" :devices="devices"/>
      </div>
    </div>

    <div class="w-1/2 flex flex-row-reverse">
      <BubbleChart :devices="devices" :linesData="linesData" :showLines="showLines" />
    </div>

  </div>
</template>

<script>
import servicesEnvironment from '../services/api/environment';

import HeaderTemplate from '../components/HeaderTemplate.vue';
import BubbleChart from '../components/BubbleChart.vue';
import MessageTemplate from '../components/messages/MessageTemplate.vue';
import FindBestRoute from '../components/FindBestRoute.vue';
import AddDevice from '../components/AddDevice.vue';
import ListDevice from '@/components/ListDevice.vue';

export default {
  name: 'MainPage',
  components: {
    HeaderTemplate,
    BubbleChart,
    AddDevice,
    FindBestRoute,
    MessageTemplate,
    ListDevice,
  },
  data() {
    return {
      devices: [],
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
          this.devices = response.data.devices;
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
