<template>
  <div class="container">
    <h2 class="my-4">Adicionar dispositivo</h2>
    <form @submit.prevent="handleSubmit">
      <div class="mb-3">
        <label for="id" class="form-label">ID</label>
        <input type="number" v-model="formData.id" class="form-control" id="id" required>
      </div>

      <div class="form-container">
        <div class="form-group">
          <label for="power" class="form-label">Power</label>
          <input type="number" v-model="formData.power" class="form-control" id="power" required>
        </div>

        <div class="form-group">
          <label for="pos_x" class="form-label">Posição X</label>
          <input type="number" v-model="formData.pos_x" class="form-control" id="pos_x" required>
        </div>

        <div class="form-group">
          <label for="pos_y" class="form-label">Posição Y</label>
          <input type="number" v-model="formData.pos_y" class="form-control" id="pos_y" required>
        </div>
      </div>

      <div class="form-container">
        <div class="form-group">
          <label for="walking_speed" class="form-label">Velocidade de Caminhada</label>
          <input type="number" v-model="formData.walking_speed" class="form-control" id="walking_speed" required>
        </div>

        <div class="form-group">
          <label for="message_freq" class="form-label">Frequência de Mensagens</label>
          <input type="number" v-model="formData.message_freq" class="form-control" id="message_freq" required>
        </div>
      </div>

      <button type="submit" class="btn btn-primary">Enviar</button>
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
  name: 'AddDevice',
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
        await servicesDevices.insertDevice(this.formData)
        this.$emit("update-devices");
      } catch (error) {
        console.error('Erro ao enviar os dados:', error);
      }

      this.formData = {
        id: null,
        power: null,
        pos_x: null,
        pos_y: null,
        walking_speed: null,
        message_freq: null
      }
    }
  }
};
</script>

<style scoped>
.form-container {
  display: flex;
  gap: 20px;
}

.form-group {
  flex: 1;
}

.form-label {
  display: block;
  margin-bottom: 5px;
}

.form-control {
  width: 100%;
}
</style>
