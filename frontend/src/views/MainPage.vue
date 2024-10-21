<template>
  <div class="main-container">
    <div class="form-container left-column">
      <AddDevice v-on:update-devices="getDevices"/>
      <FindBestRoute v-on:get-route="getRoute"/>
    </div>

    <div class="right-column">
      <BubbleChart 
        :routesData="routesData" :linesData="linesData" :showLines="showLines" 
        class="max-size" />
    </div>

  </div>
</template>

<script>
import BubbleChart from '../components/BubbleChart.vue';
import AddDevice from '../components/AddDevice.vue';
import FindBestRoute from '../components/FindBestRoute.vue';
import servicesEnvironment from '../services/api/environment';

export default {
  name: 'MainPage',
  components: {
    BubbleChart,
    AddDevice,
    FindBestRoute,
  },
  data() {
    return {
      routesData: [],
      linesData: [],
      showLines: false,
      intervalId: null
    };
  },
  methods: {
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

<style scoped>
.main-container {
  display: flex;
  height: 100vh;
}

.left-column {
  flex: 1;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;
  height: 100%;
}

.right-column {
  flex: 2;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
  align-items: center;
}

.max-size {
  width: 50vw;
  height: 100vh;
}

.form-wrapper {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  box-sizing: border-box;
}

</style>
