<template>
  <div class="main-container">
    <div class="form-container left-column">
      <div class="form-wrapper">
        <AddDevice v-on:update-devices="getDevices"/>
      </div>
      <div class="form-wrapper">
        <FindBestRoute v-on:get-route="getRoute"/>
      </div>
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
  },
  beforeUnmount() {
    this.stopUpdates();
  }
};
</script>

<style scoped>
.main-container {
  display: flex;
  height: 100vh; /* Altura total da tela */
}

.left-column {
  flex: 1;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;
}

.right-column {
  flex: 2;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: flex 0.3s ease;
  width: 100%;
  height: 100%;
}

.max-size {
  width: 100%;
  height: 100%;
}

.form-wrapper {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  box-sizing: border-box;
}

</style>
