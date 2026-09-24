<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { PageHeader } from '@/components/shared'
import { toast } from 'vue-sonner'
import {
  Palette,
  Loader2,
  Upload,
  Trash2,
  Eye,
  Mail,
  ExternalLink,
  Info,
  Sparkles
} from 'lucide-vue-next'
import { useBrandingStore } from '@/stores/branding'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const brandingStore = useBrandingStore()
const authStore = useAuthStore()

// Super-admin guard — branding is platform-level, not a tenant resource.
// Other roles land here only via the sidebar submenu; we render a clear
// "not allowed" panel instead of silently failing on save.
const isSuperAdmin = computed(() => authStore.user?.is_super_admin || false)

const isLoading = ref(true)
const isSaving = ref(false)
const isUploadingLogo = ref(false)
const isUploadingFavicon = ref(false)

const logoInput = ref<HTMLInputElement | null>(null)
const faviconInput = ref<HTMLInputElement | null>(null)

// Local form state — initialized from store, edited in place, sent as a
// diff to PUT on save. Keeping the form state local avoids the store
// rebinding the DOM mid-edit when the user is typing.
const form = ref({
  brand_name: '',
  brand_tagline: '',
  logo_url: '',
  favicon_url: '',
  primary_color: '142 71% 45%',
  accent_color: '173 80% 40%',
  support_email: '',
  support_url: '',
  footer_text: '',
  login_tagline: ''
})

function syncFromStore() {
  const b = brandingStore.branding
  form.value = {
    brand_name: b.brand_name || '',
    brand_tagline: b.brand_tagline || '',
    logo_url: b.logo_url || '',
    favicon_url: b.favicon_url || '',
    primary_color: b.primary_color || '142 71% 45%',
    accent_color: b.accent_color || '173 80% 40%',
    support_email: b.support_email || '',
    support_url: b.support_url || '',
    footer_text: b.footer_text || '',
    login_tagline: b.login_tagline || ''
  }
}

onMounted(async () => {
  if (!brandingStore.loaded) {
    await brandingStore.load()
  }
  syncFromStore()
  isLoading.value = false
})

watch(
  () => brandingStore.branding,
  () => syncFromStore(),
  { deep: true }
)

// Convenience preview values for the live side-panel
const previewPrimary = computed(() => `hsl(${form.value.primary_color})`)
const previewAccent = computed(() => `hsl(${form.value.accent_color})`)

async function handleSave() {
  if (!isSuperAdmin.value) return
  isSaving.value = true
  try {
    await brandingStore.update({
      brand_name: form.value.brand_name,
      brand_tagline: form.value.brand_tagline,
      logo_url: form.value.logo_url,
      favicon_url: form.value.favicon_url,
      primary_color: form.value.primary_color,
      accent_color: form.value.accent_color,
      support_email: form.value.support_email,
      support_url: form.value.support_url,
      footer_text: form.value.footer_text,
      login_tagline: form.value.login_tagline
    })
    toast.success(t('branding.saved'))
  } catch (err: any) {
    toast.error(err?.response?.data?.message || t('branding.saveFailed'))
  } finally {
    isSaving.value = false
  }
}

async function handleAssetUpload(file: File, type: 'logo' | 'favicon') {
  if (type === 'logo') isUploadingLogo.value = true
  else isUploadingFavicon.value = true
  try {
    const url = await brandingStore.uploadAsset(file, type)
    if (type === 'logo') form.value.logo_url = url
    else form.value.favicon_url = url
    toast.success(t('branding.assetUploaded', { type }))
  } catch (err: any) {
    // Last-resort fallback: surface the raw error so a Vue runtime issue
    // (e.g. an i18n interpolation throw) isn't silently swallowed under a
    // generic "Upload failed" toast.
    toast.error(
      err?.response?.data?.message ||
        err?.message ||
        t('branding.assetUploadFailed')
    )
  } finally {
    if (type === 'logo') isUploadingLogo.value = false
    else isUploadingFavicon.value = false
  }
}

function pickFile(type: 'logo' | 'favicon') {
  const input = type === 'logo' ? logoInput.value : faviconInput.value
  input?.click()
}

function onFileSelected(type: 'logo' | 'favicon', event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    toast.error(t('branding.fileTooLarge'))
    input.value = ''
    return
  }
  handleAssetUpload(file, type).finally(() => {
    input.value = ''
  })
}

