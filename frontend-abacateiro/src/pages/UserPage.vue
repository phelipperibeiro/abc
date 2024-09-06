<template>
  <q-page padding>
    <!-- content -->
    
    <table-actions class="q-mt-lg"></table-actions>
    
  </q-page>
</template>

<script>
import {defineComponent, defineAsyncComponent} from 'vue';

export default defineComponent({
  name: "UserPage",
  components: {
    TableActions: defineAsyncComponent(() => import('components/tables/TableActions.vue')),
  },
})
</script>
