<template>
  <q-card>
    <q-card-section>
      <q-table :rows="users" :columns="columns">
        <template v-slot:body-cell-actions="props">
          <q-td align="right" name="actions">
            <q-btn
              icon="edit"
              size="sm"
              class="q-ml-sm"
              flat
              dense
              @click="editUser(props.row)"
            />
            <q-btn
              icon="delete"
              size="sm"
              class="q-ml-sm"
              flat
              dense
              @click="deleteUser(props.row.id)"
            />
          </q-td>
        </template>
      </q-table>
    </q-card-section>

    <UserFormModal
      :is-modal-open="isModalOpen"
      :is-edit-mode="isEditMode"
      :user-data="selectedUser"
      @update:isModalOpen="isModalOpen = $event"
      @saveUser="handleSaveUser"
    />
  </q-card>
</template>

<script>
import { defineComponent, ref, onMounted, defineAsyncComponent } from "vue";
import { EventBus } from '@/plugins/eventBus';
import axios from "axios";

export default defineComponent({
  name: "TableActions",
  components: {
    UserFormModal: defineAsyncComponent(() =>
      import("components/user/form/UserFormModal.vue")
    ),
  },
  setup() {
    const users = ref([]);
    const columns = ref([
      {
        name: "user_name",
        required: true,
        label: "Name",
        align: "left",
        field: (row) => row.user_name,
        format: (val) => `${val}`,
        sortable: true,
      },
      {
        name: "user_email",
        required: true,
        label: "Email",
        align: "left",
        field: (row) => row.user_email,
        format: (val) => `${val}`,
        sortable: true,
      },
      {
        name: "user_document",
        label: "Document",
        align: "left",
        field: (row) => row.user_document,
        format: (val) => `${val}`,
        sortable: true,
      },
      { name: "actions", label: "Actions", align: "right" },
    ]);
    const isModalOpen = ref(false);
    const isEditMode = ref(true);

    const selectedUser = ref({
      user_name: "",
      user_email: "",
      user_password: "",
      user_password_confirmation: "",
      user_document: "",
    });

    const fetchUsers = async () => {
      try {
        const response = await axios.get("http://localhost:8080/users");
        users.value = response.data;
      } catch (error) {
        console.error("Error fetching users:", error);
      }
    };

    const handleUserSaved = () => {
      fetchUsers();
    };

    onMounted(() => {
      fetchUsers();
      EventBus.on('user-saved', handleUserSaved);
    });

    const handleSaveUser = (user) => {
      if (isEditMode.value) {
        console.log("Edit user:", user);
      }
    };

    const deleteUser = (id) => {
      console.log(`Delete user with id: ${id}`);
    };

    const openModal = () => {
      isModalOpen.value = true;
    };

    const editUser = (user) => {
      console.log("Edit user:", user);

      // essa linha só funciona com hot reload
      selectedUser.value = {
        user_name: user.user_name,
        user_email: user.user_email,
        user_document: user.user_document,
      };
      openModal();
    };

    return {
      users,
      isModalOpen,
      isEditMode,
      selectedUser,
      columns,
      handleSaveUser,
      deleteUser,
      editUser,
    };
  },
});
</script>

<style scoped></style>
