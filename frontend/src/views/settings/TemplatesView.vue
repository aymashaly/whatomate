<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import DOMPurify from 'dompurify'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command'
import { PageHeader, SearchInput, DataTable, DeleteConfirmDialog, type Column } from '@/components/shared'
import { api, templatesService } from '@/services/api'
import { useOrganizationsStore } from '@/stores/organizations'
import { toast } from 'vue-sonner'
import { Plus, RefreshCw, FileText, Eye, Pencil, Trash2, Loader2, MessageSquare, Image, FileIcon, Video, X, Check, AlertCircle, Send, Upload, ChevronsUpDown } from 'lucide-vue-next'
import { getErrorMessage } from '@/lib/api-utils'
import { useDebounceFn } from '@vueuse/core'

const { t } = useI18n()

interface WhatsAppAccount {
  id: string
  name: string
  phone_id: string
}

interface Template {
  id: string
  whatsapp_account: string
  meta_template_id: string
  name: string
  display_name: string
  language: string
  category: string
  status: string
  header_type: string
  header_content: string
  body_content: string
  footer_content: string
  buttons: any[]
  sample_values: any[]
  created_at: string
  updated_at: string
}

const organizationsStore = useOrganizationsStore()

const templates = ref<Template[]>([])
const accounts = ref<WhatsAppAccount[]>([])
const isLoading = ref(true)
const isSyncing = ref(false)
const searchQuery = ref('')
const selectedAccount = ref<string>(localStorage.getItem('templates_selected_account') || 'all')

// Dialog state
const isDialogOpen = ref(false)
const isSubmitting = ref(false)
const editingTemplate = ref<Template | null>(null)
const isPreviewOpen = ref(false)
const previewTemplate = ref<Template | null>(null)
const deleteDialogOpen = ref(false)
const templateToDelete = ref<Template | null>(null)
const publishDialogOpen = ref(false)
const templateToPublish = ref<Template | null>(null)

// Header media upload state
const headerMediaFile = ref<File | null>(null)
const headerMediaUploading = ref(false)
const headerMediaHandle = ref('')
const headerMediaFilename = ref('')

const formData = ref({
  whatsapp_account: '',
  name: '',
  display_name: '',
  language: 'en',
  category: 'UTILITY',
  header_type: 'NONE',
  header_content: '',
  body_content: '',
  footer_content: '',
  buttons: [] as any[],
  sample_values: [] as any[]
})

// Pagination state
const currentPage = ref(1)
const totalItems = ref(0)
const pageSize = 20

const columns = computed<Column<Template>[]>(() => [
  { key: 'name', label: t('templates.name'), sortable: true },
  { key: 'category', label: t('templates.category'), sortable: true },
  { key: 'status', label: t('templates.status'), sortable: true },
  { key: 'language', label: t('templates.language'), sortable: true },
  { key: 'header_type', label: t('templates.header') },
  { key: 'actions', label: t('common.actions'), align: 'right' },
])

const sortKey = ref('name')
const sortDirection = ref<'asc' | 'desc'>('asc')

