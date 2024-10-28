<template>
  <div class="container mx-auto p-4">
    <h3 class="text-2xl font-bold mb-4">Encontrar o caminho</h3>
    <form @submit.prevent="handleSubmit">

      <div class="flex gap-4 mb-4">
        <div class="flex-1">
          <label for="walking_speed" class="block text-sm font-medium text-gray-700">Source</label>
          <input type="text" v-model="source" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="walking_speed" required>
        </div>

        <div class="flex-1">
          <label for="message_freq" class="block text-sm font-medium text-gray-700">Target</label>
          <input type="text" v-model="target" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-md p-1" id="message_freq" required>
        </div>
      </div>

      <div class="flex justify-end">
        <button type="submit" class="btn btn-primary bg-indigo-600 text-white py-2 px-4 rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">Enviar</button>
      </div>
    </form>

    <div v-if="responseData" class="alert alert-success mt-4">
      <strong>Resposta do servidor:</strong>
      <pre>{{ responseData }}</pre>
    </div>
  </div>
</template>

<script>
import servicesDevices from '../services/api/devices';

export default {
  name: 'FindBestRoute',
  data() {
    return {
      source: null,
      target: null,
      responseData: null
    };
  },
  methods: {
    async handleSubmit() {
      try {
        const route = await servicesDevices.getRoute(this.source, this.target);
        this.$emit("get-route", route.data);

        this.source = null
        this.target = null
      } catch (error) {
        console.error('Erro ao enviar os dados:', error);
      }
    }
  }
};
</script>
