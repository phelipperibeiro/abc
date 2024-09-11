<template>
  <q-page class="flex flex-center">
    <img
      alt="Quasar logo"
      src="~assets/quasar-logo-vertical.svg"
      style="width: 200px; height: 200px"
    >
  </q-page>
</template>

<script setup>

import UserTable from 'src/components/table/UserTable.vue';
defineOptions({
  name: 'IndexPage'
});
</script>
