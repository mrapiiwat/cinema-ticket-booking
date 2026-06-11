<template>
  <div class="min-h-screen bg-stone-50 flex overflow-hidden">
    <AdminSidebar />

    <div class="flex-1 flex flex-col min-w-0 h-screen">
      <AdminNavbar />

      <div class="flex-1 overflow-y-auto">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import AdminSidebar from '../components/AdminSidebar.vue'
import AdminNavbar from '../components/AdminNavbar.vue'
</script>
