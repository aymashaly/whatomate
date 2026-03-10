# Campaign Media URL Fallback Feature

## Overview
Added the ability to use a public URL as a fallback for campaign media instead of uploading files. This is useful when:
- File upload is not working in development/production
- You want to use media hosted on a CDN
- You need to quickly test campaigns without uploading files
- You have media already hosted elsewhere

## Changes Made

### 1. Frontend - Campaign Media Tabs (`frontend/src/views/settings/CampaignsView.vue`)

#### Added State Variables:
```typescript
const campaignMediaUrl = ref('')
const campaignMediaMode = ref<'upload' | 'url'>('upload')
const savingMediaUrl = ref(false)
```

#### Added Function:
```typescript
async function saveCampaignMediaUrl() {
  // Validates URL format
  // Checks if URL is publicly accessible (http/https)
  // Updates campaign with URL
  // Stores URL in header_media_id field
}
```

#### UI Changes:
- **Tabs Interface**: Added tabs to switch between "Upload File" and "Use URL"
- **Upload Tab**: Original file upload functionality
- **URL Tab**: 
  - Input field for public URL
  - Validation hints
  - Requirements info box
  - Save button

#### URL Requirements Display:
```vue
<div class="bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-800 rounded p-3">
  <div class="flex gap-2">
    <AlertCircle class="h-4 w-4 text-blue-600" />
    <div class="text-xs text-blue-800">
      <p class="font-medium">URL Requirements:</p>
      <ul class="list-disc list-inside">
        <li>Must use HTTPS (or HTTP)</li>
        <li>Must be publicly accessible</li>
        <li>Recommended: Use CDN for better performance</li>
      </ul>
    </div>
  </div>
</div>
```

### 2. Backend - Campaign Update Handler (`internal/handlers/campaigns.go`)

#### Updated CampaignRequest Struct:
```go
type CampaignRequest struct {
    Name            string     `json:"name" validate:"required"`
    WhatsAppAccount string     `json:"whatsapp_account" validate:"required"`
    TemplateID      string     `json:"template_id" validate:"required"`
    HeaderMediaID   string     `json:"header_media_id"`
    HeaderMediaURL  string     `json:"header_media_url"` // NEW: Public URL for media
    ScheduledAt     *time.Time `json:"scheduled_at"`
}
```

#### Updated UpdateCampaign Handler:
```go
// Handle media URL (alternative to file upload)
if req.HeaderMediaURL != "" {
    // Validate URL format
    if !strings.HasPrefix(req.HeaderMediaURL, "http://") && !strings.HasPrefix(req.HeaderMediaURL, "https://") {
        return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Media URL must start with http:// or https://", nil, "")
    }
    // Store URL in header_media_id field (worker will detect it's a URL)
    updates["header_media_id"] = req.HeaderMediaURL
    updates["header_media_filename"] = "External URL"
    updates["header_media_mime_type"] = "" // Unknown for external URLs
    a.Log.Info("Campaign media URL set", "campaign_id", id, "url", req.HeaderMediaURL)
}
```

### 3. Worker - Already Supports URLs (`internal/worker/worker.go`)

The worker already intelligently handles URLs vs media IDs:

```go
if strings.HasPrefix(template.HeaderContent, "http://") || strings.HasPrefix(template.HeaderContent, "https://") {
    // It's a URL, use link
    headerParam = buildMediaParameter(template.HeaderType, "link", template.HeaderContent)
} else {
    // It's a handle from Meta's resumable upload, use id
    headerParam = buildMediaParameter(template.HeaderType, "id", template.HeaderContent)
}
```

This same logic applies to `campaignHeaderMediaID`, so URLs work automatically!

## How It Works

### User Flow:

1. **Create Campaign** with a template that has IMAGE/VIDEO/DOCUMENT header
2. **Open "View Recipients"** dialog
3. **See Media Upload Section** with two tabs:
   - **Upload File**: Traditional file upload
   - **Use URL**: Enter public URL

4. **Choose URL Tab**:
   - Enter public URL (e.g., `https://cdn.example.com/image.jpg`)
   - Click "Save URL"
   - System validates URL format
   - URL is stored in campaign

5. **Start Campaign**:
   - Worker detects URL (starts with http)
   - Uses "link" parameter in WhatsApp API
   - Messages sent with media from URL

### API Flow:

```
Frontend                    Backend                     WhatsApp API
   |                           |                              |
   |-- POST /campaigns/:id --> |                              |
   |   { header_media_url:     |                              |
   |     "https://..." }       |                              |
   |                           |                              |
   |                           |-- Validate URL               |
   |                           |-- Store in header_media_id   |
   |                           |                              |
   |<-- 200 OK ----------------|                              |
   |                           |                              |
   |-- POST /campaigns/:id/    |                              |
   |    start                  |                              |
   |                           |                              |
   |                           |-- Worker processes           |
   |                           |-- Detects URL                |
   |                           |-- Builds template with       |
   |                           |   "link" parameter           |
   |                           |                              |
   |                           |-- POST /messages ----------->|
   |                           |   {                          |
   |                           |     "template": {            |
   |                           |       "components": [{       |
   |                           |         "type": "header",    |
   |                           |         "parameters": [{     |
   |                           |           "type": "image",   |
   |                           |           "image": {         |
   |                           |             "link": "https://..."|
   |                           |           }                  |
   |                           |         }]                   |
   |                           |       }]                     |
   |                           |     }                        |
   |                           |   }                          |
   |                           |                              |
   |                           |<-- Message sent -------------|
```

