<template>
  <div class="container mx-auto p-4">
    <h3 class="text-2xl font-bold mb-4">Adicionar dispositivo</h3>
    <form @submit.prevent="handleSubmit">
      <div class="mb-4">
        <label for="label" class="block text-sm font-medium text-gray-700">Label</label>
        <input type="text" v-model="formData.label" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="label" required>
      </div>

      <div class="flex gap-4 mb-4">
        <div class="flex-1">
          <label for="power" class="block text-sm font-medium text-gray-700">Power</label>
          <input type="number" v-model="formData.power" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="power" required>
        </div>

        <div class="flex-1">
          <label for="battery" class="block text-sm font-medium text-gray-700">Battery</label>
          <input type="number" v-model="formData.battery" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="battery" required>
        </div>
      </div>

      <div class="flex gap-4 mb-4">
        <div class="flex-1">
          <label for="walking_speed" class="block text-sm font-medium text-gray-700">Velocidade de Caminhada</label>
          <input type="number" v-model="formData.walking_speed" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="walking_speed" required>
        </div>

        <div class="flex-1">
          <label for="message_freq" class="block text-sm font-medium text-gray-700">Frequência de Mensagens</label>
          <input type="number" v-model="formData.message_freq" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="message_freq" required>
        </div>
      </div>

      <div class="flex justify-end">
        <button type="submit" class="btn btn-primary bg-indigo-600 text-white py-2 px-4 rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">Enviar</button>
      </div>
    </form>

    <div v-if="responseData" class="alert alert-success mt-4 bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative">
      <strong class="font-bold">Resposta do servidor:</strong>
      <pre>{{ responseData }}</pre>
    </div>
  </div>
</template>

<script>
import servicesDevices from '../services/api/devices';

export default {
  name: 'AddDevice',
  data() {
    return {
      formData: {
        label: null,
        power: null,
        battery: null,
        walking_speed: null,
        message_freq: null
      },
      responseData: null
    };
  },
  methods: {
    async handleSubmit() {
      try {
        await servicesDevices.insertDevice(this.formData)
        this.$emit("update-devices");
      } catch (error) {
        console.error('Erro ao enviar os dados:', error);
      }

      this.formData = {
        label: null,
        power: null,
        battery: null,
        walking_speed: null,
        message_freq: null
      }
    }
  }
};
</script>
