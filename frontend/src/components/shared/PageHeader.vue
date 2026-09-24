<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { ArrowLeft } from 'lucide-vue-next'
import type { Component } from 'vue'

defineProps<{
  title: string
  description?: string
  icon?: Component
  iconGradient?: string
  backLink?: string
  breadcrumbs?: Array<{ label: string; href?: string }>
}>()
</script>

<template>
  <!--
    Floating header — minimal glass + neumorphic icon mark.
    No ambient animation. Subtle top edge highlight only.
  -->
  <header
    class="glass mx-2 mt-2 flex h-14 items-center rounded-2xl px-4 light:bg-white/85"
  >
    <RouterLink v-if="backLink" :to="backLink">
      <Button
        variant="ghost"
        size="icon"
        class="spring-pressable mr-3 h-9 w-9 rounded-xl text-muted-foreground hover:text-foreground hover:bg-black/[0.04] light:hover:bg-gray-100"
      >
        <ArrowLeft class="h-4 w-4" />
      </Button>
    </RouterLink>
    <div
      v-if="icon"
      class="relative h-9 w-9 shrink-0 overflow-hidden rounded-xl bg-gradient-to-br from-emerald-400 to-cyan-500 flex items-center justify-center mr-3 shadow-sm ring-1 ring-white/[0.10] light:ring-black/[0.05]"
      :class="iconGradient"
    >
      <component :is="icon" class="h-4 w-4 text-white relative" />
    </div>
    <div class="flex-1 min-w-0">
      <h1 class="display text-[14px] tracking-tight truncate text-foreground">{{ title }}</h1>
      <template v-if="breadcrumbs?.length">
        <Breadcrumb>
          <BreadcrumbList>
            <template v-for="(crumb, index) in breadcrumbs" :key="index">
              <BreadcrumbItem>
                <BreadcrumbLink v-if="crumb.href" :href="crumb.label">
                  {{ crumb.label }}
                </BreadcrumbLink>
                <BreadcrumbPage v-else>{{ crumb.label }}</BreadcrumbPage>
              </BreadcrumbItem>
              <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" />
            </template>
          </BreadcrumbList>
        </Breadcrumb>
      </template>
      <p v-else-if="description" class="text-[11.5px] text-muted-foreground mt-0.5 truncate">
        {{ description }}
      </p>
    </div>
    <slot name="actions" />
  </header>
</template>
