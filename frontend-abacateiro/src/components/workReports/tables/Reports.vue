<template>
  <q-card class="no-shadow" bordered>

    <q-card-section>
      <div class="text-h6 text-grey-8">
        Passagens de serviço
      </div>
    </q-card-section>

    <q-separator></q-separator>

    <q-card-section class="q-pa-none">
      <q-table square class="no-shadow"
        title="Treats"
        :rows="workReports"
        :columns="columns"
        row-key="work_report_id"
        :filter="filter"
      >
        <template v-slot:top-right>
          <q-input v-if="show_filter" filled borderless dense debounce="300" v-model="filter" placeholder="Search">
            <template v-slot:append>
              <q-icon name="search"/>
            </template>
          </q-input>

          <q-btn class="q-ml-sm" icon="filter_list" @click="show_filter=!show_filter" flat/>
        </template>

        <template v-slot:body-cell-Action="props">
          <q-td :props="props">
            <!-- <q-btn icon="edit" size="sm" flat dense/> -->
            <!-- <q-btn icon="delete" size="sm" class="q-ml-sm" flat dense/> -->
            <q-btn icon="download" size="sm" flat round @click="downloadReport(props.row)" />
          </q-td>
        </template>

      </q-table>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref, onMounted } from "vue";
import axios from "axios";

export default defineComponent({
  name: "ReportsComponent",
  setup() {
    const workReports = ref([]);

    const columns = ref([
      {
       name: "unit_id",
       label: 'Unidade',
       field: row => row.unit.unit_name,
       format: val => `${val}`,
       sortable: true,
       align: 'left',
      },
      {
        name: "work_report_from",
        label: 'Embarque',
        field: row => row.work_report_from,
        format: val => {
          const date = new Date(val);
          return `${date.toLocaleDateString('pt-BR')} ${date.toLocaleTimeString('pt-BR')}`;
        },
        sortable: true,
        align: 'left',
      },
      {
        name: "work_report_to",
        label: 'Desembarque',
        field: row => row.unit_id,
        format: val => {
          const date = new Date(val);
          return `${date.toLocaleDateString('pt-BR')} ${date.toLocaleTimeString('pt-BR')}`;
        },
        sortable: true,
        align: 'left',
      },
      {
        name: "work_report_docname",
        label: 'Titulo',
        field: row => row.work_report_docname,
        format: val => `${val}`,
        sortable: true,
        align: 'left',
      },
      {name: 'Action', label: '', field: 'Action', sortable: false, align: 'center'}
    ]);

    const downloadReport = (row) => {
      console.log('Download report:', row);
      // Adicione aqui a lógica para download do relatório
    }

    const fetchWorkReports = async () => {
      try {
        const response = await axios.get("http://localhost:8888/work-reports");
        workReports.value = Array.isArray(response.data.work_reports) ? response.data.work_reports : [];
      } catch (error) {
        console.error("Erro ao buscar relatórios de trabalho:", error);
      }
    };

    onMounted(() => {
      fetchWorkReports();
    });

    const show_filter = ref(false);
    return {
      filter: ref(''),
      show_filter,
      workReports,
      columns,
      downloadReport
    }
  }
})
</script>

<style scoped>
</style>
