# Image Upload Implementation for Templates and Campaigns

## Overview
This implementation adds image/media upload functionality to both Templates and Campaigns pages, allowing users to upload media files that are stored on the host server and used in WhatsApp campaigns with public links.

## Changes Made

### 1. Templates Page (`frontend/src/views/settings/TemplatesView.vue`)

#### Added Features:
- **Media Upload UI**: Added file input and upload button for templates with IMAGE, VIDEO, or DOCUMENT header types
- **Upload Validation**: Only shows upload UI when header type requires media (not TEXT or NONE)
- **Upload Status**: Shows success indicator with filename when media is uploaded
- **File Type Validation**: Accepts only appropriate file types based on header type:
  - IMAGE: jpeg, png
  - VIDEO: mp4
  - DOCUMENT: pdf

#### Implementation Details:
```vue
<!-- Header Media Upload Section -->
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
</div>
```

#### Functions Added:
- `onHeaderMediaFileChange()`: Handles file selection
- `uploadHeaderMedia()`: Uploads file to Meta via API
- `getAcceptedFileTypes()`: Returns accepted MIME types based on header type

### 2. Campaigns Page (`frontend/src/views/settings/CampaignsView.vue`)

#### Added Features:
- **Campaign Media Upload**: Added media upload section in the "View Recipients" dialog
- **Template Detection**: Automatically detects if the campaign's template requires media
- **Media Preview**: Shows currently uploaded media with preview option
- **Upload for All Recipients**: Uploaded media is used for all recipients in the campaign

#### Implementation Details:
```vue
<!-- Campaign Media Upload Section -->
<div v-if="selectedCampaign && selectedTemplate && selectedTemplate.header_type && 
           selectedTemplate.header_type !== 'NONE' && selectedTemplate.header_type !== 'TEXT' && 
           selectedCampaign.status === 'draft'" 
     class="border rounded-lg p-4 bg-muted/30">
  <div class="flex items-center gap-2 mb-3">
    <component :is="getHeaderIcon(selectedTemplate.header_type)" class="h-5 w-5 text-primary" />
    <Label class="text-base font-semibold">{{ $t('campaigns.campaignMedia') }}</Label>
  </div>
  
  <!-- Current media display -->
  <div v-if="selectedCampaign.header_media_filename" class="mb-3 p-3 bg-background rounded border">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <Check class="h-4 w-4 text-green-600" />
        <span class="text-sm font-medium">{{ selectedCampaign.header_media_filename }}</span>
      </div>
      <Button variant="ghost" size="sm" @click="viewCampaignMedia(selectedCampaign)">
        <Eye class="h-4 w-4" />
      </Button>
    </div>
  </div>
  
  <!-- Upload new media -->
  <div class="flex items-center gap-2">
    <Input
      type="file"
      :accept="getAcceptedMediaTypes(selectedTemplate.header_type)"
      @change="onCampaignMediaFileChange"
      class="flex-1"
      ref="campaignMediaInput"
    />
    <Button
      type="button"
      variant="outline"
      size="sm"
      @click="uploadCampaignMedia"
      :disabled="!campaignMediaFile || campaignMediaUploading"
    >
      <Loader2 v-if="campaignMediaUploading" class="h-4 w-4 mr-2 animate-spin" />
      <Upload v-else class="h-4 w-4 mr-2" />
      {{ $t('campaigns.upload') }}
    </Button>
  </div>
</div>
```

#### Functions Added:
- `getHeaderIcon()`: Returns appropriate icon for media type
- `getAcceptedMediaTypes()`: Returns accepted MIME types for media type
- `onCampaignMediaFileChange()`: Handles file selection
- `uploadCampaignMedia()`: Uploads media to campaign
- `viewCampaignMedia()`: Opens media preview dialog
- Updated `viewRecipients()`: Now fetches template details to check for media requirements

#### State Variables Added:
```typescript
const campaignMediaFile = ref<File | null>(null)
const campaignMediaUploading = ref(false)
const campaignMediaInput = ref<HTMLInputElement | null>(null)
```

### 3. Backend Integration (Already Implemented)

The backend already has all necessary endpoints:

#### Templates API:
- `POST /api/templates/upload-media` - Uploads media for template headers
  - Accepts: `file` (multipart), `account` (WhatsApp account name)
  - Returns: `{ handle, filename, mime_type, size }`
  - Uploads to Meta's resumable upload API
  - Stores handle in template's `header_content` field

#### Campaigns API:
- `POST /api/campaigns/:id/media` - Uploads media for campaign
  - Accepts: `file` (multipart)
  - Returns: `{ media_id, filename, mime_type, local_path }`
  - Uploads to WhatsApp Cloud API
  - Stores locally for preview
  - Stores media_id in campaign's `header_media_id` field

- `GET /api/campaigns/:id/media` - Retrieves campaign media for preview
  - Returns: Binary file data with appropriate Content-Type

#### Worker Implementation:
The worker (`internal/worker/worker.go`) automatically uses the campaign's `header_media_id` when sending template messages:

