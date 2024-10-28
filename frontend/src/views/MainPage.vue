<template>
  <div class="flex m-8">

    <div class="w-1/2">
      <HeaderTemplate @navigate="loadComponent" />
      <div class="container mx-auto mt-5">
        <component :is="currentComponent" @get-route="getRoute" @update-devices="getChart" :devicesLabel="devices" />
      </div>
    </div>

    <div class="w-1/2 flex flex-row-reverse">
      <BubbleChart :chartData="chart" :linesData="linesData" :showLines="showLines" />
    </div>

  </div>
</template>

<script>
import servicesChart from '@/services/api/chart';

import HeaderTemplate from '@/components/HeaderTemplate.vue';
import BubbleChart from '@/components/BubbleChart.vue';
import MessageTemplate from '@/components/messages/MessageTemplate.vue';
import FindBestRoute from '@/components/FindBestRoute.vue';
import AddDevice from '@/components/AddDevice.vue';
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
      chart: {},
      linesData: [],
      showLines: false,
      intervalId: null,
      currentComponent: 'AddDevice'
    };
  },
  computed: {
    devices() {
      return Object.keys(this.chart).map(key => key);
    }
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
    getChart() {
      servicesChart.getChart()
        .then(response => {
          this.chart = response.data;
        })
        .catch(error => {
          console.error('Erro ao buscar os dados:', error);
        });
    },
    loadChart() {
      if (this.intervalId) return;

      this.intervalId = setInterval(() => {
        this.getChart()
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
    this.getChart()
    this.loadChart()
  },
  beforeUnmount() {
    this.stopUpdates();
  }
};
</script>
