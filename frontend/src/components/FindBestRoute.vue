<template>
  <div class="container">
    <h2 class="my-3">Encontrar a melhor rota</h2>
    <form @submit.prevent="handleSubmit">
      <!-- Campos Source e Target lado a lado -->
      <div class="form-container">
        <div class="form-group">
          <label for="source" class="form-label">Source</label>
          <input type="number" v-model="source" class="form-control" required>
        </div>

        <div class="form-group">
          <label for="target" class="form-label">Target</label>
          <input type="number" v-model="target" class="form-control" required>
        </div>
      </div>

      <!-- Botão Search alinhado à direita -->
      <div class="submit-container">
        <button type="submit" class="btn btn-primary">Search</button>
      </div>
    </form>

    <!-- Resposta do servidor -->
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
        const route = await servicesDevices.getRoute(this.source, this.target)
        this.$emit("get-route", route.data);
      } catch (error) {
        console.error('Erro ao enviar os dados:', error);
      }
    }
  }
};
</script>

<style scoped>
/* Estilo para os campos lado a lado */
.form-container {
  margin-bottom: 24px;
  display: flex;
  gap: 20px;
}

.form-group {
  flex: 1;
}

.form-label {
  margin: 5px 0;
}

.form-control {
  width: 100%;
}

/* Botão Search alinhado à direita */
.submit-container {
  display: flex;
  justify-content: flex-end;
}
</style>