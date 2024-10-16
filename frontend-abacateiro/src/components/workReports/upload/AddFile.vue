<!-- Solução adaptada de https://www.digitalocean.com/community/tutorials/how-to-handle-file-uploads-in-vue-2 -->
<template>

  <div class="container q-ma-md">
    <!--UPLOAD-->

    <form enctype="multipart/form-data" novalidate v-if="isInitial || isSaving">

      <div class="dropbox">

        <input
          type="file"
          :name="uploadFieldName"
          :disabled="isSaving"
          accept="*"
          class="input-file"
          @change="filesChange($event.target.name, $event.target.files); fileCount = $event.target.files.length;" />

        <p v-if="isInitial">
          Clique para selecionar ou arraste o arquivo.
          <br />
          <br />
          Limite máximo de 200MiB.
        </p>

        <p v-if="isSaving">Enviando e processando arquivo...</p>

      </div>

    </form>

    <!--SUCCESS-->
    <div v-if="isSuccess">
        <h6>Arquivo armazenado com sucesso.</h6>
    </div>



    <!--FAILED-->
    <div style="overflow-y: hidden !important; overflow-x: hidden !important;" v-if="isFailed">
      <h6>Falha no envio.</h6>
      <p class="text-negative" v-if="uploadError.data.error">
        {{ uploadError.data.error }}
      </p>
      <q-btn color="primary" label="Tente novamente" @click="reset" />
      <pre>{{ uploadError }}</pre>
    </div>

  </div>

</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { Notify } from 'quasar';
//import { uploadData } from '@/upload/file-upload.fake';
import { uploadData } from '@/plugins/file-upload';

const STATUS_INITIAL = 0;
const STATUS_SAVING = 1;
const STATUS_SUCCESS = 2;
const STATUS_FAILED = 3;

export default {
  name: 'AddFile',
  emits: [
    "close",
    "closeUpdate"
  ],
  props: [
    "api"
  ],
  setup(props, { emit }) {

    const uploadedFiles = ref([]);
    const uploadError = ref(null);
    const currentStatus = ref(STATUS_INITIAL);
    const uploadFieldName = ref("file");

    const isInitial = computed(() => currentStatus.value === STATUS_INITIAL);
    const isSaving = computed(() => currentStatus.value === STATUS_SAVING);
    const isSuccess = computed(() => currentStatus.value === STATUS_SUCCESS);
    const isFailed = computed(() => currentStatus.value === STATUS_FAILED);

    const reset = () => {
      currentStatus.value = STATUS_INITIAL;
      uploadedFiles.value = [];
      uploadError.value = null;
    };

    const save = (formData, name) => {

      currentStatus.value = STATUS_SAVING;

      uploadData(formData, props.api + "/" + name)
        .then((x) => {
          uploadedFiles.value = [].concat(x);
          currentStatus.value = STATUS_SUCCESS;
          emit("closeUpdate");
        })
        .catch((err) => {
          uploadError.value = err.response;
          currentStatus.value = STATUS_FAILED;
        });

    };

    const filesChange = (fieldName, fileList) => {

      console.log('filesChange');
      console.log('fieldName', fieldName);
      console.log('fileList', fileList[0].name);

      const formData = new FormData();

      if (!fileList.length) return;

      if (fileList[0].size > 200 * 1024 * 1024) {
        Notify.create({
          type: "negative",
          message: "Erro: arquivo ultrapassa o limite de 200MiB.",
        });
        return;
      }

      Array.from(Array(fileList.length).keys()).map((x) => {
        formData.append(fieldName, fileList[x], fileList[x].name);
      });

      save(formData, fileList[0].name);

    };

    onMounted(() => {
      reset();
    });

    return {
      uploadedFiles,
      uploadError,
      currentStatus,
      uploadFieldName,
      isInitial,
      isSaving,
      isSuccess,
      isFailed,
      reset,
      save,
      filesChange,
    };
  },
};
</script>

<style lang="scss">
  .dropbox {
    outline: 2px dashed grey; /* the dash box */
    outline-offset: -10px;
    background: lightcyan;
    color: dimgray;
    padding: 10px 10px;
    min-height: 200px; /* minimum height */
    position: relative;
    cursor: pointer;
    overflow: hidden;
    overflow-y: hidden;
    overflow-x: hidden;
  }

  .input-file {
    opacity: 0; /* invisible but it's there! */
    width: 100%;
    height: 200px;
    position: absolute;
    cursor: pointer;
  }

  .dropbox:hover {
    background: lightblue; /* when mouse over to the drop zone, change color */
  }

  .dropbox p {
    font-size: 1.2em;
    text-align: center;
    padding: 50px 0;
    overflow: hidden;
    overflow-y: hidden;
    overflow-x: hidden;
  }
</style>
