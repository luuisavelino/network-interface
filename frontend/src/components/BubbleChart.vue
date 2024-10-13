<template>
  <div>
    <Bubble :data="getChartData" :options="getChartOptions" ref="bubble"/>
  </div>
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

ChartJS.register(LinearScale, PointElement, Tooltip, Legend)

export default {
  name: 'BubbleChart',
  components: {
    Bubble
  },
  props: {
    routesData: {
      type: Array,
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
      chartInstance: null,
    };
  },
  computed: {
    getChartData() {
      return {
        datasets: [
          {
            label: 'Routes',
            data: this.routesData,
            backgroundColor: 'rgba(255, 99, 132, 0.6)',
            borderColor: 'rgba(255, 99, 132, 1)',
            borderWidth: 1
          }
        ]
      }
    },
    getChartOptions() {
      return {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          x: {
            type: 'linear',
            position: 'bottom',
            beginAtZero: true,
          },
          y: {
            type: 'linear',
            position: 'left',
            beginAtZero: true,
          }
        },
      }
    },
  },
  methods: {
    drawLines() {
      this.chartInstance = this.$refs.bubble.chart;

      const ctx = this.chartInstance.ctx;

      this.linesData.forEach(line => {
        const point1 = this.chartInstance.getDatasetMeta(0).data[line.index1];
        const point2 = this.chartInstance.getDatasetMeta(0).data[line.index2];

        ctx.save();
        ctx.beginPath();
        ctx.moveTo(point1.x, point1.y);
        ctx.lineTo(point2.x, point2.y);
        ctx.lineWidth = 2;
        ctx.strokeStyle = 'red';
        ctx.stroke();
        ctx.restore();
      });
    }
  },
  watch: {
    showLines: {
      handler() {
        if (this.showLines) {
          this.drawLines();
        }
      },
      deep: true
    },
  }
};
</script>

