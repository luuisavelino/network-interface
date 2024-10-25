<template>
  <div>
    <ul>
      <li v-for="message in messages" :key="message.id" class="border-b p-2">
        <strong>{{ message.sender }}</strong>: {{ message.text }} <br />
        <small class="text-gray-500">{{ message.timestamp }}</small>
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
    }
  }
}
</script>
