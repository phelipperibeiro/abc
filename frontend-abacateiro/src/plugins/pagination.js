import axios from 'axios';

export async function paginate(endpoint, page = 1, perPage = 10, filters = {}) {

  try {

    const response = await axios.get(endpoint, {
      params: {
        page,
        page_size: perPage,
        ...filters,
      },
    });

    // TODO: preciso ajustar a questao da variavel $key, pode ser que os itens não esteja na primeira posição
    const $key = Object.keys(response.data)[0];

    const metadata = response.data.metadata;
    const itens = response.data[$key];

    return {
      data: itens, // Dados paginados
      currentPage: metadata.current_page, // Página atual
      perPage: metadata.page_size, // Itens por página
      totalPages: metadata.last_page, // Última página (total de páginas)
      totalItems: metadata.total_records, // Total de registros
    };

  } catch (error) {
    console.error("Erro ao buscar dados:", error);
    throw error;
  }
}
