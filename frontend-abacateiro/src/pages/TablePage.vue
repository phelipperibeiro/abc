<template>
  <q-page padding>
    <div class="button-container">
      <q-btn color="primary" label="Create User" @click="openCreateModal" />          
    </div>
    
    <UserTableActions class="q-mt-lg"></UserTableActions>
    
    <UserFormModal 
      :is-modal-open="isModalOpen" 
      :is-edit-mode="isEditMode" 
      :user-data="selectedUser" 
      @update:isModalOpen="isModalOpen = $event"
      @saveUser="handleSaveUser" />
    
  </q-page>
</template>

<script setup>

import UserTable from 'src/components/table/UserTable.vue';
defineOptions({
  name: 'TablePage'
});
</script>