const languages = [
  { code: 'af', name: 'Afrikaans' },
  { code: 'sq', name: 'Albanian' },
  { code: 'ar', name: 'Arabic' },
  { code: 'az', name: 'Azerbaijani' },
  { code: 'bn', name: 'Bengali' },
  { code: 'bg', name: 'Bulgarian' },
  { code: 'ca', name: 'Catalan' },
  { code: 'zh_CN', name: 'Chinese (CHN)' },
  { code: 'zh_HK', name: 'Chinese (HKG)' },
  { code: 'zh_TW', name: 'Chinese (TAI)' },
  { code: 'hr', name: 'Croatian' },
  { code: 'cs', name: 'Czech' },
  { code: 'da', name: 'Danish' },
  { code: 'nl', name: 'Dutch' },
  { code: 'en', name: 'English' },
  { code: 'en_GB', name: 'English (UK)' },
  { code: 'en_US', name: 'English (US)' },
  { code: 'et', name: 'Estonian' },
  { code: 'fil', name: 'Filipino' },
  { code: 'fi', name: 'Finnish' },
  { code: 'fr', name: 'French' },
  { code: 'ka', name: 'Georgian' },
  { code: 'de', name: 'German' },
  { code: 'el', name: 'Greek' },
  { code: 'gu', name: 'Gujarati' },
  { code: 'ha', name: 'Hausa' },
  { code: 'he', name: 'Hebrew' },
  { code: 'hi', name: 'Hindi' },
  { code: 'hu', name: 'Hungarian' },
  { code: 'id', name: 'Indonesian' },
  { code: 'ga', name: 'Irish' },
  { code: 'it', name: 'Italian' },
  { code: 'ja', name: 'Japanese' },
  { code: 'kn', name: 'Kannada' },
  { code: 'kk', name: 'Kazakh' },
  { code: 'rw_RW', name: 'Kinyarwanda' },
  { code: 'ko', name: 'Korean' },
  { code: 'ky_KG', name: 'Kyrgyz (Kyrgyzstan)' },
  { code: 'lo', name: 'Lao' },
  { code: 'lv', name: 'Latvian' },
  { code: 'lt', name: 'Lithuanian' },
  { code: 'mk', name: 'Macedonian' },
  { code: 'ms', name: 'Malay' },
  { code: 'ml', name: 'Malayalam' },
  { code: 'mr', name: 'Marathi' },
  { code: 'nb', name: 'Norwegian (Bokmål)' },
  { code: 'fa', name: 'Persian' },
  { code: 'pl', name: 'Polish' },
  { code: 'pt_BR', name: 'Portuguese (BR)' },
  { code: 'pt_PT', name: 'Portuguese (POR)' },
  { code: 'pa', name: 'Punjabi' },
  { code: 'ro', name: 'Romanian' },
  { code: 'ru', name: 'Russian' },
  { code: 'sr', name: 'Serbian' },
  { code: 'sk', name: 'Slovak' },
  { code: 'sl', name: 'Slovenian' },
  { code: 'es', name: 'Spanish' },
  { code: 'es_AR', name: 'Spanish (ARG)' },
  { code: 'es_MX', name: 'Spanish (MEX)' },
  { code: 'es_ES', name: 'Spanish (SPA)' },
  { code: 'sw', name: 'Swahili' },
  { code: 'sv', name: 'Swedish' },
  { code: 'ta', name: 'Tamil' },
  { code: 'te', name: 'Telugu' },
  { code: 'th', name: 'Thai' },
  { code: 'tr', name: 'Turkish' },
  { code: 'uk', name: 'Ukrainian' },
  { code: 'ur', name: 'Urdu' },
  { code: 'uz', name: 'Uzbek' },
  { code: 'vi', name: 'Vietnamese' },
  { code: 'zu', name: 'Zulu' },
]

const languageSelectorOpen = ref(false)

function getLanguageName(code: string): string {
  return languages.find(l => l.code === code)?.name || code
}

const categories = [
  { value: 'UTILITY', label: 'Utility', description: 'Order updates, account alerts' },
  { value: 'MARKETING', label: 'Marketing', description: 'Promotions, offers' },
  { value: 'AUTHENTICATION', label: 'Authentication', description: 'OTP, verification codes' },
]

const headerTypes = [
  { value: 'NONE', label: 'None' },
  { value: 'TEXT', label: 'Text' },
  { value: 'IMAGE', label: 'Image' },
  { value: 'VIDEO', label: 'Video' },
  { value: 'DOCUMENT', label: 'Document' },
]

// Refetch data when organization changes
watch(() => organizationsStore.selectedOrgId, async () => {
  await fetchAccounts()
  await fetchTemplates()
})

onMounted(async () => {
  await fetchAccounts()
  await fetchTemplates()
})

async function fetchAccounts() {
  try {
    const response = await api.get('/accounts')
    accounts.value = response.data.data?.accounts || []
    // Validate stored account still exists, fallback to 'all' if not
    if (selectedAccount.value !== 'all' && !accounts.value.some(a => a.name === selectedAccount.value)) {
      selectedAccount.value = 'all'
      localStorage.setItem('templates_selected_account', 'all')
    }
  } catch (error) {
    console.error('Failed to fetch accounts:', error)
  }
}

function onAccountChange(value: string | number | bigint | Record<string, any> | null) {
  if (typeof value !== 'string') return
  localStorage.setItem('templates_selected_account', value)
  currentPage.value = 1
  fetchTemplates()
}

async function fetchTemplates() {
  isLoading.value = true
  try {
    const response = await templatesService.list({
      account: selectedAccount.value !== 'all' ? selectedAccount.value : undefined,
      search: searchQuery.value || undefined,
      page: currentPage.value,
      limit: pageSize
    })
    const data = (response.data as any).data || response.data
    templates.value = data.templates || []
    totalItems.value = data.total ?? templates.value.length
  } catch (error: any) {
    console.error('Failed to fetch templates:', error)
    toast.error(t('common.failedLoad', { resource: t('resources.templates') }))
    templates.value = []
  } finally {
    isLoading.value = false
  }
}

