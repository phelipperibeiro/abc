<template>
  <q-select
    filled
    borderless
    dense
    v-model="model"
    :options="unitOptions"
    label="Unidade"
    class="q-mr-sm"
  />
</template>

<script>
import { defineComponent, ref } from 'vue';

export default defineComponent({

  name: 'UnitSelect',

  emits: ['update:modelValue'],

  setup(_, { emit }) {

    const model = ref('Selecione');

    // TODO: carregar as unidades do banco de dados, por enquanto está fixo uns fake para testar.
    // vefificar a melhor forma de fazer isso, se é com um service ou com um store.
    const unitOptions = ref([
      { label: 'Selecione', value: '0' },
      { label: 'Unidade 1', value: 'Unidade 1' },
      { label: 'Unidade 2', value: 'Unidade 2' },
      { label: 'Unidade 3', value: 'Unidade 3' }
    ]);

    const updateModel = (value) => {
      model.value = value;
      emit('update:modelValue', value);
    };

    return {
      model,
      unitOptions,
      updateModel
    };
  }
});
</script>
