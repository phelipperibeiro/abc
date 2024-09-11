<template>
    <div class="q-pa-md">
      <q-table
        title="Treats"
        :rows="roles"
        :columns="columns"
        row-key="name"
      />
    </div>
  </template>
  
  <script setup lang="js">
import { ref, onMounted } from 'vue';
import axios from 'axios';

// Call fetchUserRoles when the component is mounted
onMounted(() => fetchUserRoles);

  const columns = [
    { name: 'userName', align: 'center', label: 'Name', field: 'name', sortable: true },
    { name: 'userEmail', label: 'Email', field: 'email', sortable: true },
    { name: 'userGovtID', label: 'RG', field: 'RG', sortable: true },
  ]
  
  const rows = null // get rows from backend
  



// Define reactive state
const roles = ref([]);
const error = ref(null);
// Fetch user roles from the API
const fetchUserRoles = async () => {
  try {
    const response = await axios.get('localhost:8080/users') 
        .then(roles.value = response.data)
  } catch (err) {
    error.value = 'Failed to load user roles';
    console.error(err);
  }
  console.log(roles.value)
};

  </script>