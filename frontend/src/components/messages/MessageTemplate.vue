<template>
  <div class="flex h-1/2 w-full">
    <SideBar @select="selectTab" @selected-device="selectedDevice" :devices="devicesLabel"/>
    <div class="flex-grow p-4">
      <div class="border rounded-lg p-4">
        <div class="flex items-center mb-8 justify-between">
          <h2 class="text-xl font-bold">
            {{ currentTab }}
          </h2>
          <select v-if="currentTab !== 'New'"
            v-model="selectedFilter" 
            class="mt-2 w-2/5 border border-gray-300 rounded-md p-1">
            <option v-for="filter in filters" :key="filter" :value="filter">
              {{ filter }}
            </option>
          </select>
        </div>

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
          <MessageBox :currentDevice="currentDevice" :recipients="this.devicesLabel" @sendMessage="sendMessage" />
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
      devices: {},
      filters: [
        'all',
        'user-message',
        'new-connection',
        'new-connection-ack',
        'confirm-connection',
        'update-routing',
        'update-routing-ack',
      ],
      selectedFilter: 'all',
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

      return this.devices[this.currentDevice]?.messages.received.filter(message => {
        if (this.selectedFilter === 'all') {
          return message.read === status;
        }

        return message.topic === this.selectedFilter && message.read === status;
      });
    },
    sendMessage(data){
      servicesDevices.sendMessage(data)
        .catch(error => {
          console.error(error);
        });
    }
  }
}
</script>
