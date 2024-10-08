<template>
  <q-card class="no-shadow" bordered>
    <q-card-section>
      <div class="text-h6 text-grey-8">
        Passagens de serviço - Tópicos
      </div>
    </q-card-section>

    <q-separator></q-separator>

    <q-card-section class="q-pa-none">
      <q-table square class="no-shadow"
        title="Passagens"
        row-key="work_report_topic_id"
        :columns="columns"
        :rows="paginationStore.data"
        :filter="filter"
        :rows-per-page-options="[10, 20, 50, 100]"
        v-model:pagination="paginationStore.pagination"
        @request="onRequest"
      >
          <template v-slot:top-right>
            <q-input v-if="show_filter" filled borderless dense debounce="300" v-model="filter" placeholder="Search">
              <template v-slot:append>
                <q-icon name="search"/>
              </template>
            </q-input>

            <q-btn class="q-ml-sm" icon="filter_list" @click="show_filter=!show_filter" flat/>
          </template>

          <template v-slot:header="props">
            <q-tr :props="props">
              <q-th auto-width />
              <q-th
                v-for="col in props.cols"
                :key="col.name"
                :props="props"
              >
                {{ col.label }}
              </q-th>
            </q-tr>
          </template>

          <template v-slot:body="props">
            <q-tr :props="props">
              <q-td auto-width>
                <q-btn size="sm" color="primary" round dense @click="props.expand = !props.expand" :icon="props.expand ? 'remove' : 'add'" />
              </q-td>

              <q-td
                v-for="col in props.cols"
                :key="col.name"
                :props="props"
              >
                <template v-if="col.name === 'Passagem'">
                  {{ col.value }}
                  <q-btn icon="download" size="sm" flat round @click="downloadTopics(props.row)" />
                </template>
                <template v-else>
                  {{ col.value }}
                </template>
              </q-td>
            </q-tr>

            <q-tr v-show="props.expand" :props="props">
              <q-td colspan="100%">
                <span style="display: block; word-wrap: break-word; width: 100ch; white-space: pre-wrap; overflow: auto;">
                  {{props.row.work_report_topic_text}}
              </span>
              </q-td>
            </q-tr>
          </template>
      </q-table>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref, onMounted } from "vue";
import { usePaginationStore } from '@/stores/paginationStore';

export default defineComponent({
  name: 'TopicsComponent',
  setup () {
    const show_filter = ref(false);
    const paginationStore = usePaginationStore(); // Usando a store de paginação
    const filter = ref('');

    const columns = ref([
      {
        name: 'Desembarque',
        required: true,
        label: 'Desembarque',
        align: 'left',
        field: row => row.work_report.work_report_from,
        format: val => {
          const date = new Date(val);
          return `${date.toLocaleDateString('pt-BR')} ${date.toLocaleTimeString('pt-BR')}`;
        },
        sortable: true
      },
      {
        name: 'Passagem',
        required: true,
        label: 'Passagem',
        align: 'left',
        field: row => row.work_report.work_report_docname,
        format: val => `${val}`,
        sortable: true
      },
      {
        name: 'Topico',
        required: true,
        label: 'Tópico',
        align: 'left',
        field: row => row.work_report_topic_title,
        format: val => `${val}`,
        sortable: true
      },
    ]);

    const fetchTopics = () => {
      const endpoint = 'http://localhost:8888/work-report-topics'; // Definindo o endpoint
      const filters = {
        search: filter.value, // Filtro de pesquisa
      };
      paginationStore.fetchData(endpoint, filters); // Usando a função da store para buscar os dados
    };

    const onRequest = (props) => {
      const { page, rowsPerPage } = props.pagination;

      // Atualiza a paginação no store
      paginationStore.pagination.page = page;
      paginationStore.pagination.rowsPerPage = rowsPerPage;

      // Requisita os novos dados
      fetchTopics();
    };

    onMounted(() => {
      fetchTopics(); // Carrega os dados ao montar o componente
    });

    const downloadTopics = (row) => {
      console.log('Download report:', row);
      // Lógica para download do relatório
    };

    return {
      filter,
      show_filter,
      columns,
      paginationStore, // Retorna a store para ser usada no template
      onRequest,
      downloadTopics
    }
  }
});
</script>
