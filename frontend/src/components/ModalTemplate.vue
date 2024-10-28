<template>
  <div class="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50" @click.self="closeModal">
    <div class="bg-white w-1/3 rounded-lg shadow-lg">
      <!-- Header -->
      <div class="px-4 py-2 bg-gray-800 text-white rounded-t-lg">
        <h2 class="text-lg font-semibold">Device {{ deviceLabel }}</h2>
        <button @click="closeModal" class="text-white float-right">&times;</button>
      </div>

      <!-- Body -->
      <div class="p-4">
        <label class="block text-gray-700">X:</label>
        <input v-model="x" type="number" class="w-full px-3 py-2 border rounded-md" />

        <label class="block text-gray-700 mt-4">Y:</label>
        <input v-model="y" type="number" class="w-full px-3 py-2 border rounded-md" />
      </div>

      <!-- Footer -->
      <div class="px-4 py-2 bg-gray-100 text-right rounded-b-lg">
        <button @click="send" class="bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-700">Send</button>
      </div>
    </div>
  </div>
</template>

<script>

import servicesChart from '@/services/api/chart';

export default {
  data() {
    return {
      x: null,
      y: null
    };
  },
  props: {
    deviceLabel: {
      type: String,
      required: true
    },
  },
  methods: {
    closeModal() {
      this.$emit("close");
    },
    send() {
      servicesChart.setDeviceInChart(this.deviceLabel, { x: this.x, y: this.y })
        .catch(err => console.error(err));
      this.closeModal();
    }
  }
};
</script>
