import { defineStore } from 'pinia';
import { ref } from 'vue';
import { paginate } from '@/plugins/pagination';

export const usePaginationStore = defineStore('pagination', () => {

  const data = ref([]);

  const pagination = ref({
    page: 1,
    rowsPerPage: 10,
    totalPages: 1,
    rowsNumber: 0, // Total de registros
  });

  const loading = ref(false);

  const fetchData = async (endpoint, filters = {}) => {

    loading.value = true;

    try {

      const result = await paginate(
        endpoint,
        pagination.value.page,
        pagination.value.rowsPerPage,
        filters
      );

      data.value = result.data;
      pagination.value.page = result.currentPage;
      pagination.value.rowsPerPage = result.perPage;
      pagination.value.totalPages = result.totalPages;
      pagination.value.rowsNumber = result.totalItems;

    } catch (error) {

      console.error("Erro ao carregar os dados:", error);

    } finally {

      loading.value = false;
    }

  };

  return {
    data,
    pagination,
    loading,
    fetchData,
  };
});
