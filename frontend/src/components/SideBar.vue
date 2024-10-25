<template>
  <div class="w-1/4 bg-gray-200 p-4 border-r">
    <select 
      v-model="selectedDevice" 
      @change="selectDevice" 
      class="mt-2 w-full border border-gray-300 rounded-md p-2 mb-10">
      <option v-for="device in devices" :key="device" :value="device">
        {{ device }}
      </option>
    </select>

    <ul>
      <li @click="select('Read')" class="cursor-pointer hover:bg-gray-300 p-2 rounded">Read</li>
      <li @click="select('Unread')" class="cursor-pointer hover:bg-gray-300 p-2 rounded">Unread</li>
      <li @click="select('Enviados')" class="cursor-pointer hover:bg-gray-300 p-2 rounded">Sent</li>
      <li @click="select('New')" class="cursor-pointer hover:bg-gray-300 mt-6 p-2 rounded">New</li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedDevice: 'Desktop', // Armazena o dispositivo selecionado
      devices: ['Desktop', 'Tablet', 'Mobile'] // Lista de dispositivos disponíveis
    };
  },
  methods: {
    select(tab) {
      this.$emit('select', tab);
    },
    selectDevice() {
      this.$emit('device', this.selectedDevice);
    }
  }
}
</script>