## Use Cases

### 1. Development Environment Issues
```
Problem: File upload not working in dev
Solution: Use URL from CDN or public hosting
Example: https://i.imgur.com/example.jpg
```

### 2. Production Fallback
```
Problem: Upload service temporarily down
Solution: Quickly switch to URL mode
Example: https://cdn.yoursite.com/campaign-image.jpg
```

### 3. CDN-Hosted Media
```
Benefit: Better performance, global distribution
Solution: Upload to CDN first, then use URL
Example: https://cloudinary.com/your-account/image.jpg
```

### 4. Testing
```
Benefit: Quick testing without uploading
Solution: Use any public image URL
Example: https://picsum.photos/800/600
```

## URL Requirements

### Must Have:
- ✅ Start with `http://` or `https://` (HTTPS recommended)
- ✅ Be publicly accessible (no authentication required)
- ✅ Point directly to the media file
- ✅ Return correct Content-Type header

### Recommended:
- 🌟 Use HTTPS for security
- 🌟 Use CDN for better performance
- 🌟 Use appropriate file formats:
  - Images: JPEG, PNG, WebP
  - Videos: MP4, 3GP
  - Documents: PDF
- 🌟 Optimize file sizes (< 5MB for images, < 16MB for videos)

### Examples of Valid URLs:
```
✅ https://cdn.example.com/images/campaign-banner.jpg
✅ https://i.imgur.com/abc123.png
✅ https://cloudinary.com/demo/image/upload/sample.jpg
✅ https://storage.googleapis.com/bucket/image.jpg
✅ https://s3.amazonaws.com/bucket/media/video.mp4
```

### Examples of Invalid URLs:
```
❌ example.com/image.jpg (missing protocol)
❌ ftp://example.com/image.jpg (wrong protocol)
❌ https://example.com/login-required (requires auth)
❌ https://example.com/page.html (not direct media link)
```

## Validation

### Frontend Validation:
```typescript
// Check URL format
try {
  new URL(campaignMediaUrl.value)
} catch {
  toast.error('Invalid URL format')
  return
}

// Check protocol
if (!campaignMediaUrl.value.startsWith('http://') && 
    !campaignMediaUrl.value.startsWith('https://')) {
  toast.error('URL must be publicly accessible (http/https)')
  return
}
```

### Backend Validation:
```go
// Validate URL format
if !strings.HasPrefix(req.HeaderMediaURL, "http://") && 
   !strings.HasPrefix(req.HeaderMediaURL, "https://") {
    return r.SendErrorEnvelope(fasthttp.StatusBadRequest, 
        "Media URL must start with http:// or https://", nil, "")
}
```

## Testing Checklist

- [ ] Switch between Upload and URL tabs
- [ ] Enter valid HTTPS URL and save
- [ ] Enter invalid URL (no protocol) - should show error
- [ ] Enter URL with HTTP (not HTTPS) - should work but warn
- [ ] Save URL and verify it's stored in campaign
- [ ] Start campaign with URL-based media
- [ ] Verify message is sent with correct media
- [ ] Check media appears in WhatsApp
- [ ] Test with different media types (image, video, document)
- [ ] Test with CDN URLs (Cloudinary, Imgur, etc.)

## Error Handling

### Frontend Errors:
- Empty URL → "Please enter a media URL"
- Invalid format → "Invalid URL format"
- Missing protocol → "URL must be publicly accessible (http/https)"

### Backend Errors:
- Invalid protocol → "Media URL must start with http:// or https://"
- Campaign not draft → "Can only update draft campaigns"

### WhatsApp API Errors:
- URL not accessible → "Failed to fetch media from URL"
- Invalid media type → "Media type not supported"
- File too large → "Media file exceeds size limit"

## Benefits

1. **Flexibility**: Choose between upload or URL based on situation
2. **Reliability**: Fallback option when upload fails
3. **Performance**: Use CDN for better global delivery
4. **Speed**: Quick testing without uploading
5. **Cost**: Leverage existing media hosting
6. **Simplicity**: Just paste a URL and go

## Future Enhancements

1. **URL Preview**: Show thumbnail of URL before saving
2. **URL Validation**: Check if URL is accessible before saving
3. **Media Type Detection**: Auto-detect media type from URL
4. **URL History**: Remember recently used URLs
5. **Bulk URL Import**: Import multiple URLs from CSV
6. **URL Templates**: Save frequently used URL patterns

## Conclusion

The URL fallback feature provides a robust alternative to file uploads, ensuring campaigns can always be launched even when upload functionality is unavailable. The system intelligently handles both uploaded media IDs and public URLs, making it flexible and reliable for production use.
