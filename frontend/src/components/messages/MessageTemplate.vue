<template>
  <div class="flex h-full w-full">
    <SideBar @select="selectTab" @selected-device="selectedDevice" :devices="recipients"/>
    <div class="flex-grow p-4">
      <div class="border rounded-lg p-4">
        <h2 class="text-xl font-bold mb-4">{{ currentTab }}</h2>
        <div v-if="currentTab === 'Read'">
          <MessageList :type="'received'" :messages="filteredMessages('received', true)" />
        </div>
        <div v-if="currentTab === 'Unread'">
          <MessageList :type="'received'" :messages="filteredMessages('received', false)" />
        </div>
        <div v-if="currentTab === 'Sent'">
          <MessageList :type="'sent'" :messages="filteredMessages('sent', false)" />
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
  props: {
    devices: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      currentTab: 'New',
      currentDevice: this.devices[0]?.label,
    };
  },
  computed: {
    recipients() {
      return this.devices?.map(device => device?.label);
    },
  },
  methods: {
    selectTab(tab) {
      this.currentTab = tab;
    },
    selectedDevice(device) {
      this.currentDevice = device;
    },
    filteredMessages(type, status) {
      const device = this.devices.filter(device => device?.label === this.currentDevice)[0];

      if (type === 'sent') {
        return device?.messages.sent;
      }

      return device?.messages.received.filter(message => message.read === status);
    }
  }
}
</script>