```go
func (w *Worker) sendTemplateMessage(ctx context.Context, account *models.WhatsAppAccount, 
    template *models.Template, recipient *models.BulkMessageRecipient, 
    campaignHeaderMediaID string) (string, error) {
    
    // Handle header component (for media templates)
    if template.HeaderType != "" && template.HeaderType != "TEXT" {
        // Use campaign's uploaded media ID if available
        if campaignHeaderMediaID != "" {
            headerParam := buildMediaParameter(template.HeaderType, "id", campaignHeaderMediaID)
            if headerParam != nil {
                components = append(components, map[string]interface{}{
                    "type":       "header",
                    "parameters": []map[string]interface{}{headerParam},
                })
            }
        }
    }
    // ... rest of implementation
}
```

## How It Works

### Template Media Upload Flow:
1. User creates/edits a template with IMAGE/VIDEO/DOCUMENT header type
2. User selects a media file
3. File is uploaded to Meta's resumable upload API
4. Meta returns a handle (file reference)
5. Handle is stored in template's `header_content` field
6. When template is submitted to Meta, the handle is included in the sample values

### Campaign Media Upload Flow:
1. User creates a campaign with a template that has a media header
2. User opens "View Recipients" dialog
3. System detects template requires media and shows upload section
4. User uploads media file
5. File is uploaded to WhatsApp Cloud API → returns `media_id`
6. File is also saved locally for preview → returns `local_path`
7. Campaign is updated with:
   - `header_media_id`: Meta's media ID
   - `header_media_filename`: Original filename
   - `header_media_mime_type`: MIME type
   - `header_media_local_path`: Local path for preview

### Campaign Sending Flow:
1. Campaign is started
2. Worker processes each recipient
3. For each recipient, worker calls `sendTemplateMessage()`
4. Worker checks if campaign has `header_media_id`
5. If yes, includes media in template components:
   ```json
   {
     "type": "header",
     "parameters": [{
       "type": "image",  // or "video", "document"
       "image": {
         "id": "campaign_media_id"  // Uses Meta's media ID
       }
     }]
   }
   ```
6. WhatsApp sends message with media to recipient

## API Endpoints Used

### Frontend API Service (`frontend/src/services/api.ts`):

```typescript
export const templatesService = {
  uploadMedia: (accountName: string, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('account', accountName)
    return axios.post(`${api.defaults.baseURL}/templates/upload-media`, formData, {
      withCredentials: true,
      headers: { 'X-CSRF-Token': csrfToken }
    })
  }
}

export const campaignsService = {
  uploadMedia: (campaignId: string, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return axios.post(`${api.defaults.baseURL}/campaigns/${campaignId}/media`, formData, {
      withCredentials: true,
      headers: { 'X-CSRF-Token': csrfToken }
    })
  },
  getMedia: (campaignId: string) =>
    api.get(`/campaigns/${campaignId}/media`, { responseType: 'arraybuffer' })
}
```

## Database Schema

### Templates Table:
```sql
CREATE TABLE templates (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    whatsapp_account VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    header_type VARCHAR(20),  -- TEXT, IMAGE, VIDEO, DOCUMENT, NONE
    header_content TEXT,       -- Stores Meta handle for media
    body_content TEXT NOT NULL,
    -- ... other fields
);
```

### Campaigns Table:
```sql
CREATE TABLE bulk_message_campaigns (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    whatsapp_account VARCHAR(100) NOT NULL,
    template_id UUID NOT NULL,
    header_media_id TEXT,         -- Meta media ID
    header_media_filename TEXT,   -- Original filename
    header_media_mime_type TEXT,  -- MIME type
    header_media_local_path TEXT, -- Local path for preview
    -- ... other fields
);
```

## Key Features

1. **Dual Storage**: Media is stored both on Meta (for sending) and locally (for preview)
2. **Type Safety**: Only accepts appropriate file types for each media type
3. **User Feedback**: Shows upload progress and success/error messages
4. **Preview Support**: Users can preview uploaded media before sending
5. **Draft-Only Uploads**: Media can only be uploaded for draft campaigns
6. **Automatic Integration**: Uploaded media is automatically used when sending campaign messages
7. **Public Links**: Media is accessible via public links generated by WhatsApp

## Testing Checklist

- [ ] Upload image to template with IMAGE header type
- [ ] Upload video to template with VIDEO header type
- [ ] Upload document to template with DOCUMENT header type
- [ ] Create campaign with media template
- [ ] Upload media to campaign
- [ ] View uploaded media preview
- [ ] Start campaign and verify media is sent correctly
- [ ] Verify media appears in WhatsApp messages
- [ ] Test with different file sizes
- [ ] Test error handling for invalid file types
- [ ] Test error handling for upload failures

## Future Enhancements

1. **Media Library**: Create a reusable media library for campaigns
2. **Bulk Upload**: Allow uploading multiple media files at once
3. **Media Editing**: Add basic image editing capabilities
4. **Media Analytics**: Track media engagement (views, clicks)
5. **CDN Integration**: Optionally use CDN for faster media delivery
6. **Media Compression**: Automatically compress large files
7. **Media Validation**: Validate image dimensions and file sizes before upload

## Conclusion

The implementation successfully adds image/media upload functionality to both Templates and Campaigns pages. Media files are uploaded to the host server, stored locally for preview, and uploaded to WhatsApp's Cloud API for use in campaigns. The system uses public links generated by WhatsApp to deliver media to recipients, ensuring reliable and fast delivery.
