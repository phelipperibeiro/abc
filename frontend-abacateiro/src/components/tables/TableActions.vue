<template>
  <q-card>
    <q-card-section>
      <q-table :rows="users" :columns="columns">
        <template v-slot:body-cell-actions="props">
          <q-td align="right" name="actions">
            <q-btn icon="delete" size="sm" class="q-ml-sm" flat dense @click="deleteUser(props.row.id)"/>
          </q-td>
        </template>
      </q-table>
    </q-card-section>
  </q-card>
</template>

<script>

import { defineComponent, ref, onMounted } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: "TableActions",
  setup() {

    const users = ref([]);
    const columns = ref([
      { name: 'user_name', required: true, label: 'Name', align: 'left', field: row => row.user_name, format: val => `${val}`, sortable: true },
      { name: 'user_email', required: true, label: 'Email', align: 'left', field: row => row.user_email, format: val => `${val}`, sortable: true },
      { name: 'user_document', label: 'Document', align: 'left', field: row => row.user_document, format: val => `${val}`, sortable: true },
      { name: 'actions', label: 'Actions', align: 'right' }
    ]);

    const fetchUsers = async () => {
      try {
        const response = await axios.get('http://localhost:8080/users');
        users.value = response.data;
      } catch (error) {
        console.error('Error fetching users:', error);
      }
    };

    onMounted(() => {
      fetchUsers();
    });

    const deleteUser = (id) => {
      console.log(`Delete user with id: ${id}`);
    };

    return {
      users,
      columns,
      deleteUser
    }
  }
})
</script>

<style scoped>

</style>
