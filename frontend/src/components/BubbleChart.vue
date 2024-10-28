<template>
  <Bubble ref="bubble" :data="getChartData" :options="getChartOptions" style="height:600px; width:600px" />
</template>

<script>
import {
  Chart as ChartJS,
  Tooltip,
  Legend,
  PointElement,
  LinearScale
} from 'chart.js'
import { Bubble } from 'vue-chartjs'

export default {
  name: 'BubbleChart',
  components: {
    Bubble
  },
  props: {
    chartData: {
      type: Object,
      required: true
    },
    linesData: {
      type: Array,
      required: true
    },
    showLines: {
      type: Boolean,
      required: true
    }
  },
  data() {
    return {
      chart: null,
      bubbleColor: {}
    };
  },
  computed: {
    customData() {
      return Object.keys(this.chartData).map(key => {
        const route = this.chartData[key];
        return {
          backgroundColor: 'rgba(1, 0, 132, 0.6)',
          x: route.x,
          y: route.y,
          r: this.calculateRadiusInChartUnits(route.r),
          label: key
        };
      });
    },
    getChartData() {
      return {
        datasets: this.customData.map(data => {
          return {
            label: data.label,
            data: [data],
            backgroundColor: this.getBubbleColor(data.label),
            borderColor: this.getBubbleColor(data.label),
            borderWidth: 1,
          }
        })
      }
    },
    getChartOptions() {
      return {
        responsive: false,
        maintainAspectRatio: true,
        scales: {
          x: {
            type: 'linear',
            position: 'bottom',
            beginAtZero: true,
            max: 50,
            tickLength: 100,
          },
          y: {
            type: 'linear',
            position: 'left',
            beginAtZero: true,
            max: 50,
            tickLength: 100,
          }
        },
        layout: {
          autoPadding: false,
        },
        plugins: {
          tooltip: {
            callbacks: {
              label: function(context) {
                const label = context.raw.label || '';
                return `${label} (X: ${context.raw.x}, Y: ${context.raw.y}, R: ${context.raw.r.toFixed(2)})`;
              },
            }
          },
        }
      }
    },
    bubblePosition() {
      if (!this.chart) return;
      const result = {};

      for (let i = 0; i < this.customData.length; i++) {
        result[this.customData[i].label] = i;
      }

      return result;
    },
  },
  beforeMount() {
    ChartJS.register(LinearScale, PointElement, Tooltip, Legend)
  },
  mounted() {
    this.$nextTick(() => {
      this.chart = this.$refs.bubble.chart;
    });
  },
  methods: {
    getBubbleColor(bubbleLabel) {
      if (!this.bubbleColor[bubbleLabel]) {
        this.bubbleColor[bubbleLabel] = this.randomColor();
      }
      return this.bubbleColor[bubbleLabel];
    },
    randomColor() {
      const r = Math.floor(Math.random() * 256);
      const g = Math.floor(Math.random() * 256);
      const b = Math.floor(Math.random() * 256);
      return `rgb(${r}, ${g}, ${b}, 0.4)`;
    },
    calculateRadiusInChartUnits(baseRadius = 10) {
      const xScale = this.chart.scales.x;
      const xPixelsPerUnit = (xScale.width) / 50;
      return (baseRadius * xPixelsPerUnit);
    },
    drawLines() {
      const ctx = this.chart.ctx;
      this.linesData.forEach(line => {
        const point1 = this.chart.getDatasetMeta(this.bubblePosition[line.source]).data[0];
        const point2 = this.chart.getDatasetMeta(this.bubblePosition[line.target]).data[0];

        ctx.save();
        ctx.beginPath();
        ctx.moveTo(point1.x, point1.y);
        ctx.lineTo(point2.x, point2.y);
        ctx.lineWidth = 2;
        ctx.strokeStyle = 'red';
        ctx.stroke();
        ctx.restore();
      });
    },
  },
  watch: {
    linesData: {
      handler() {
        this.drawLines();
      },
      deep: true
    },
  }
};
</script>
