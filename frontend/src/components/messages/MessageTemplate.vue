<template>
  <div class="flex h-full w-full">
    <SideBar @select="selectTab" @selected-device="selectedDevice" :devices="devicesLabel"/>
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
          <MessageBox :recipients="devices" @sendMessage="sendMessage" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import SideBar from './SideBar.vue';
import MessageList from './MessageList.vue';
import MessageBox from './MessageBox.vue';

import servicesDevices from '@/services/api/devices';

export default {
  components: {
    SideBar,
    MessageList,
    MessageBox,
  },
  props: {
    devicesLabel: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      currentTab: 'New',
      currentDevice: this.devicesLabel[0],
      devices: {}
    };
  },
  methods: {
    selectTab(tab) {
      this.currentTab = tab;
    },
    selectedDevice(device) {
      this.currentDevice = device;
    },
    filteredMessages(type, status) {
      if (!this.devices[this.currentDevice]) {
        servicesDevices.getDeviceById(this.currentDevice)
        .then(response => this.devices[this.currentDevice] = response.data);
      }

      if (type === 'sent') {
        return this.devices[this.currentDevice]?.messages.sent;
      }

      return this.devices[this.currentDevice]?.messages.received.filter(message => message.read === status);
    }
  }
}
</script>
