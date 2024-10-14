<template>
  <q-card class="no-shadow" bordered>
    <q-card-section>
      <div class="text-h6 text-grey-8">
        Passagens de serviço - Busca Avançada
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
                  <span v-html="col.value"></span>
                </template>
              </q-td>
            </q-tr>
          </template>

      </q-table>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref, watch, onMounted } from "vue";
import { usePaginationStore } from '@/stores/advancedSearchPaginationStore';

export default defineComponent({
  name: 'TopicsComponent',
  setup () {
    const show_filter = ref(false);
    const paginationStore = usePaginationStore(); // Usando a store de paginação
    const filter = ref('');

    const columns = ref([
      {
        name: 'Rank',
        required: true,
        label: 'Rank',
        align: 'left',
        field: row => row.rank,
        sortable: true
      },
      {
        name: 'Passagem',
        required: true,
        label: 'Passagem',
        align: 'left',
        field: row => row.work_report_docname,
        format: val => `${val}`,
        sortable: true
      },
      {
        name: 'Topico',
        required: true,
        label: 'Tópico',
        align: 'left',
        field: row => row.highlight_title,
        sortable: true
      },
    ]);

    const fetchTopicsAdvSearch = () => {
      const endpoint = 'http://localhost:8888/work-report-topics/adv-search';
      const filters = {
        search: filter.value, // Filtro de pesquisa
      };
      paginationStore.fetchData(endpoint, filters);
    };

    const onRequest = (props) => {
      const { page, rowsPerPage } = props.pagination;

      // Atualiza a paginação no store
      paginationStore.pagination.page = page;
      paginationStore.pagination.rowsPerPage = rowsPerPage;

      // Requisita os novos dados
      fetchTopicsAdvSearch();
    };

    onMounted(() => {
      fetchTopicsAdvSearch(); // Carrega os dados ao montar o componente
    });

    watch(filter, (newValue, oldValue) => {
      console.log('Filter changed from', oldValue, 'to', newValue);
      fetchTopicsAdvSearch();
    });

    const downloadTopics = (row) => {
      console.log('Download report:', row);
      // Lógica para download do relatório
    };

    return {
      filter,
      show_filter,
      columns,
      paginationStore,
      onRequest,
      downloadTopics
    }
  }
});
</script>
