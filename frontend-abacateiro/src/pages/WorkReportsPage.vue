<template>

    <q-card class="fit no-shadow" bordered>

      <q-tabs
        v-model="tab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="justify"
      >
      <q-tab name="reports" :class="tab=='contact'?'text-blue':''" icon="receipt_long" label="Passagens"/>
      <q-tab name="topics" :class="tab=='message'?'text-blue':''" icon="list_alt" label="Tópicos" />
      <q-tab name="advanced-search" :class="tab=='advanced-search'?'text-blue':''" icon="filter_list" label="Busca Avançada" />

      </q-tabs>

      <q-separator/>

      <q-tab-panels v-model="tab" animated>
        <q-tab-panel name="reports" class="q-pa-sm">
          <ReportsComponent/>
        </q-tab-panel>

        <q-tab-panel name="topics" class="q-pa-sm">
          <TopicsComponent/>
        </q-tab-panel>

        <q-tab-panel name="advanced-search" class="q-pa-sm">
          <AdvancedSearchComponent/>
        </q-tab-panel>

      </q-tab-panels>

    </q-card>

</template>

<script>
import { defineComponent, defineAsyncComponent } from "vue";
import { ref } from 'vue'
import { EventBus } from "@/plugins/eventBus";

export default defineComponent({
  name: "WorkReportsPage",
  components: {
    ReportsComponent: defineAsyncComponent(() => import('components/workReports/tables/Reports.vue')),
    AdvancedSearchComponent: defineAsyncComponent(() => import('components/workReports/tables/AdvancedSearch.vue')),
    TopicsComponent: defineAsyncComponent(() => import('components/workReports/tables/Topics.vue'))
  },
  setup() {
    return {
      tab: ref('reports')
    };
  },
});
</script>

<style scoped>
.button-container {
  display: flex;
  justify-content: flex-end;
}
</style>