// Debounced search
const debouncedSearch = useDebounceFn(() => {
  currentPage.value = 1
  fetchTemplates()
}, 300)

watch(searchQuery, () => debouncedSearch())

function handlePageChange(page: number) {
  currentPage.value = page
  fetchTemplates()
}

async function syncTemplates() {
  if (!selectedAccount.value || selectedAccount.value === 'all') {
    toast.error(t('templates.selectAccountFirst'))
    return
  }

  isSyncing.value = true
  try {
    const response = await api.post('/templates/sync', {
      whatsapp_account: selectedAccount.value
    })
    toast.success(response.data.data.message || t('templates.syncSuccess'))
    await fetchTemplates()
  } catch (error) {
    toast.error(getErrorMessage(error, t('templates.syncFailed')))
  } finally {
    isSyncing.value = false
  }
}

function openCreateDialog() {
  editingTemplate.value = null
  formData.value = {
    whatsapp_account: (selectedAccount.value && selectedAccount.value !== 'all') ? selectedAccount.value : (accounts.value[0]?.name || ''),
    name: '',
    display_name: '',
    language: 'en',
    category: 'UTILITY',
    header_type: 'NONE',
    header_content: '',
    body_content: '',
    footer_content: '',
    buttons: [],
    sample_values: []
  }
  // Reset header media state
  headerMediaFile.value = null
  headerMediaHandle.value = ''
  headerMediaFilename.value = ''
  isDialogOpen.value = true
}

function openEditDialog(template: Template) {
  editingTemplate.value = template
  formData.value = {
    whatsapp_account: template.whatsapp_account,
    name: template.name,
    display_name: template.display_name,
    language: template.language,
    category: template.category,
    header_type: template.header_type || 'NONE',
    header_content: template.header_content || '',
    body_content: template.body_content,
    footer_content: template.footer_content || '',
    buttons: template.buttons || [],
    sample_values: template.sample_values || []
  }
  // Reset header media state (will show existing handle if present)
  headerMediaFile.value = null
  headerMediaHandle.value = template.header_content || ''
  headerMediaFilename.value = ''
  isDialogOpen.value = true
}

function openPreview(template: Template) {
  previewTemplate.value = template
  isPreviewOpen.value = true
}

async function saveTemplate() {
  if (!formData.value.name.trim() || !formData.value.body_content.trim()) {
    toast.error(t('templates.nameBodyRequired'))
    return
  }

  if (!formData.value.whatsapp_account) {
    toast.error(t('templates.selectAccountRequired'))
    return
  }

  isSubmitting.value = true
  try {
    if (editingTemplate.value) {
      await api.put(`/templates/${editingTemplate.value.id}`, formData.value)
      toast.success(t('common.updatedSuccess', { resource: t('resources.Template') }))
    } else {
      await api.post('/templates', formData.value)
      toast.success(t('common.createdSuccess', { resource: t('resources.Template') }))
    }
    isDialogOpen.value = false
    await fetchTemplates()
  } catch (error) {
    toast.error(getErrorMessage(error, t('common.failedSave', { resource: t('resources.template') })))
  } finally {
    isSubmitting.value = false
  }
}

function openDeleteDialog(template: Template) {
  templateToDelete.value = template
  deleteDialogOpen.value = true
}

async function confirmDeleteTemplate() {
  if (!templateToDelete.value) return

  try {
    await api.delete(`/templates/${templateToDelete.value.id}`)
    toast.success(t('common.deletedSuccess', { resource: t('resources.Template') }))
    deleteDialogOpen.value = false
    templateToDelete.value = null
    await fetchTemplates()
  } catch (error) {
    toast.error(getErrorMessage(error, t('common.failedDelete', { resource: t('resources.template') })))
  }
}

const publishingTemplateId = ref<string | null>(null)

function openPublishDialog(template: Template) {
  templateToPublish.value = template
  publishDialogOpen.value = true
}

async function confirmPublishTemplate() {
  if (!templateToPublish.value) return

  publishingTemplateId.value = templateToPublish.value.id
  try {
    const response = await api.post(`/templates/${templateToPublish.value.id}/publish`)
    toast.success(response.data.data?.message || t('templates.publishSuccess'))
    publishDialogOpen.value = false
    templateToPublish.value = null
    await fetchTemplates()
  } catch (error) {
    toast.error(getErrorMessage(error, t('templates.publishFailed')), { duration: 8000 })
  } finally {
    publishingTemplateId.value = null
  }
}

