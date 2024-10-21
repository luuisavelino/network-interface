<template>
  <div class="container">
    <h2 class="my-4">Adicionar dispositivo</h2>
    <form @submit.prevent="handleSubmit">
      <!-- Campo ID -->
      <div class="form-group">
        <label for="id" class="form-label">ID</label>
        <input type="number" v-model="formData.id" class="form-control" id="id" required>
      </div>

      <!-- Campos Power, Posição X, Posição Y -->
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

      <!-- Campos Velocidade de Caminhada, Frequência de Mensagens -->
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

      <!-- Botão Enviar -->
      <div class="submit-container">
        <button type="submit" class="btn btn-primary">Enviar</button>
      </div>
    </form>

    <!-- Resposta do Servidor -->
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
/* Estilo para o layout dos campos */
.form-container {
  /* margin-bottom: 24px; */
  display: flex;
  gap: 20px;
}

.form-group {
  flex: 1;
  margin-bottom: 24px;
}

.form-label {
  margin: 5px 0;
}

.form-control {
  width: 100%;
}

/* Estilo para o botão de envio no lado direito */
.submit-container {
  display: flex;
  justify-content: flex-end;
}
</style>
