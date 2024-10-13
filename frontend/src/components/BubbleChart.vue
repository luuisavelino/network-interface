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
        maintainAspectRatio: true,
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
        plugins: {
          tooltip: {
            callbacks: {
              label: function(context) {
                const label = context.raw.label || '';
                return `${label} (X: ${context.raw.x}, Y: ${context.raw.y}, R: ${context.raw.r.toFixed(2)})`;
              }
            }
          },
          onResize: function(chart) {
            chart.data.datasets[0].data.forEach((item) => {
              item.r = this.calculateRadiusInChartUnits(item.x, item.y, chart);
            });
            chart.update();
          }
        }
      }
    },
    bubblePosition() {
      if (!this.chartInstance) return;
      const result = {};
      const dataPoints = this.chartInstance.getDatasetMeta(0).data;
      dataPoints.forEach((point, index) => {
        result[point.$context.raw.label] = index;
      });

      return result;
    }
  },
  methods: {
    calculateRadiusInChartUnits(x, y, chart, baseRadius = 10) {
      const xScale = chart.scales.x;
      const yScale = chart.scales.y;

      const xPixelsPerUnit = (xScale.right - xScale.left) / (xScale.max - xScale.min);
      const yPixelsPerUnit = (yScale.bottom - yScale.top) / (yScale.max - yScale.min);

      const pixelsPerUnit = (xPixelsPerUnit + yPixelsPerUnit) / 2;

      return baseRadius / pixelsPerUnit;
    },
    drawLines() {
      this.chartInstance = this.$refs.bubble.chart;
      const ctx = this.chartInstance.ctx;
      this.linesData.forEach(line => {
        const point1 = this.chartInstance.getDatasetMeta(0).data[this.bubblePosition[line.source]];
        const point2 = this.chartInstance.getDatasetMeta(0).data[this.bubblePosition[line.target]];

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