// Dark-first: default is dark mode, light: prefix for light mode
function getStatusBadgeClass(status: string) {
  switch (status) {
    case 'APPROVED':
      return 'bg-green-900 text-green-300 light:bg-green-100 light:text-green-800'
    case 'PENDING':
      return 'bg-yellow-900 text-yellow-300 light:bg-yellow-100 light:text-yellow-800'
    case 'REJECTED':
      return 'bg-red-900 text-red-300 light:bg-red-100 light:text-red-800'
    case 'DRAFT':
      return 'bg-gray-800 text-gray-300 light:bg-gray-100 light:text-gray-800'
    default:
      return 'bg-gray-800 text-gray-300 light:bg-gray-100 light:text-gray-800'
  }
}

function getCategoryBadgeClass(category: string) {
  switch (category) {
    case 'UTILITY':
      return 'bg-blue-900 text-blue-300 light:bg-blue-100 light:text-blue-800'
    case 'MARKETING':
      return 'bg-purple-900 text-purple-300 light:bg-purple-100 light:text-purple-800'
    case 'AUTHENTICATION':
      return 'bg-orange-900 text-orange-300 light:bg-orange-100 light:text-orange-800'
    default:
      return 'bg-gray-800 text-gray-300 light:bg-gray-100 light:text-gray-800'
  }
}

function getHeaderIcon(type: string) {
  switch (type) {
    case 'IMAGE':
      return Image
    case 'VIDEO':
      return Video
    case 'DOCUMENT':
      return FileIcon
    default:
      return MessageSquare
  }
}


// Extract all parameter names (both positional {{1}} and named {{name}})
function extractParamNames(content: string): string[] {
  const matches = content.match(/\{\{([^}]+)\}\}/g) || []
  const seen = new Set<string>()
  const names: string[] = []
  for (const m of matches) {
    const name = m.replace(/[{}]/g, '').trim()
    if (name && !seen.has(name)) {
      seen.add(name)
      names.push(name)
    }
  }
  return names
}

// Get variable names from body content (supports both {{1}} and {{name}})
const bodyVariables = computed(() => {
  return extractParamNames(formData.value.body_content)
})

// Get variable names from header content
const headerVariables = computed(() => {
  if (formData.value.header_type !== 'TEXT') return []
  return extractParamNames(formData.value.header_content)
})

// Button types for template
const buttonTypes = [
  { value: 'QUICK_REPLY', label: 'Quick Reply', description: 'Simple reply button' },
  { value: 'URL', label: 'URL', description: 'Opens a website' },
  { value: 'PHONE_NUMBER', label: 'Phone Number', description: 'Calls a number' },
]

function addButton() {
  if (formData.value.buttons.length >= 3) {
    toast.error(t('templates.maxButtons'))
    return
  }
  formData.value.buttons.push({
    type: 'QUICK_REPLY',
    text: ''
  })
}

function removeButton(index: number) {
  formData.value.buttons.splice(index, 1)
}

// Handle header media file selection
function onHeaderMediaFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files && input.files.length > 0) {
    headerMediaFile.value = input.files[0]
    headerMediaFilename.value = input.files[0].name
    // Clear previous handle when new file is selected
    headerMediaHandle.value = ''
    formData.value.header_content = ''
  }
}

// Upload header media file to Meta
async function uploadHeaderMedia() {
  if (!headerMediaFile.value) {
    toast.error(t('templates.selectAccountFirst'))
    return
  }

  if (!formData.value.whatsapp_account) {
    toast.error(t('templates.selectAccountFirst'))
    return
  }

  headerMediaUploading.value = true
  try {
    const response = await templatesService.uploadMedia(formData.value.whatsapp_account, headerMediaFile.value)
    const data = response.data.data
    headerMediaHandle.value = data.handle
    formData.value.header_content = data.handle
    toast.success(t('templates.mediaUploadedSuccess'))
  } catch (error) {
    toast.error(getErrorMessage(error, t('templates.uploadFailed')))
  } finally {
    headerMediaUploading.value = false
  }
}

// Get accepted file types for header type
function getAcceptedFileTypes(): string {
  switch (formData.value.header_type) {
    case 'IMAGE':
      return 'image/jpeg,image/png'
    case 'VIDEO':
      return 'video/mp4'
    case 'DOCUMENT':
      return 'application/pdf'
    default:
      return '*/*'
  }
}

