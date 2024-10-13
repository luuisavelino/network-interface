<template>
  <div>

    <b-form-group label="Email">
      <b-form-input
      v-model="email"
      type="email"
      placeholder="Digite seu e-mail"
      ></b-form-input>
    </b-form-group>

    <button @click="drawLines">Desenhar Linhas</button>
    <button @click="startRandomUpdates">Iniciar Atualizações Aleatórias</button>

    <BubbleChart :routesData="routesData" :linesData="linesData" :showLines="showLines" />
  </div>
</template>

<script>
import BubbleChart from '../components/BubbleChart.vue';

export default {
  name: 'MainPage',
  components: {
    BubbleChart,
  },
  data() {
    return {
      routesData: [
        { x: 10, y: 20, r: 10 },
        { x: 10, y: 30, r: 10 },
        { x: 10, y: 40, r: 10 },
        { x: 10, y: 50, r: 10 },
        { x: 10, y: 60, r: 10 }
      ],
      linesData: [
        { index1: 0, index2: 1 },
        { index1: 1, index2: 2 }
      ],
      showLines: false,
      intervalId: null
    };
  },
  methods: {
    drawLines() {
      this.showLines = !this.showLines;
    },
    generateRandomData() {
      return [
        { x: Math.floor(Math.random() * 10), y: Math.floor(Math.random() * 10), r: 15 },
        { x: Math.floor(Math.random() * 10), y: Math.floor(Math.random() * 10), r: Math.random() * 5 },
        { x: Math.floor(Math.random() * 10), y: Math.floor(Math.random() * 10), r: Math.random() * 5 }
      ];
    },
    startRandomUpdates() {
      if (this.intervalId) return;

      this.intervalId = setInterval(() => {
        console.log('Updating routes data');
        this.routesData = this.generateRandomData();
      }, 5000);
    },
    stopRandomUpdates() {
      if (this.intervalId) {
        clearInterval(this.intervalId);
        this.intervalId = null;
      }
    }
  },
  beforeUnmount() {
    this.stopRandomUpdates();
  }
};
</script>
