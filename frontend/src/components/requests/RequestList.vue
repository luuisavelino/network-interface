<template>
  <div class="h-screen flex flex-col">
    <div class="flex-grow h-4/5 overflow-y-auto">
      <ul>
        <li v-for="message in messasgesOrdered" :key="message.id" class="border-b p-2">
          <strong>{{ targetRequest }}</strong> {{ targetField(message) }} <br />
          {{ message.content }} <br />
          <div class="flex flex-row">
            <small class="text-gray-500">{{ this.formatDate(message.date) }}</small>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    requests: {
      type: Array,
      required: true
    },
    type: {
      type: String,
      required: true
    }
  },
  computed: {
    targetRequest() {
      return this.type === 'sent' ? 'To:' : 'From:';
    },
    messasgesOrdered() {
      const requests = this.requests;
      return requests?.sort((a, b) => new Date(b.date) - new Date(a.date));
    }
  },
  methods: {
    formatDate(isoDate) {
      const date = new Date(isoDate);
      return date.toISOString().slice(0, 10).replace(/-/g, '/') + ' ' +
        date.toTimeString().slice(0, 8);
    },
    targetField(message) {
      return this.type === 'sent' ? message.header.destination : message.header.sender;
    }
  }
}
</script>
