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

<script>
import { ref, defineComponent, defineAsyncComponent } from 'vue';

export default defineComponent({
  name: "UserPage",
  components: {
    UserTableActions: defineAsyncComponent(() => import('components/user/tables/UserTableActions.vue')),
    UserFormModal: defineAsyncComponent(() => import('components/user/form/UserFormModal.vue'))
  },
  setup() {

    const isModalOpen = ref(false);

    const isEditMode = ref(false);

    const selectedUser = ref({
      user_name: '',
      user_email: '',
      user_password: '',
      user_password_confirmation: '',
      user_document: ''
    });

    const openModal = () => {
      isModalOpen.value = true;
    };

    const openCreateModal = () => {

      selectedUser.value = {
        user_name: '',
        user_email: '',
        user_password: '',
        user_password_confirmation: '',
        user_document: ''
      };

      openModal();
    };

    const handleSaveUser = (user) => {
      if (!isEditMode.value) {
        console.log('Create user:', user);      
      } 
    };

    return {
      isModalOpen,
      isEditMode,
      selectedUser,
      openCreateModal,
      handleSaveUser
    };
  }
});
</script>

<style scoped>
.button-container {
  display: flex;
  justify-content: flex-end;
}
</style>