function clearAsset(type: 'logo' | 'favicon') {
  if (type === 'logo') form.value.logo_url = ''
  else form.value.favicon_url = ''
}

// Quick HSL picker presets — pure HSL channels, drop into the form on click.
const colorPresets: Array<{ name: string; primary: string; accent: string }> = [
  { name: 'Emerald (default)', primary: '142 71% 45%', accent: '173 80% 40%' },
  { name: 'Royal Blue', primary: '217 91% 60%', accent: '262 83% 58%' },
  { name: 'Sunset Orange', primary: '21 90% 53%', accent: '340 82% 52%' },
  { name: 'Forest Green', primary: '142 76% 36%', accent: '84 81% 44%' },
  { name: 'Royal Purple', primary: '271 81% 56%', accent: '322 81% 52%' },
  { name: 'Slate', primary: '215 25% 27%', accent: '199 89% 48%' }
]

function applyPreset(preset: { primary: string; accent: string }) {
  form.value.primary_color = preset.primary
  form.value.accent_color = preset.accent
}

// CSS color for the picker preview. The form values are HSL channel triplets
// ("142 71% 45%") so we wrap them for display only — the backend stores the
// raw triplet so CSS variables can consume it without parsing.
function hslFromTriplet(triplet: string) {
  if (!triplet) return ''
  return `hsl(${triplet})`
}
</script>

