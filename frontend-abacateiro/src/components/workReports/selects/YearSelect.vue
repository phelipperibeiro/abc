<template>
  <q-select
    filled
    borderless
    dense
    v-model="model"
    :options="yearOptions"
    label="Ano"
    class="q-mr-sm"
  />
</template>

<script>
import { defineComponent, ref } from 'vue';

const generateYearOptions = (startYear) => {
  const currentYear = new Date().getFullYear();
  const years = [{ label: 'Selecione', value: '0' }];
  for (let year = startYear; year <= currentYear; year++) {
    years.push({ label: year, value: year });
  }
  return years;
};

export default defineComponent({

  name: 'YearSelect',

  emits: ['update:modelValue'],

  setup(_, { emit }) {

    const model = ref('Selecione');

    const yearOptions = ref(generateYearOptions(2018));

    const updateModel = (value) => {
      model.value = value;
      alert(value);
      emit('update:modelValue', value);
    };

    return {
      model,
      yearOptions,
      updateModel
    };

  }
});
</script>
