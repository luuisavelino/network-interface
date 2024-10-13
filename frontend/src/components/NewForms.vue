<template>
  <div class="container">
    <h2 class="my-4">Formulário de Dados</h2>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="id">ID</label>
        <input type="number" v-model="formData.id" class="form-control" id="id" required>
      </div>

      <div class="form-group">
        <label for="power">Power</label>
        <input type="number" v-model="formData.power" class="form-control" id="power" required>
      </div>

      <div class="form-group">
        <label for="pos_x">Posição X</label>
        <input type="number" v-model="formData.pos_x" class="form-control" id="pos_x" required>
      </div>

      <div class="form-group">
        <label for="pos_y">Posição Y</label>
        <input type="number" v-model="formData.pos_y" class="form-control" id="pos_y" required>
      </div>

      <div class="form-group">
        <label for="walking_speed">Velocidade de Caminhada</label>
        <input type="number" v-model="formData.walking_speed" class="form-control" id="walking_speed" required>
      </div>

      <div class="form-group">
        <label for="message_freq">Frequência de Mensagens</label>
        <input type="number" v-model="formData.message_freq" class="form-control" id="message_freq" required>
      </div>

      <button type="submit" class="btn btn-primary mt-3">Enviar</button>
    </form>

    <div v-if="responseData" class="alert alert-success mt-4">
      <strong>Resposta do servidor:</strong>
      <pre>{{ responseData }}</pre>
    </div>
  </div>
</template>

<script>
import servicesDevices from '../services/devices';

export default {
  data() {
    return {
      formData: {
        id: null,
        power: null,
        pos_x: null,
        pos_y: null,
        walking_speed: null,
        message_freq: null
      },
      responseData: null
    };
  },
  methods: {
    async handleSubmit() {
      try {
        const response = await servicesDevices.insertDevice(JSON.stringify(this.formData))

        const data = await response.json();
        this.responseData = data;
      } catch (error) {
        console.error('Erro ao enviar os dados:', error);
      }
    }
  }
};
</script>