<template>
  <div class="flex flex-col h-full">
    <PageHeader
      :title="$t('branding.title')"
      :subtitle="$t('branding.subtitle')"
      :icon="Palette"
      icon-gradient="bg-gradient-to-br from-fuchsia-500 to-purple-600 shadow-fuchsia-500/20"
    />

    <ScrollArea class="flex-1">
      <div class="p-6 space-y-4 max-w-6xl mx-auto">
        <!-- Non-super-admin guard -->
        <div
          v-if="!isSuperAdmin"
          class="rounded-2xl border border-border bg-card p-6 flex items-start gap-3"
        >
          <Info class="h-5 w-5 text-amber-500 shrink-0 mt-0.5" />
          <div>
            <p class="font-medium">{{ $t('branding.superAdminOnlyTitle') }}</p>
            <p class="text-sm text-muted-foreground mt-1">
              {{ $t('branding.superAdminOnlyDesc') }}
            </p>
          </div>
        </div>

        <div v-else-if="isLoading" class="rounded-2xl border border-border bg-card p-10 flex items-center justify-center">
          <Loader2 class="h-5 w-5 animate-spin text-muted-foreground" />
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <!-- Form column -->
          <div class="lg:col-span-2 space-y-4">
            <!-- Identity card -->
            <Card class="neu-card !rounded-2xl">
              <CardHeader>
                <CardTitle class="flex items-center gap-2">
                  <Sparkles class="h-4 w-4 text-fuchsia-500" />
                  {{ $t('branding.identity') }}
                </CardTitle>
                <CardDescription>{{ $t('branding.identityDesc') }}</CardDescription>
              </CardHeader>
              <CardContent class="space-y-4">
                <div class="space-y-2">
                  <Label for="brand_name">{{ $t('branding.brandName') }}</Label>
                  <Input
                    id="brand_name"
                    v-model="form.brand_name"
                    :placeholder="$t('branding.brandNamePlaceholder')"
                    maxlength="100"
                  />
                  <p class="text-xs text-muted-foreground">
                    {{ $t('branding.brandNameHelp') }}
                  </p>
                </div>

                <div class="space-y-2">
                  <Label for="brand_tagline">{{ $t('branding.brandTagline') }}</Label>
                  <Input
                    id="brand_tagline"
                    v-model="form.brand_tagline"
                    :placeholder="$t('branding.brandTaglinePlaceholder')"
                    maxlength="255"
                  />
                  <p class="text-xs text-muted-foreground">
                    {{ $t('branding.brandTaglineHelp') }}
                  </p>
                </div>
              </CardContent>
            </Card>

            <!-- Logo + Favicon card -->
            <Card class="neu-card !rounded-2xl">
              <CardHeader>
                <CardTitle class="flex items-center gap-2">
                  <Upload class="h-4 w-4 text-fuchsia-500" />
                  {{ $t('branding.assets') }}
                </CardTitle>
                <CardDescription>{{ $t('branding.assetsDesc') }}</CardDescription>
              </CardHeader>
              <CardContent class="space-y-5">
                <!-- Logo -->
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <Label>{{ $t('branding.logo') }}</Label>
                    <span class="text-xs text-muted-foreground">
                      {{ $t('branding.logoHint') }}
                    </span>
                  </div>
                  <div class="flex items-center gap-4">
                    <div
                      class="h-20 w-20 rounded-2xl border border-border bg-background flex items-center justify-center overflow-hidden"
                    >
                      <img
                        v-if="form.logo_url"
                        :src="form.logo_url"
                        :alt="$t('branding.logo')"
                        class="h-full w-full object-contain"
                      />
                      <span
                        v-else
                        class="text-xs text-muted-foreground text-center px-2"
                      >
                        {{ $t('branding.noLogo') }}
                      </span>
                    </div>
                    <div class="flex flex-col gap-2">
                      <input
                        ref="logoInput"
                        type="file"
                        accept="image/png,image/jpeg,image/webp,image/svg+xml,image/x-icon,image/vnd.microsoft.icon"
                        class="hidden"
                        @change="onFileSelected('logo', $event)"
                      />
                      <Button
                        variant="outline"
                        size="sm"
                        class="spring-pressable"
                        :disabled="isUploadingLogo"
                        @click="pickFile('logo')"
                      >
                        <Loader2 v-if="isUploadingLogo" class="mr-2 h-4 w-4 animate-spin" />
                        <Upload v-else class="mr-2 h-4 w-4" />
                        {{ $t('branding.uploadLogo') }}
                      </Button>
                      <Button
                        v-if="form.logo_url"
                        variant="ghost"
                        size="sm"
                        class="text-muted-foreground"
                        @click="clearAsset('logo')"
                      >
                        <Trash2 class="mr-2 h-4 w-4" />
                        {{ $t('branding.removeLogo') }}
                      </Button>
                    </div>
                  </div>
                </div>

                <Separator />

                <!-- Favicon -->
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <Label>{{ $t('branding.favicon') }}</Label>
                    <span class="text-xs text-muted-foreground">
                      {{ $t('branding.faviconHint') }}
                    </span>
                  </div>
                  <div class="flex items-center gap-4">
                    <div
                      class="h-14 w-14 rounded-xl border border-border bg-background flex items-center justify-center overflow-hidden"
                    >
                      <img
                        v-if="form.favicon_url"
                        :src="form.favicon_url"
                        :alt="$t('branding.favicon')"
                        class="h-full w-full object-contain"
                      />
                      <span
                        v-else
                        class="text-[10px] text-muted-foreground text-center px-1"
                      >
                        {{ $t('branding.noFavicon') }}
                      </span>
                    </div>
                    <div class="flex flex-col gap-2">
                      <input
                        ref="faviconInput"
                        type="file"
                        accept="image/png,image/jpeg,image/svg+xml,image/x-icon,image/vnd.microsoft.icon"
                        class="hidden"
                        @change="onFileSelected('favicon', $event)"
                      />
                      <Button
                        variant="outline"
                        size="sm"
                        class="spring-pressable"
                        :disabled="isUploadingFavicon"
                        @click="pickFile('favicon')"
                      >
                        <Loader2 v-if="isUploadingFavicon" class="mr-2 h-4 w-4 animate-spin" />
                        <Upload v-else class="mr-2 h-4 w-4" />
                        {{ $t('branding.uploadFavicon') }}
                      </Button>
                      <Button
                        v-if="form.favicon_url"
                        variant="ghost"
                        size="sm"
                        class="text-muted-foreground"
                        @click="clearAsset('favicon')"
                      >
                        <Trash2 class="mr-2 h-4 w-4" />
                        {{ $t('branding.removeFavicon') }}
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Colors card -->
            <Card class="neu-card !rounded-2xl">
              <CardHeader>
                <CardTitle class="flex items-center gap-2">
                  <Palette class="h-4 w-4 text-fuchsia-500" />
                  {{ $t('branding.colors') }}
                </CardTitle>
                <CardDescription>{{ $t('branding.colorsDesc') }}</CardDescription>
              </CardHeader>
              <CardContent class="space-y-4">
                <div class="grid grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <Label for="primary_color">{{ $t('branding.primaryColor') }}</Label>
                    <div class="flex items-center gap-2">
                      <div
                        class="h-10 w-10 rounded-xl border border-border shrink-0"
                        :style="{ background: hslFromTriplet(form.primary_color) }"
                      />
                      <Input
                        id="primary_color"
                        v-model="form.primary_color"
                        placeholder="142 71% 45%"
                        class="font-mono text-sm"
                      />
                    </div>
                  </div>
                  <div class="space-y-2">
                    <Label for="accent_color">{{ $t('branding.accentColor') }}</Label>
                    <div class="flex items-center gap-2">
                      <div
                        class="h-10 w-10 rounded-xl border border-border shrink-0"
                        :style="{ background: hslFromTriplet(form.accent_color) }"
                      />
                      <Input
                        id="accent_color"
                        v-model="form.accent_color"
                        placeholder="173 80% 40%"
                        class="font-mono text-sm"
                      />
                    </div>
                  </div>
                </div>
                <p class="text-xs text-muted-foreground">
                  {{ $t('branding.colorFormatHint') }}
                </p>

                <Separator />

                <div class="space-y-2">
                  <Label>{{ $t('branding.presets') }}</Label>
                  <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
                    <button
                      v-for="preset in colorPresets"
                      :key="preset.name"
                      type="button"
                      class="spring-pressable flex items-center gap-2.5 rounded-xl border border-border bg-card px-3 py-2 text-left text-sm hover:bg-accent transition-colors"
                      @click="applyPreset(preset)"
                    >
                      <div class="flex">
                        <div
                          class="h-5 w-5 rounded-l-full"
                          :style="{ background: hslFromTriplet(preset.primary) }"
                        />
                        <div
                          class="h-5 w-5 rounded-r-full"
                          :style="{ background: hslFromTriplet(preset.accent) }"
                        />
                      </div>
                      <span class="truncate">{{ preset.name }}</span>
                    </button>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Support + Footer card -->
            <Card class="neu-card !rounded-2xl">
              <CardHeader>
                <CardTitle class="flex items-center gap-2">
                  <Mail class="h-4 w-4 text-fuchsia-500" />
                  {{ $t('branding.support') }}
                </CardTitle>
                <CardDescription>{{ $t('branding.supportDesc') }}</CardDescription>
              </CardHeader>
              <CardContent class="space-y-4">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <Label for="support_email">{{ $t('branding.supportEmail') }}</Label>
                    <Input
                      id="support_email"
                      v-model="form.support_email"
                      type="email"
                      placeholder="support@example.com"
                    />
                  </div>
                  <div class="space-y-2">
                    <Label for="support_url">{{ $t('branding.supportURL') }}</Label>
                    <Input
                      id="support_url"
                      v-model="form.support_url"
                      type="url"
                      placeholder="https://help.example.com"
                    />
                  </div>
                </div>

                <Separator />

                <div class="space-y-2">
                  <Label for="footer_text">{{ $t('branding.footerText') }}</Label>
                  <Textarea
                    id="footer_text"
                    v-model="form.footer_text"
                    :placeholder="$t('branding.footerTextPlaceholder')"
                    rows="2"
                    maxlength="500"
                  />
                </div>

                <div class="space-y-2">
                  <Label for="login_tagline">{{ $t('branding.loginTagline') }}</Label>
                  <Textarea
                    id="login_tagline"
                    v-model="form.login_tagline"
                    :placeholder="$t('branding.loginTaglinePlaceholder')"
                    rows="2"
                    maxlength="500"
                  />
                  <p class="text-xs text-muted-foreground">
                    {{ $t('branding.loginTaglineHelp') }}
                  </p>
                </div>
              </CardContent>
            </Card>

            <!-- Save bar -->
            <div class="flex justify-end gap-2 pb-6">
              <Button
                variant="outline"
                size="sm"
                class="spring-pressable"
                @click="syncFromStore"
                :disabled="isSaving"
              >
                {{ $t('branding.reset') }}
              </Button>
              <Button
                size="sm"
                class="spring-pressable"
                :style="{ background: previewPrimary, borderColor: previewPrimary }"
                :disabled="isSaving"
                @click="handleSave"
              >
                <Loader2 v-if="isSaving" class="mr-2 h-4 w-4 animate-spin" />
                {{ $t('branding.save') }}
              </Button>
            </div>
          </div>

          <!-- Preview column -->
          <div class="lg:col-span-1">
            <div class="sticky top-2 space-y-4">
              <Card class="neu-card !rounded-2xl">
                <CardHeader>
                  <CardTitle class="flex items-center gap-2 text-sm">
                    <Eye class="h-4 w-4 text-fuchsia-500" />
                    {{ $t('branding.preview') }}
                  </CardTitle>
                  <CardDescription>{{ $t('branding.previewDesc') }}</CardDescription>
                </CardHeader>
                <CardContent class="space-y-4">
                  <!-- Sidebar preview -->
                  <div class="rounded-xl border border-border bg-background p-3 space-y-2">
                    <div class="flex items-center gap-2.5 px-1">
                      <div
                        class="h-8 w-8 rounded-[10px] overflow-hidden flex items-center justify-center"
                        :style="{ background: `linear-gradient(135deg, ${previewPrimary}, ${previewAccent})` }"
                      >
                        <img
                          v-if="form.logo_url"
                          :src="form.logo_url"
                          class="h-full w-full object-contain bg-white"
                        />
                        <span v-else class="text-white text-[10px] font-semibold">
                          {{ (form.brand_name || 'W').slice(0, 1) }}
                        </span>
                      </div>
                      <div class="leading-none">
                        <div class="text-sm font-semibold truncate">
                          {{ form.brand_name || 'Whatomate' }}
                        </div>
                        <div class="text-[9px] uppercase tracking-wider text-muted-foreground">
                          {{ form.brand_tagline || 'WhatsApp Platform' }}
                        </div>
                      </div>
                    </div>
                    <div class="h-7 rounded-lg bg-accent/40 px-3 flex items-center text-xs">
                      Dashboard
                    </div>
                    <div
                      class="h-7 rounded-lg px-3 flex items-center text-xs"
                      :style="{ background: `${previewPrimary}22`, color: previewPrimary }"
                    >
                      Chat (active)
                    </div>
                  </div>

                  <!-- Login preview -->
                  <div
                    class="rounded-xl border border-border p-4 text-center"
                    :style="{
                      background: `linear-gradient(135deg, ${previewPrimary}1a, ${previewAccent}1a)`
                    }"
                  >
                    <div
                      class="mx-auto h-12 w-12 rounded-2xl overflow-hidden flex items-center justify-center mb-2"
                      :style="{ background: `linear-gradient(135deg, ${previewPrimary}, ${previewAccent})` }"
                    >
                      <img
                        v-if="form.logo_url"
                        :src="form.logo_url"
                        class="h-full w-full object-contain bg-white"
                      />
                      <span v-else class="text-white text-lg font-semibold">
                        {{ (form.brand_name || 'W').slice(0, 1) }}
                      </span>
                    </div>
                    <div class="font-semibold">
                      {{ $t('auth.welcomeTitle') }} {{ form.brand_name || 'Whatomate' }}
                    </div>
                    <div class="text-xs text-muted-foreground mt-0.5">
                      {{ form.login_tagline || $t('auth.welcomeSubtitle') }}
                    </div>
                    <div class="mt-3 h-8 rounded-lg border border-border bg-background px-3 flex items-center text-xs text-muted-foreground">
                      email@example.com
                    </div>
                    <button
                      class="mt-2 h-8 w-full rounded-lg text-xs font-medium text-white"
                      :style="{ background: previewPrimary }"
                    >
                      {{ $t('auth.signIn') }}
                    </button>
                  </div>

                  <!-- Support + Footer preview -->
                  <div class="rounded-xl border border-border bg-background p-3 space-y-2 text-xs">
                    <div class="font-medium">{{ $t('branding.previewSupport') }}</div>
                    <a
                      v-if="form.support_email"
                      :href="`mailto:${form.support_email}`"
                      class="flex items-center gap-1.5 text-muted-foreground hover:text-foreground"
                    >
                      <Mail class="h-3 w-3" /> {{ form.support_email }}
                    </a>
                    <a
                      v-if="form.support_url"
                      :href="form.support_url"
                      target="_blank"
                      rel="noopener"
                      class="flex items-center gap-1.5 text-muted-foreground hover:text-foreground"
                    >
                      <ExternalLink class="h-3 w-3" /> {{ form.support_url }}
                    </a>
                    <div v-if="form.footer_text" class="pt-2 border-t border-border text-muted-foreground">
                      {{ form.footer_text }}
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </ScrollArea>
  </div>
</template>