function getSampleValue(component: string, paramName: string): string {
  const sample = formData.value.sample_values.find(
    (s: any) => s.component === component && s.param_name === paramName
  )
  return sample?.value || ''
}

function setSampleValue(component: string, paramName: string, value: string) {
  const existingIndex = formData.value.sample_values.findIndex(
    (s: any) => s.component === component && s.param_name === paramName
  )
  if (existingIndex >= 0) {
    formData.value.sample_values[existingIndex].value = value
  } else {
    formData.value.sample_values.push({ component, param_name: paramName, value })
  }
}

function formatVariableLabel(paramName: string): string {
  return `{{${paramName}}}`
}

// Format template preview with sample values (sanitized to prevent XSS)
function formatPreview(text: string, samples: any[]): string {
  // Sanitize the base text first
  let result = DOMPurify.sanitize(text, { ALLOWED_TAGS: [] })

  // Handle named parameters with param_name field
  samples.forEach((sample) => {
    if (sample && sample.param_name && sample.value) {
      const sanitizedSample = DOMPurify.sanitize(String(sample.value), { ALLOWED_TAGS: [] })
      result = result.replace(`{{${sample.param_name}}}`, `<span class="bg-green-900 light:bg-green-100 px-1 rounded">${sanitizedSample}</span>`)
    }
  })

  // Replace remaining variables (both named and positional)
  result = result.replace(/\{\{([^}]+)\}\}/g, '<span class="bg-yellow-900 light:bg-yellow-100 px-1 rounded">{{$1}}</span>')
  return result
}
</script>

