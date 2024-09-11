<template>
  <q-dialog v-model="localIsModalOpen" @hide="onDialogHide">
    <q-card
      v-bind:style="$q.screen.lt.sm ? { width: '80%' } : { width: '40%' }"
    >
      <q-card-section>
        <div class="text-center q-pt-lg">
          <div class="col text-h6 ellipsis">
            {{ isEditMode ? "Edit User" : "Create User" }}
          </div>
        </div>
      </q-card-section>
      <q-card-section>
        <q-form class="q-gutter-md">
          <div class="col-6">
            <q-item>
              <q-input
                dense
                outlined
                class="full-width"
                v-model="user.user_name"
                label="Nome"
              />
            </q-item>
          </div>
          <div class="col-6">
            <q-item>
              <q-input
                dense
                outlined
                class="full-width"
                v-model="user.user_email"
                label="Email"
              />
            </q-item>
          </div>
          <div class="col-6">
            <q-item>
              <q-input
                dense
                outlined
                class="full-width"
                v-model="user.user_document"
                label="Documento"
              />
            </q-item>
          </div>
          <div class="col-6" v-if="!isEditMode">
            <q-item>
              <q-input
                dense
                outlined
                class="full-width"
                v-model="user.user_password"
                label="Senha"
                type="password"
              />
            </q-item>
          </div>
          <div class="col-6" v-if="!isEditMode">
            <q-item>
              <q-input
                dense
                outlined
                class="full-width"
                v-model="user.user_password_confirmation"
                label="Confirme a Senha"
                type="password"
              />
            </q-item>
          </div>
          <div class="col-6">
            <q-card-actions align="right">
              <q-btn label="Fechar" color="primary" @click="closeModal" />
              <q-btn label="Salvar" color="primary" @click="saveUser" />
            </q-card-actions>
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script>
import { ref, watch, defineComponent } from "vue";

export default defineComponent({
  name: "UserModal",
  props: {
    isEditMode: {
      type: Boolean,
      default: false,
    },
    isModalOpen: {
      type: Boolean,
      required: true,
    },
    userData: {
      type: Object,
      default: () => ({
        user_name: "",
        user_email: "",
        user_password: "",
        user_password_confirmation: "",
        user_document: "",
      }),
    },
  },
  emits: ["update:isModalOpen", "saveUser"],
  setup(props, { emit }) {
    const user = ref({ ...props.userData });

    const localIsModalOpen = ref(props.isModalOpen);

    watch(
      () => props.isModalOpen,
      (newVal) => {
        localIsModalOpen.value = newVal;
      }
    );

    watch(
      () => props.userData,
      (newVal) => {
        user.value = { ...newVal };
      }
    );

    const closeModal = () => {
      localIsModalOpen.value = false;
      // user.value = {
      //   user_name: '',
      //   user_email: '',
      //   user_password: '',
      //   user_password_confirmation: '',
      //   user_document: ''
      // };
      emit("update:isModalOpen", false);
    };

    const onDialogHide = () => {
      closeModal();
    };

    const saveUser = () => {
      emit("saveUser", user.value);
      closeModal();
    };

    return {
      user,
      localIsModalOpen,
      closeModal,
      saveUser,
      onDialogHide,
    };
  },
});
</script>
