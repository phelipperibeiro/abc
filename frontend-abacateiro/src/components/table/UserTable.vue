<template>
  <div class="q-pa-md">
    <!-- <q-responsive> -->
    <q-table title="Usuários" :rows="rows" :columns="columns" row-key="name">
      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn icon="mode_edit" @click="onEdit(props.row)"></q-btn>
          <q-btn icon="delete" @click="onDelete(props.row)"></q-btn>
        </q-td>
      </template>
    </q-table>
    <q-dialog v-model="editModal">
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Edit User</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="editUser.user_name" label="Name" />
          <q-input v-model="editUser.user_email" label="Email" />
          <q-input v-model="editUser.user_document" label="RG" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" color="primary" v-close-popup />
          <q-btn flat label="Save" color="primary" @click="saveEdit" />
        </q-card-actions>
      </q-card>
    </q-dialog>
    <!-- </q-responsive>  -->
  </div>
</template>

<script setup lang="js">
import { ref, onMounted } from 'vue';
import axios from 'axios';

// Call fetchUserRoles when the component is mounted
onMounted(() => fetchUserRoles());

  const columns = [
    { name: 'userName', align: 'left', label: 'Name', field: (row) => row.user_name, sortable: true },
    { name: 'userEmail', align: 'left',label: 'Email', field: (row) => row.user_email, sortable: true },
    { name: 'userGovtID', align: 'left', label: 'RG', field: (row) => row.user_document, sortable: true },
    { name: 'actions', label: 'Action' }
  ]

// Define reactive state
const rows = ref([]);
const error = ref(null);
// Fetch user roles from the API
const fetchUserRoles = async () => {
  try {
    const response = await axios.get('http://localhost:8080/users')
    rows.value = response.data
  } catch (err) {
    error.value = 'Failed to load user roles';
    console.error(err);
  }
};

defineExpose({fetchUserRoles})

var editModal = ref(false)
const onEdit = (row) => {
  editUser.value = {...row};
  editModal.value = true;
};

const editUser = ref({
  name: "",
  email: "",
  document: "",
});

const saveEdit = async () => {
  try {
    await axios.put(`http://localhost:8080/users/${editUser.value.id}`, {
      user_name: editUser.value.user_name,
      user_email: editUser.value.user_email,
      user_document: editUser.value.user_document
    });
    await fetchUserRoles();
    editModal.value = false;
  } catch (err) {
    console.error('Failed to update user:', err);
  }
};

const onDelete = async (row) => {
  try {
    await axios.delete(`http://localhost:8080/users/${row.id}`);
    await fetchUserRoles();
  } catch (err) {
    console.error('Failed to delete user:', err);
  }
};
</script>