<template>
  <div class="flex flex-col h-full bg-[#0a0a0b] light:bg-gray-50">
    <PageHeader :title="$t('templates.title')" :subtitle="$t('templates.subtitle')" :icon="FileText" icon-gradient="bg-gradient-to-br from-blue-500 to-cyan-600 shadow-blue-500/20">
      <template #actions>
        <Button variant="outline" size="sm" @click="syncTemplates" :disabled="isSyncing || !selectedAccount || selectedAccount === 'all'">
          <Loader2 v-if="isSyncing" class="h-4 w-4 mr-2 animate-spin" />
          <RefreshCw v-else class="h-4 w-4 mr-2" />
          {{ $t('templates.syncFromMeta') }}
        </Button>
        <Button variant="outline" size="sm" @click="openCreateDialog">
          <Plus class="h-4 w-4 mr-2" />
          {{ $t('templates.createTemplate') }}
        </Button>
      </template>
    </PageHeader>

    <ScrollArea class="flex-1">
      <div class="p-6">
        <div class="max-w-6xl mx-auto">
          <Card>
            <CardHeader>
              <div class="flex items-center justify-between flex-wrap gap-4">
                <div>
                  <CardTitle>{{ $t('templates.yourTemplates') }}</CardTitle>
                  <CardDescription>{{ $t('templates.yourTemplatesDesc') }}</CardDescription>
                </div>
                <div class="flex items-center gap-4 flex-wrap">
                  <div class="flex items-center gap-2">
                    <Label class="text-sm text-muted-foreground">{{ $t('templates.account') }}:</Label>
                    <Select v-model="selectedAccount" @update:model-value="onAccountChange">
                      <SelectTrigger class="w-[180px]">
                        <SelectValue :placeholder="$t('templates.allAccounts')" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="all">{{ $t('templates.allAccounts') }}</SelectItem>
                        <SelectItem v-for="account in accounts" :key="account.id" :value="account.name">
                          {{ account.name }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <SearchInput v-model="searchQuery" :placeholder="$t('templates.searchTemplates') + '...'" class="w-64" />
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <DataTable
                :items="templates"
                :columns="columns"
                :is-loading="isLoading"
                :empty-icon="FileText"
                :empty-title="$t('templates.noTemplatesFound')"
                :empty-description="$t('templates.noTemplatesFoundDesc')"
                server-pagination
                :current-page="currentPage"
                :total-items="totalItems"
                :page-size="pageSize"
                item-name="templates"
                @page-change="handlePageChange"
                v-model:sort-key="sortKey"
                v-model:sort-direction="sortDirection"
              >
                <template #cell-name="{ item: template }">
                  <div>
                    <span class="font-medium">{{ template.display_name || template.name }}</span>
                    <p class="text-xs font-mono text-muted-foreground">{{ template.name }}</p>
                  </div>
                </template>
                <template #cell-category="{ item: template }">
                  <Badge :class="getCategoryBadgeClass(template.category)" class="text-xs">
                    {{ template.category }}
                  </Badge>
                </template>
                <template #cell-status="{ item: template }">
                  <Badge :class="getStatusBadgeClass(template.status)" class="text-xs">
                    {{ template.status }}
                  </Badge>
                </template>
                <template #cell-language="{ item: template }">
                  <span class="text-muted-foreground">{{ getLanguageName(template.language) }}</span>
                </template>
                <template #cell-header_type="{ item: template }">
                  <div class="flex items-center gap-1">
                    <component :is="getHeaderIcon(template.header_type)" class="h-4 w-4 text-muted-foreground" />
                    <span class="text-muted-foreground text-sm">{{ template.header_type || 'NONE' }}</span>
                  </div>
                </template>
                <template #cell-actions="{ item: template }">
                  <div class="flex items-center justify-end gap-1">
                    <Button variant="ghost" size="icon" class="h-8 w-8" @click="openPreview(template)">
                      <Eye class="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      class="h-8 w-8"
                      @click="openEditDialog(template)"
                      :disabled="template.status === 'PENDING'"
                    >
                      <Pencil class="h-4 w-4" />
                    </Button>
                    <Button
                      v-if="template.status === 'DRAFT' || template.status === 'REJECTED'"
                      variant="ghost"
                      size="icon"
                      class="h-8 w-8 text-blue-600 hover:text-blue-700"
                      @click="openPublishDialog(template)"
                      :disabled="publishingTemplateId === template.id"
                    >
                      <Loader2 v-if="publishingTemplateId === template.id" class="h-4 w-4 animate-spin" />
                      <Send v-else class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" class="h-8 w-8 text-destructive" @click="openDeleteDialog(template)">
                      <Trash2 class="h-4 w-4" />
                    </Button>
                  </div>
                </template>
                <template #empty-action>
                  <div class="flex items-center justify-center gap-2">
                    <Button variant="outline" size="sm" @click="syncTemplates" :disabled="!selectedAccount || selectedAccount === 'all'">
                      <RefreshCw class="h-4 w-4 mr-2" />
                      {{ $t('templates.syncFromMeta') }}
                    </Button>
                    <Button variant="outline" size="sm" @click="openCreateDialog">
                      <Plus class="h-4 w-4 mr-2" />
                      {{ $t('templates.createTemplate') }}
                    </Button>
                  </div>
                </template>
              </DataTable>
            </CardContent>
          </Card>
        </div>
      </div>
    </ScrollArea>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:open="isDialogOpen">
      <DialogContent class="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{{ editingTemplate ? $t('templates.editDialogTitle') : $t('templates.createDialogTitle') }}</DialogTitle>
          <DialogDescription>
            {{ editingTemplate ? $t('templates.editDialogDesc') : $t('templates.createDialogDesc') }}
          </DialogDescription>
        </DialogHeader>

        <div class="space-y-4 py-4">
          <!-- Account Selection -->
          <div class="space-y-2">
            <Label>{{ $t('templates.whatsappAccount') }} <span class="text-destructive">*</span></Label>
            <select
              v-model="formData.whatsapp_account"
              class="w-full h-10 rounded-md border bg-background px-3 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="!!editingTemplate"
            >
              <option value="">{{ $t('templates.selectAccount') }}...</option>
              <option v-for="account in accounts" :key="account.id" :value="account.name">
                {{ account.name }}
              </option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <!-- Template Name -->
            <div class="space-y-2">
              <Label>{{ $t('templates.templateName') }} <span class="text-destructive">*</span></Label>
              <Input
                v-model="formData.name"
                placeholder="order_confirmation"
                :disabled="!!editingTemplate"
              />
              <p class="text-xs text-muted-foreground">{{ $t('templates.templateNameLowercase') }}</p>
            </div>

            <!-- Display Name -->
            <div class="space-y-2">
              <Label>{{ $t('templates.displayName') }}</Label>
              <Input
                v-model="formData.display_name"
                placeholder="Order Confirmation"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <!-- Language -->
            <div class="space-y-2">
              <Label>{{ $t('templates.language') }} <span class="text-destructive">*</span></Label>
              <Popover v-model:open="languageSelectorOpen">
                <PopoverTrigger as-child>
                  <Button
                    variant="outline"
                    role="combobox"
                    class="w-full justify-between"
                    :disabled="!!editingTemplate"
                  >
                    <span>{{ getLanguageName(formData.language) }}</span>
                    <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent class="w-[300px] p-0">
                  <Command>
                    <CommandInput placeholder="Search language..." />
                    <CommandList>
                      <CommandEmpty>No language found.</CommandEmpty>
                      <CommandGroup>
                        <CommandItem
                          v-for="lang in languages"
                          :key="lang.code"
                          :value="lang.code"
                          @select="() => { formData.language = lang.code; languageSelectorOpen = false }"
                        >
                          <Check
                            :class="['mr-2 h-4 w-4', formData.language === lang.code ? 'opacity-100' : 'opacity-0']"
                          />
                          {{ lang.name }}
                        </CommandItem>
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
            </div>

            <!-- Category -->
            <div class="space-y-2">
              <Label>{{ $t('templates.category') }} <span class="text-destructive">*</span></Label>
              <Select v-model="formData.category">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="cat in categories" :key="cat.value" :value="cat.value">
                    <div>
                      <div class="font-medium">{{ cat.label }}</div>
                      <div class="text-xs text-muted-foreground">{{ cat.description }}</div>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <!-- Header Type -->
          <div class="space-y-2">
            <Label>{{ $t('templates.headerType') }}</Label>
            <Select v-model="formData.header_type">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="type in headerTypes" :key="type.value" :value="type.value">
                  {{ type.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- Header Content (Text) -->
          <div v-if="formData.header_type === 'TEXT'" class="space-y-2">
            <Label>{{ $t('templates.headerText') }}</Label>
            <Input v-model="formData.header_content" :placeholder="$t('templates.headerTextPlaceholder')" />
          </div>

          <!-- Header Media Upload (Image/Video/Document) -->
          <div v-if="formData.header_type !== 'NONE' && formData.header_type !== 'TEXT'" class="space-y-2">
            <Label>{{ $t('templates.headerMedia') }}</Label>
            <div class="flex items-center gap-2">
              <Input
                type="file"
                :accept="getAcceptedFileTypes()"
                @change="onHeaderMediaFileChange"
                class="flex-1"
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                @click="uploadHeaderMedia"
                :disabled="!headerMediaFile || headerMediaUploading || !formData.whatsapp_account"
              >
                <Loader2 v-if="headerMediaUploading" class="h-4 w-4 mr-2 animate-spin" />
                <Upload v-else class="h-4 w-4 mr-2" />
                {{ $t('templates.upload') }}
              </Button>
            </div>
            <p class="text-xs text-muted-foreground">
              {{ $t('templates.uploadMediaHint') }}
            </p>
            <div v-if="headerMediaHandle" class="flex items-center gap-2 text-sm text-green-600">
              <Check class="h-4 w-4" />
              <span>{{ $t('templates.mediaUploaded') }}: {{ headerMediaFilename || headerMediaHandle }}</span>
            </div>
          </div>

          <!-- Body Content -->
          <div class="space-y-2">
            <Label>{{ $t('templates.bodyContent') }} <span class="text-destructive">*</span></Label>
            <Textarea
              v-model="formData.body_content"
              :placeholder="$t('templates.bodyContentPlaceholder')"
              rows="4"
            />
            <p class="text-xs text-muted-foreground">
              {{ $t('templates.variablesHint') }}
            </p>
          </div>

          <!-- Footer Content -->
          <div class="space-y-2">
            <Label>{{ $t('templates.footerContent') }}</Label>
            <Input v-model="formData.footer_content" :placeholder="$t('templates.footerPlaceholder')" />
          </div>

          <!-- Sample Values for Variables -->
          <div v-if="bodyVariables.length > 0 || headerVariables.length > 0" class="space-y-2">
            <Label>{{ $t('templates.sampleValues') }}</Label>
            <p class="text-xs text-muted-foreground mb-2">
              {{ $t('templates.sampleValuesHint') }}
            </p>
            <div class="space-y-2">
              <div v-for="paramName in headerVariables" :key="'header-' + paramName" class="flex items-center gap-2">
                <Label class="w-32 text-sm">{{ formatVariableLabel(paramName) }} (Header)</Label>
                <Input
                  :value="getSampleValue('header', paramName)"
                  @input="(e: any) => setSampleValue('header', paramName, e.target.value)"
                  :placeholder="$t('templates.sampleValuePlaceholder')"
                  class="flex-1"
                />
              </div>
              <div v-for="paramName in bodyVariables" :key="'body-' + paramName" class="flex items-center gap-2">
                <Label class="w-32 text-sm">{{ formatVariableLabel(paramName) }}</Label>
                <Input
                  :value="getSampleValue('body', paramName)"
                  @input="(e: any) => setSampleValue('body', paramName, e.target.value)"
                  :placeholder="$t('templates.sampleValuePlaceholder')"
                  class="flex-1"
                />
              </div>
            </div>
          </div>

          <!-- Buttons -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <Label>{{ $t('templates.buttons') }}</Label>
              <Button type="button" variant="outline" size="sm" @click="addButton">
                <Plus class="h-4 w-4 mr-2" />
                {{ $t('templates.addButton') }}
              </Button>
            </div>
            <div v-if="formData.buttons.length > 0" class="space-y-2">
              <div
                v-for="(button, index) in formData.buttons"
                :key="index"
                class="flex items-center gap-2 p-3 border rounded-md"
              >
                <div class="flex-1 space-y-2">
                  <Select v-model="button.type">
                    <SelectTrigger class="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="type in buttonTypes" :key="type.value" :value="type.value">
                        <div>
                          <div class="font-medium">{{ type.label }}</div>
                          <div class="text-xs text-muted-foreground">{{ type.description }}</div>
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <Input v-model="button.text" :placeholder="$t('templates.buttonText')" />
                  <Input
                    v-if="button.type === 'URL'"
                    v-model="button.url"
                    :placeholder="$t('templates.buttonUrl')"
                  />
                  <Input
                    v-if="button.type === 'PHONE_NUMBER'"
                    v-model="button.phone_number"
                    :placeholder="$t('templates.buttonPhone')"
                  />
                </div>
                <Button type="button" variant="ghost" size="icon" @click="removeButton(index)">
                  <X class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="isDialogOpen = false">{{ $t('common.cancel') }}</Button>
          <Button @click="saveTemplate" :disabled="isSubmitting">
            <Loader2 v-if="isSubmitting" class="h-4 w-4 mr-2 animate-spin" />
            {{ editingTemplate ? $t('common.update') : $t('common.create') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Preview Dialog -->
    <Dialog v-model:open="isPreviewOpen">
      <DialogContent class="max-w-md">
        <DialogHeader>
          <DialogTitle>{{ $t('templates.preview') }}</DialogTitle>
          <DialogDescription>{{ previewTemplate?.name }}</DialogDescription>
        </DialogHeader>
        <div v-if="previewTemplate" class="space-y-4">
          <!-- Header Preview -->
          <div v-if="previewTemplate.header_type && previewTemplate.header_type !== 'NONE'" class="border-b pb-3">
            <div v-if="previewTemplate.header_type === 'TEXT'" class="font-semibold">
              <div v-html="formatPreview(previewTemplate.header_content, previewTemplate.sample_values.filter((s: any) => s.component === 'header'))"></div>
            </div>
            <div v-else class="flex items-center gap-2 text-muted-foreground">
              <component :is="getHeaderIcon(previewTemplate.header_type)" class="h-5 w-5" />
              <span>{{ previewTemplate.header_type }} {{ $t('templates.attachment') }}</span>
            </div>
          </div>
          <!-- Body Preview -->
          <div class="whitespace-pre-wrap" v-html="formatPreview(previewTemplate.body_content, previewTemplate.sample_values.filter((s: any) => s.component === 'body'))"></div>
          <!-- Footer Preview -->
          <div v-if="previewTemplate.footer_content" class="text-sm text-muted-foreground border-t pt-3">
            {{ previewTemplate.footer_content }}
          </div>
          <!-- Buttons Preview -->
          <div v-if="previewTemplate.buttons && previewTemplate.buttons.length > 0" class="space-y-2 border-t pt-3">
            <Button
              v-for="(button, index) in previewTemplate.buttons"
              :key="index"
              variant="outline"
              class="w-full"
              size="sm"
            >
              {{ button.text }}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ $t('templates.deleteConfirmTitle') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ $t('templates.deleteConfirmDesc', { name: templateToDelete?.name }) }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ $t('common.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmDeleteTemplate" class="bg-destructive text-destructive-foreground hover:bg-destructive/90">
            {{ $t('common.delete') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Publish Confirmation Dialog -->
    <AlertDialog v-model:open="publishDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ $t('templates.publishConfirmTitle') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ $t('templates.publishConfirmDesc', { name: templateToPublish?.name }) }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ $t('common.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmPublishTemplate">
            {{ $t('templates.publish') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
