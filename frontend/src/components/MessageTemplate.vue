<template>
  <div class="flex h-full w-full">
    <SideBar @select="selectTab" />
    <div class="flex-grow p-4">
      <div class="border rounded-lg p-4">
        <h2 class="text-xl font-bold mb-4">{{ currentTab }}</h2>
        <div v-if="currentTab === 'Read'">
          <MessageList :messages="filteredMessages('read')" />
        </div>
        <div v-if="currentTab === 'Unread'">
          <MessageList :messages="filteredMessages('unread')" />
        </div>
        <div v-if="currentTab === 'Sent'">
          <MessageList :messages="sentMessages" />
        </div>
        <div v-if="currentTab === 'New'">
          <MessageBox :recipients="recipients" @sendMessage="sendMessage" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import SideBar from './SideBar.vue';
import MessageList from './MessageList.vue';
import MessageBox from './MessageBox.vue';

export default {
  components: {
    SideBar,
    MessageList,
    MessageBox,
  },
  data() {
    return {
      currentTab: 'New',
      messages: [
        { id: 1, sender: 'João', text: 'Olá!', timestamp: '2024-10-01 10:00', status: 'read' },
        { id: 2, sender: 'Maria', text: 'Como você está?', timestamp: '2024-10-02 11:00', status: 'unread' },
      ],
      sentMessages: [
        { id: 3, sender: 'Você', text: 'Oi, tudo bem?', timestamp: '2024-10-01 12:00' },
      ],
      recipients: ['João', 'Maria', 'Carlos'],
    };
  },
  methods: {
    selectTab(tab) {
      this.currentTab = tab;
    },
    sendMessage(message) {
      this.sentMessages.push({
        id: this.sentMessages.length + 1,
        sender: 'Você',
        text: message.message,
        timestamp: new Date().toLocaleString()
      });
      this.messages.push({
        id: this.messages.length + 1, // Gerar um novo ID
        sender: 'Você',
        text: message.message,
        timestamp: new Date().toLocaleString(),
        status: 'unread' // Status ao enviar
      });
      console.log('Mensagem enviada:', message);
    },
    filteredMessages(status) {
      return this.messages.filter(message => message.status === status);
    }
  }
}
</script>
