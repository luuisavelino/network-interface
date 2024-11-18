<template>
  <div>
    <div class="mb-4">
      <label for="recipient" class="block text-sm font-medium">Destinatário:</label>
      <select id="recipient" v-model="selectedRecipient" class="mt-1 block w-full border border-gray-300 rounded-md p-1">
        <option v-for="recipient in filteredRecipient" :key="recipient" :value="recipient">
          {{ recipient }}
        </option>
      </select>
    </div>
    <textarea v-model="message" rows="4" placeholder="Digite sua mensagem..." class="block w-full border border-gray-300 rounded-md p-2"></textarea>
    <button @click="send" class="mt-4 px-4 py-2 bg-blue-500 text-white rounded">Enviar</button>
  </div>
</template>

<script>
export default {
  props: {
    currentDevice: {
      type: String,
      required: true
    },
    recipients: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      selectedRecipient: '',
      message: '',
    };
  },
  computed: {
    filteredRecipient() {
      if (!this.recipients || this.recipients.length === 0) return [];
      return this.recipients.filter(recipient => recipient !== this.currentDevice);
    }
  },
  methods: {
    send() {
      if (this.message && this.selectedRecipient) {
        this.$emit('sendRequest', {
          header: {
            sender: this.currentDevice,
            destination: this.selectedRecipient,
            'content-type': "text"
          },
          body: this.message
        });
        this.message = '';
      } else {
        alert("Por favor, selecione um destinatário e digite uma mensagem.");
      }
    }
  }
}
</script>
