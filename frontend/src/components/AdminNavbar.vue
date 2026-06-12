<template>
  <header class="bg-white border-b border-stone-100 px-8 py-4 flex items-center justify-end">
    <div class="flex items-center gap-4">
      <span class="text-xs text-stone-400 font-mono bg-stone-50 border border-stone-100 px-2 py-1 rounded-md">
        {{ user?.email || 'admin@cinepass.local' }}
      </span>
      <span class="text-xs text-emerald-700 font-mono bg-emerald-50 border border-emerald-100 px-2 py-1 rounded-md">
        Role: ADMIN
      </span>
    </div>
  </header>
</template>

<script setup>
import { getSession } from '../api/client'

const user = getSession()?.user
</script>
