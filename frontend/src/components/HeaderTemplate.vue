<template>
  <header class="bg-gray-900 shadow-lg">
    <nav class="container mx-auto flex justify-between items-center p-4">
      <h1 class="text-white text-2xl font-semibold">Network Emulator</h1>
      <ul class="flex space-x-6">
        <li @click="navigate('AddDevice')" 
            :class="{'text-blue-400 font-bold': selectedTab === 'AddDevice', 'text-white': selectedTab !== 'AddDevice'}"
            class="cursor-pointer transition duration-200">Add</li>
        <li @click="navigate('ListDevice')"
            :class="{'text-blue-400 font-bold': selectedTab === 'ListDevice', 'text-white': selectedTab !== 'ListDevice'}"
            class="cursor-pointer transition duration-200">List</li>
        <li @click="navigate('FindBestRoute')" 
            :class="{'text-blue-400 font-bold': selectedTab === 'FindBestRoute', 'text-white': selectedTab !== 'FindBestRoute'}"
            class="cursor-pointer transition duration-200">Find</li>
        <li @click="navigate('RequestTemplate')" 
            :class="{'text-blue-400 font-bold': selectedTab === 'RequestTemplate', 'text-white': selectedTab !== 'RequestTemplate'}"
            class="cursor-pointer transition duration-200">Requests</li>
      </ul>
    </nav>
  </header>
</template>

<script>
export default {
  data() {
    return {
      selectedTab: ''
    };
  },
  methods: {
    navigate(component) {
      this.selectedTab = component; // Define a aba selecionada
      this.$emit('navigate', component);
    }
  }
}
</script>
