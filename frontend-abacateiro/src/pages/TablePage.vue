<template>
  <q-page padding>
    <q-btn
      label="Create User"
      color="primary"
      @click="showModal = true"
      icon="person_add"
    />
    <!-- Create User -->
    <q-dialog v-model="showModal" :backdrop-filter="'blur(4px)'">
      <q-card>
        <q-card-section>
          <div class="text-h6">Create User</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit.prevent="submitForm">
            <q-input class="q-pa-sm" v-model="form.name" label="Name" filled />
            <q-input
              class="q-pa-sm"
              v-model="form.email"
              label="Email"
              type="email"
              filled
            />
            <q-input
              class="q-pa-sm"
              v-model="form.document"
              label="RG"
              type="GovernmentID"
              filled
            />
            <q-input
              v-model="form.password"
              class="q-pa-sm"
              label="Password"
              type="password"
              filled
            />
            <q-input
              v-model="form.confirmPassword"
              class="q-pa-sm"
              label="Confirm Password"
              type="password"
              filled
            />
            <q-card-actions>
              <q-btn
                flat
                label="Cancel"
                color="secondary"
                @click="showModal = false"
              />
              <q-btn label="Submit" color="primary" type="submit" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <UserTable />
  </q-page>
</template>

<script setup>
import UserTable from "src/components/table/UserTable.vue";
import { ref, onMounted } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";

const router = useRouter();

defineOptions({
  name: "TablePage",
});

const showModal = ref(false);
const form = ref({
  name: "",
  email: "",
  document: "",
  password: "",
  confirmPassword: "",
});

const submitForm = async () => {
  if (form.value.password !== form.value.confirmPassword) {
    alert("Passwords do not match!");
    return;
  }

  try {
    await axios.post("http://localhost:8080/users", {
      user_name: form.value.name,
      user_email: form.value.email,
      user_password: form.value.password,
      user_document: form.value.document,
    });
    alert("User created successfully!");
    router.go();
  } catch (error) {
    console.error("Error creating user:", error);
    alert("Failed to create user.");
  }

  showModal.value = false;
};
</script>
