<template>
<div class="bg-white shadow-md rounded-lg overflow-hidden">
  <table class="min-w-full divide-y divide-gray-200">
    <thead class="bg-gray-100">
      <tr>
        <th class="px-6 py-3 text-left text-gray-900 font-medium">Label</th>
        <th class="px-6 py-3 text-left text-gray-900 font-medium">Battery</th>
        <th class="px-6 py-3 text-left text-gray-900 font-medium">Power</th>
        <th class="px-6 py-3 text-left text-gray-900 font-medium">Status</th>
        <th class="px-6 py-3 text-left text-gray-900 font-medium">Options</th>
      </tr>
    </thead>
    <tbody class="bg-white divide-y divide-gray-200">
      <tr>
        <td class="px-6 py-4 text-gray-600">Descrição do Label</td>
        <td class="px-6 py-4 text-green-500">80%</td>
        <td class="px-6 py-4 text-red-500">Off</td>
        <td class="px-6 py-4 text-blue-500">Active</td>
        <td class="px-6 py-4 flex space-x-2">
          <!-- Ícone Edit -->
          <button class="text-blue-500 hover:text-blue-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232a2.828 2.828 0 014 4L7.5 21H3v-4.5L15.232 5.232z" />
            </svg>
          </button>
          <!-- Ícone Inativar -->
          <button class="text-yellow-500 hover:text-yellow-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 11-12.728 0" />
            </svg>
          </button>
          <!-- Ícone Delete -->
          <button class="text-red-500 hover:text-red-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </td>
      </tr>
      <!-- Adicione mais linhas conforme necessário -->
      <tr>
        <td class="px-6 py-4 text-gray-600">Outro Label</td>
        <td class="px-6 py-4 text-green-500">60%</td>
        <td class="px-6 py-4 text-red-500">On</td>
        <td class="px-6 py-4 text-blue-500">Inactive</td>
        <td class="px-6 py-4 flex space-x-2">
          <!-- Ícone Edit -->
          <button class="text-blue-500 hover:text-blue-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232a2.828 2.828 0 014 4L7.5 21H3v-4.5L15.232 5.232z" />
            </svg>
          </button>
          <!-- Ícone Inativar -->
          <button class="text-yellow-500 hover:text-yellow-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 11-12.728 0" />
            </svg>
          </button>
          <!-- Ícone Delete -->
          <button class="text-red-500 hover:text-red-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</div>

      
</template>

<script>
export default {
  name: 'ListDevice',
  props: {
    device: {
      type: Object,
      required: true
    }
  }
}
</script>