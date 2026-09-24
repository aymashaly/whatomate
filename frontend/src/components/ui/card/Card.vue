<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: string
}>()
</script>

<template>
  <!--
    Soft-Tile surface. Same color as the page canvas, with paired
    shadows that make the card appear extruded. Subtle hover lift
    when paired with `.transition`.
  -->
  <div
    :class="cn(
      'neu-card rounded-2xl border border-transparent text-card-foreground transition',
      props.class
    )"
  >
    <slot />
  </div>
</template>
