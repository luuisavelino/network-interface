<template>
  <div>
    <ul>
      <li v-for="message in messages" :key="message.id" class="border-b p-2">
        <strong> {{ targetMessage }}</strong> {{ targetField(message) }} <br />
        {{ message.content }} <br />
        <div class="flex flex-row">
          <small class="text-gray-500">{{ this.formatDate(message.date) }}</small>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    messages: {
      type: Array,
      required: true
    },
    type: {
      type: String,
      required: true
    }
  },
  computed: {
    targetMessage() {
      return this.type === 'sent' ? 'To:' : 'From:';
    },
  },
  methods: {
    formatDate(isoDate) {
      const date = new Date(isoDate);
      return date.toISOString().slice(0, 10).replace(/-/g, '/') + ' ' +
        date.toTimeString().slice(0, 8).replace(/:/g, '-');
    },
    targetField(message) {
      return this.type === 'sent' ? message.destination : message.sender;
    }
  }
}
</script>
