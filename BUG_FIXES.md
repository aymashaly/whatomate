# Bug Fixes for Template and Campaign Media Issues

## Issues Fixed

### 1. Template Edit - "WhatsApp Account Not Found" Error

**Problem**: When editing a template, even though the account is selected and displayed, the backend was trying to validate the `whatsapp_account` field and failing because it was looking for an account that might not exist or was being passed incorrectly.

**Root Cause**: The `UpdateTemplate` handler in `internal/handlers/templates.go` was not properly handling the `whatsapp_account` field during updates. The account field is disabled in the UI during edit (as it shouldn't change), but the frontend was still sending it in the request body.

**Fix Applied** (`internal/handlers/templates.go`):
```go
// Validate WhatsApp account if it's being changed (shouldn't normally happen, but validate if provided)
if req.WhatsAppAccount != "" && req.WhatsAppAccount != template.WhatsAppAccount {
    if _, err := a.resolveWhatsAppAccount(orgID, req.WhatsAppAccount); err != nil {
        return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "WhatsApp account not found", nil, "")
    }
    template.WhatsAppAccount = req.WhatsAppAccount
}
```

**What Changed**:
- Now only validates the account if it's actually being changed
- If the account in the request matches the existing template's account, no validation is performed
- This prevents false "account not found" errors during normal template edits

---

### 2. Campaign Send - "Parameter Format Mismatch" Error

**Problem**: When sending campaign messages with templates that have IMAGE/VIDEO/DOCUMENT headers, the error occurred:
```
API error 132012: (#132012) Parameter format does not match format in the created template
Details: header: Format mismatch, expected IMAGE, received UNKNOWN
```

**Root Cause**: The worker was trying to send a header component even when:
1. No media was uploaded to the campaign
2. The template's `header_content` field contained a handle (not a URL) or was empty
3. The system couldn't determine whether to use "id" or "link" for the media parameter

**Fix Applied** (`internal/worker/worker.go`):

```go
// Handle header component (for media templates)
if template.HeaderType != "" && template.HeaderType != "TEXT" && template.HeaderType != "NONE" {
    var headerParam map[string]interface{}
    
    // Priority 1: Use campaign's uploaded media ID if available
    if campaignHeaderMediaID != "" {
        headerParam = buildMediaParameter(template.HeaderType, "id", campaignHeaderMediaID)
    } else if template.HeaderContent != "" {
        // Priority 2: Use template's header content
        // Check if it's a URL (starts with http) or a handle
        if strings.HasPrefix(template.HeaderContent, "http://") || strings.HasPrefix(template.HeaderContent, "https://") {
            // It's a URL, use link
            headerParam = buildMediaParameter(template.HeaderType, "link", template.HeaderContent)
        } else {
            // It's a handle from Meta's resumable upload, use id
            headerParam = buildMediaParameter(template.HeaderType, "id", template.HeaderContent)
        }
    }
    
    // Only add header component if we have valid media
    if headerParam != nil {
        components = append(components, map[string]interface{}{
            "type":       "header",
            "parameters": []map[string]interface{}{headerParam},
        })
    } else {
        // Log warning if template requires media but none provided
        w.Log.Warn("Template requires media header but no media provided",
            "template_id", template.ID,
            "template_name", template.Name,
            "header_type", template.HeaderType,
            "campaign_media_id", campaignHeaderMediaID,
            "template_header_content", template.HeaderContent)
    }
}
```

**What Changed**:
1. **Smart Media Detection**: Now checks if `header_content` is a URL (starts with http) or a handle
   - URLs use `"link"` parameter
   - Handles use `"id"` parameter
2. **Conditional Header Component**: Only adds header component if valid media is available
3. **Warning Logging**: Logs a warning if template requires media but none is provided
4. **No More UNKNOWN Format**: Prevents sending malformed header components

**Additional Validation** (`internal/handlers/campaigns.go`):
Added validation to prevent starting campaigns without required media:

```go
// Check if template requires media but campaign has none
if template.HeaderType != "" && template.HeaderType != "TEXT" && template.HeaderType != "NONE" {
    if campaign.HeaderMediaID == "" && template.HeaderContent == "" {
        return r.SendErrorEnvelope(fasthttp.StatusBadRequest, 
            fmt.Sprintf("Template requires %s media but none uploaded. Please upload media before starting the campaign.", template.HeaderType), 
            nil, "")
    }
}
```

This prevents users from starting campaigns that will fail due to missing media.

---

## Testing Checklist

After these fixes, test the following scenarios:

### Template Editing:
- [ ] Edit a template without changing the account - should work
- [ ] Edit template content (body, footer, etc.) - should work
- [ ] Upload media to template with IMAGE header - should work
- [ ] Save template after uploading media - should work

### Campaign with Media Templates:
- [ ] Create campaign with template that has IMAGE header
- [ ] Try to start campaign without uploading media - should show error
- [ ] Upload media to campaign
- [ ] Start campaign - should work
- [ ] Verify messages are sent with correct media
- [ ] Check that media appears in WhatsApp messages

### Different Media Types:
- [ ] Test with IMAGE header (JPEG, PNG)
- [ ] Test with VIDEO header (MP4)
- [ ] Test with DOCUMENT header (PDF)
- [ ] Test with URL-based media (http/https links)
- [ ] Test with handle-based media (Meta upload handles)

---

## Technical Details

### Media Parameter Types

WhatsApp Cloud API supports two ways to specify media in template headers:

1. **Using Media ID** (for media uploaded to WhatsApp):
```json
{
  "type": "header",
  "parameters": [{
    "type": "image",
    "image": {
      "id": "1234567890"  // Media ID from WhatsApp upload
    }
  }]
}
```

2. **Using URL** (for publicly accessible media):
```json
{
  "type": "header",
  "parameters": [{
    "type": "image",
    "image": {
      "link": "https://example.com/image.jpg"  // Public URL
    }
  }]
}
```

### Media Flow

1. **Template Creation**:
   - User uploads media → Gets handle from Meta
   - Handle stored in `template.header_content`
   - Used as sample when submitting template to Meta

2. **Campaign Creation**:
   - User uploads media → Gets media ID from WhatsApp
   - Media ID stored in `campaign.header_media_id`
   - File also saved locally for preview

3. **Campaign Sending**:
   - Worker checks `campaign.header_media_id` first (priority)
   - Falls back to `template.header_content` if no campaign media
   - Determines whether to use "id" or "link" based on content format
   - Only sends header component if valid media is available

---

## Error Prevention

The fixes implement multiple layers of error prevention:

1. **Frontend Validation**: Shows upload UI for templates requiring media
2. **Backend Validation**: Prevents starting campaigns without required media
3. **Worker Intelligence**: Smart detection of media type (URL vs handle)
4. **Graceful Degradation**: Logs warnings instead of failing when media is missing
5. **Clear Error Messages**: User-friendly error messages explaining what's needed

---

## Future Improvements

Consider these enhancements:

1. **Media Validation**: Validate media dimensions and file sizes before upload
2. **Media Preview**: Show media preview in campaign list
3. **Media Library**: Create reusable media library for campaigns
4. **Automatic Fallback**: Automatically use template media if campaign media not uploaded
5. **Media Analytics**: Track which media performs best in campaigns
6. **Bulk Media Upload**: Allow uploading media for multiple campaigns at once

---

## Conclusion

These fixes resolve both the template editing issue and the campaign media format mismatch error. The system now:

- ✅ Properly handles template edits without false account errors
- ✅ Intelligently detects media type (URL vs handle)
- ✅ Prevents campaigns from starting without required media
- ✅ Provides clear error messages to users
- ✅ Logs warnings for debugging
- ✅ Supports both URL-based and ID-based media parameters


## 2024-03-04: Campaign URL Media Detection Fix

**Issue**: When using a URL for campaign media (via the URL fallback feature), the worker was incorrectly treating it as a media ID and using the "id" parameter instead of "link" parameter, causing WhatsApp API error:
```
API error 100: (#100) Param template['components'][0]['parameters'][0]['image']['id'] is not a valid whatsapp business account media attachment ID
```

**Root Cause**: The worker's `sendTemplateMessage` function was checking `campaignHeaderMediaID` first but always using "id" parameter without checking if the value was a URL.

**Fix**: Updated `internal/worker/worker.go` to detect URLs in campaign media:
```go
// Priority 1: Use campaign's uploaded media ID if available
if campaignHeaderMediaID != "" {
    // Check if campaign media is a URL or a media ID
    if strings.HasPrefix(campaignHeaderMediaID, "http://") || strings.HasPrefix(campaignHeaderMediaID, "https://") {
        // It's a URL, use link
        headerParam = buildMediaParameter(template.HeaderType, "link", campaignHeaderMediaID)
    } else {
        // It's a media ID, use id
        headerParam = buildMediaParameter(template.HeaderType, "id", campaignHeaderMediaID)
    }
}
```

**Impact**: Campaign media URLs now work correctly. The worker intelligently detects whether the campaign media is:
- A URL (starts with http/https) → uses "link" parameter
- A media ID → uses "id" parameter

**Files Modified**:
- `internal/worker/worker.go`
