<template>
  <div class="container">
    <h3 class="my-3">Encontrar o caminho</h3>
    <form @submit.prevent="handleSubmit">
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

      <div class="submit-container">
        <button type="submit" class="btn btn-primary">Search</button>
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

.submit-container {
  display: flex;
  justify-content: flex-end;
}
</style>