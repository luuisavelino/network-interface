<template>
  <div class="flex h-1/2 w-full">
    <SideBar @select="selectTab" @selected-device="selectedDevice" :devices="devicesLabel ?? []"/>
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
          <RequestList :type="'received'" :requests="filteredRequests('received', true)" />
        </div>
        <div v-if="currentTab === 'Unread'">
          <RequestList :type="'received'" :requests="filteredRequests('received', false)" />
        </div>
        <div v-if="currentTab === 'Sent'">
          <RequestList :type="'sent'" :requests="filteredRequests('sent', false)" />
        </div>
        <div v-if="currentTab === 'New'">
          <RequestBox :currentDevice="currentDevice" :recipients="this.devicesLabel" @sendRequest="sendRequest" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import SideBar from './SideBar.vue';
import RequestList from './RequestList.vue';
import RequestBox from './RequestBox.vue';

import servicesDevices from '@/services/api/devices';

export default {
  components: {
    SideBar,
    RequestList,
    RequestBox,
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
    filteredRequests(type, status) {
      if (!this.devices[this.currentDevice]) {
        servicesDevices.getDeviceById(this.currentDevice)
        .then(response => this.devices[this.currentDevice] = response.data);
      }

      if (type === 'sent') {
        return this.devices[this.currentDevice]?.requests?.sent;
      }
      
      return this.devices[this.currentDevice]?.requests?.received.filter(message => {
        if (this.selectedFilter === 'all') {
          return message.read === status;
        }

        return message.header.topic === this.selectedFilter && message.read === status;
      });
    },
    sendRequest(data){
      servicesDevices.sendRequest(data)
        .catch(error => {
          console.error(error);
        });
    }
  }
}
</script>
