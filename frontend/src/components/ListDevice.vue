<template>
  <div class="bg-white shadow-md rounded-lg overflow-hidden">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-100">
        <tr>
          <th class="px-6 py-3 text-left text-gray-900 font-medium">Label</th>
          <th class="px-6 py-3 text-left text-gray-900 font-medium">Battery</th>
          <th class="px-6 py-3 text-left text-gray-900 font-medium">Power</th>
          <th class="px-6 py-3 text-left text-gray-900 font-medium">Status</th>
          <th class="px-6 py-3 text-left text-gray-900 font-medium">Options</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr v-for="device in devices" :key="device.label">
          <td class="px-6 py-4 text-gray-600">{{ device.label }}</td>
          <td class="px-6 py-4 text-green-500">{{ device.battery }}%</td>
          <td class="px-6 py-4 text-red-500">{{ device.power }}</td>
          <td class="px-6 py-4 text-blue-500">{{ device.status }}</td>
          <td class="px-6 py-4 flex space-x-2">
            <!-- Ícone Mapa -->
            <button @click="openModal(device.label)" class="text-green-500 hover:text-green-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 2l6 3 6-3v16l-6 3-6-3-6 3V5l6-3z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 2v16m6-13v16" />
              </svg>
            </button>

            <!-- Ícone Edit -->
            <!-- <button class="text-blue-500 hover:text-blue-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232a2.828 2.828 0 014 4L7.5 21H3v-4.5L15.232 5.232z" />
              </svg>
            </button> -->
            <!-- Ícone Delete -->
            <button @click="deleteDevice(device.label)" class="text-red-500 hover:text-red-700">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    
    <ModalTemplate v-if="isModalVisible" :deviceLabel="selectedLabel" @close="closeModal" />

  </div>
</template>

<script>

import servicesDevices from '@/services/api/devices';

import ModalTemplate from './ModalTemplate.vue';

export default {
  name: 'ListDevice',
  components: {
    ModalTemplate
  },
  data() {
    return {
      devices: [],
      isModalVisible: false,
    };
  },
  beforeMount() {
    this.getDevices();
  },
  methods: {
    deleteDevice(label) {
      servicesDevices.deleteDevice(label)
        .then(() => {
          this.getDevices();
        })
        .catch(error => {
          console.error(error);
        });
    },
    getDevices() {
      servicesDevices.getDevices()
        .then(response => {
          this.devices = response.data;
        })
        .catch(error => {
          console.error(error);
        });
    },
    openModal(label) {
      this.selectedLabel = label;
      this.isModalVisible = true;
    },
    closeModal() {
      this.isModalVisible = false;
      this.selectedLabel = '';
    }
  }
}
</script>
