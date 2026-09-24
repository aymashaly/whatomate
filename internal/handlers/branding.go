package handlers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/shridarpatil/whatomate/internal/models"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
	gormpkg "gorm.io/gorm"
)

// defaultPlatformBranding is returned when the platform_settings table is
// empty (fresh install, pre-migration, or the row was wiped). The frontend
// uses this to render a sensible UI before the super admin has saved anything.
func defaultPlatformBranding() models.PlatformSettings {
	return models.PlatformSettings{
		ID:           uuid.MustParse(models.SingletonPlatformSettingsID),
		BrandName:    "Whatomate",
		BrandTagline: "WhatsApp Platform",
		PrimaryColor: "142 71% 45%",
		AccentColor:  "173 80% 40%",
	}
}

// loadPlatformBranding returns the singleton branding row, creating it with
// defaults on first access. Every public read goes through this so callers
// (including the login page, which is unauthenticated) always get a
// complete payload.
func (a *App) loadPlatformBranding() (models.PlatformSettings, error) {
	var b models.PlatformSettings
	singletonID := models.SingletonPlatformSettingsID
	err := a.DB.Where("id = ?", singletonID).First(&b).Error
	if err == nil {
		return b, nil
	}
	if !errors.Is(err, gormpkg.ErrRecordNotFound) {
		return b, err
	}
	// No row yet — insert defaults so subsequent reads see a real record.
	b = defaultPlatformBranding()
	if err := a.DB.Create(&b).Error; err != nil {
		// Two racing inserts can collide on the PK; re-read on duplicate.
		if strings.Contains(err.Error(), "duplicate") {
			return a.loadPlatformBranding()
		}
		return b, err
	}
	return b, nil
}

// GetPlatformBranding returns the current platform-wide branding.
// Public — the login page needs this before any auth state exists.
func (a *App) GetPlatformBranding(r *fastglue.Request) error {
	b, err := a.loadPlatformBranding()
	if err != nil {
		a.Log.Error("Failed to load platform branding", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load branding", nil, "")
	}
	return r.SendEnvelope(b)
}

// UpdatePlatformBranding updates the platform-wide branding.
// Super-admin only — these fields belong to the SaaS provider's identity,
// not to any tenant. Other users get a clean 403.
func (a *App) UpdatePlatformBranding(r *fastglue.Request) error {
	userID, err := a.getUserIDFromCtx(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}
	if !a.IsSuperAdmin(userID) {
		return r.SendErrorEnvelope(fasthttp.StatusForbidden,
			"Only the platform administrator can change platform branding", nil, "")
	}

	var req struct {
		BrandName    *string `json:"brand_name"`
		BrandTagline *string `json:"brand_tagline"`
		LogoURL      *string `json:"logo_url"`
		FaviconURL   *string `json:"favicon_url"`
		PrimaryColor *string `json:"primary_color"`
		AccentColor  *string `json:"accent_color"`
		SupportEmail *string `json:"support_email"`
		SupportURL   *string `json:"support_url"`
		FooterText   *string `json:"footer_text"`
		LoginTagline *string `json:"login_tagline"`
	}
	if err := r.Decode(&req, "json"); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid request body", nil, "")
	}

	current, err := a.loadPlatformBranding()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load branding", nil, "")
	}

	if req.BrandName != nil {
		name := strings.TrimSpace(*req.BrandName)
		if name == "" {
			return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "brand_name cannot be empty", nil, "")
		}
		if len(name) > 100 {
			return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "brand_name is too long (max 100 chars)", nil, "")
		}
		current.BrandName = name
	}
	if req.BrandTagline != nil {
		current.BrandTagline = strings.TrimSpace(*req.BrandTagline)
	}
	if req.LogoURL != nil {
		current.LogoURL = strings.TrimSpace(*req.LogoURL)
	}
	if req.FaviconURL != nil {
		current.FaviconURL = strings.TrimSpace(*req.FaviconURL)
	}
	if req.PrimaryColor != nil {
		c := strings.TrimSpace(*req.PrimaryColor)
		if c != "" && !isValidHSLChannels(c) {
			return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
				"primary_color must be an HSL channel triplet like '142 71% 45%'", nil, "")
		}
		current.PrimaryColor = c
	}
	if req.AccentColor != nil {
		c := strings.TrimSpace(*req.AccentColor)
		if c != "" && !isValidHSLChannels(c) {
			return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
				"accent_color must be an HSL channel triplet like '173 80% 40%'", nil, "")
		}
		current.AccentColor = c
	}
	if req.SupportEmail != nil {
		current.SupportEmail = strings.TrimSpace(*req.SupportEmail)
	}
	if req.SupportURL != nil {
		current.SupportURL = strings.TrimSpace(*req.SupportURL)
	}
	if req.FooterText != nil {
		current.FooterText = strings.TrimSpace(*req.FooterText)
	}
	if req.LoginTagline != nil {
		current.LoginTagline = strings.TrimSpace(*req.LoginTagline)
	}

	if err := a.DB.Save(&current).Error; err != nil {
		a.Log.Error("Failed to update platform branding", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to update branding", nil, "")
	}

	a.Log.Info("Platform branding updated", "actor", userID, "brand", current.BrandName)
	return r.SendEnvelope(current)
}

// getUserIDFromCtx extracts the user UUID from request context.
func (a *App) getUserIDFromCtx(r *fastglue.Request) (uuid.UUID, error) {
	v := r.RequestCtx.UserValue("user_id")
	if v == nil {
		return uuid.Nil, errors.New("user_id not in context")
	}
	switch id := v.(type) {
	case uuid.UUID:
		return id, nil
	case string:
		return uuid.Parse(id)
	default:
		return uuid.Nil, errors.New("user_id has unexpected type")
	}
}

// isValidHSLChannels accepts the loose "H S% L%" channel triplet form used
// by CSS variables (no commas, no hsl() wrapper).
func isValidHSLChannels(s string) bool {
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return false
	}
	if !strings.HasSuffix(parts[1], "%") || !strings.HasSuffix(parts[2], "%") {
		return false
	}
	return true
}

// brandingAssetDir returns the on-disk directory for uploaded branding assets.
// Lives under the persisted `./uploads/` volume mount so logos and favicons
// survive container recreates (the audio dir is `./audio`, which is NOT
// mounted, so we deliberately ignore it here and hardcode the path under
// uploads).
func (a *App) brandingAssetDir() string {
	storagePath := a.Config.Storage.LocalPath
	if storagePath == "" {
		storagePath = "./uploads"
	}
	return filepath.Join(storagePath, "branding")
}

// UploadBrandingAsset handles multipart file uploads for branding logo and
// favicon. The ?type= query parameter must be "logo" or "favicon".
// Super-admin only.
func (a *App) UploadBrandingAsset(r *fastglue.Request) error {
	userID, err := a.getUserIDFromCtx(r)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Unauthorized", nil, "")
	}
	if !a.IsSuperAdmin(userID) {
		return r.SendErrorEnvelope(fasthttp.StatusForbidden,
			"Only the platform administrator can upload branding assets", nil, "")
	}

	assetType := string(r.RequestCtx.QueryArgs().Peek("type"))
	if assetType != "logo" && assetType != "favicon" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
			"Query parameter 'type' must be 'logo' or 'favicon'", nil, "")
	}

	form, err := r.RequestCtx.MultipartForm()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid multipart form: "+err.Error(), nil, "")
	}
	files := form.File["file"]
	if len(files) == 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "No file provided", nil, "")
	}
	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Failed to open file", nil, "")
	}
	defer func() { _ = file.Close() }()

	// 2 MB cap — branding assets are small icons / SVGs.
	const maxSize = 2 << 20
	data, err := io.ReadAll(io.LimitReader(file, maxSize+1))
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to read file", nil, "")
	}
	if len(data) > maxSize {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "File too large. Maximum size is 2MB", nil, "")
	}

	mime := strings.ToLower(fileHeader.Header.Get("Content-Type"))
	allowed := map[string]bool{
		"image/svg+xml":            true,
		"image/png":                true,
		"image/jpeg":               true,
		"image/jpg":                true,
		"image/webp":               true,
		"image/x-icon":             true,
		"image/vnd.microsoft.icon": true,
		"image/ico":                true,
		"application/octet-stream": true,
	}
	if !allowed[mime] {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest,
			"Unsupported image type: "+mime, nil, "")
	}

	ext := pickImageExtension(mime, fileHeader.Filename)
	if ext == "" {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Could not determine file extension", nil, "")
	}

	assetDir := a.brandingAssetDir()
	if err := os.MkdirAll(assetDir, 0755); err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to create asset directory", nil, "")
	}

	// Stable filename so the URL is predictable for caching.
	filename := fmt.Sprintf("platform_%s%s", assetType, ext)
	dest := filepath.Join(assetDir, filename)
	if err := os.WriteFile(dest, data, 0644); err != nil {
		a.Log.Error("Failed to write branding asset", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to save asset", nil, "")
	}

	url := "/api/platform/branding/asset/" + filename

	current, err := a.loadPlatformBranding()
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to load branding", nil, "")
	}
	if assetType == "logo" {
		current.LogoURL = url
	} else {
		current.FaviconURL = url
	}
	if err := a.DB.Save(&current).Error; err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to save asset URL", nil, "")
	}

	a.Log.Info("Branding asset uploaded", "actor", userID, "type", assetType, "filename", filename, "size", len(data))
	return r.SendEnvelope(map[string]any{
		"filename": filename,
		"type":     assetType,
		"url":      url,
		"size":     len(data),
	})
}

func pickImageExtension(mime, originalName string) string {
	switch mime {
	case "image/svg+xml":
		return ".svg"
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	case "image/x-icon", "image/vnd.microsoft.icon", "image/ico":
		return ".ico"
	}
	ext := strings.ToLower(filepath.Ext(originalName))
	switch ext {
	case ".svg", ".png", ".jpg", ".jpeg", ".webp", ".ico":
		return ext
	}
	return ""
}

// ServeBrandingAsset serves an uploaded branding file. Public (the favicon
// and logo are fetched by the browser without auth). Path traversal guards
// mirror the IVR audio serve pattern.
func (a *App) ServeBrandingAsset(r *fastglue.Request) error {
	filenameVal := r.RequestCtx.UserValue("filename")
	filename, _ := filenameVal.(string)
	filename = sanitizeFilename(filename)

	assetDir := a.brandingAssetDir()
	baseDir, err := filepath.Abs(assetDir)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Storage configuration error", nil, "")
	}
	fullPath, err := filepath.Abs(filepath.Join(baseDir, filename))
	if err != nil || !strings.HasPrefix(fullPath, baseDir+string(os.PathSeparator)) {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid file path", nil, "")
	}

	info, err := os.Lstat(fullPath)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusNotFound, "File not found", nil, "")
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return r.SendErrorEnvelope(fasthttp.StatusBadRequest, "Invalid file path", nil, "")
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "Failed to read file", nil, "")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	r.RequestCtx.Response.Header.Set("Content-Type", brandingMimeType(ext))
	// No-store: when the super admin replaces the logo, the new file uses
	// the same filename (platform_logo.svg / platform_favicon.*). Without
	// no-store, downstream caches (Cloudflare, browser) will keep serving
	// the previous bytes. The files are tiny so re-fetching is cheap.
	r.RequestCtx.Response.Header.Set("Cache-Control", "no-store")
	r.RequestCtx.SetBody(data)
	return nil
}

// brandingMimeType returns the right MIME for an uploaded branding asset.
// Inline the lookup (instead of reusing getMimeTypeFromExtension, which is
// for campaign media and doesn't know .svg) so the file serves with a
// real image type — important because the security middleware sets
// X-Content-Type-Options: nosniff, and browsers refuse to render an SVG
// announced as application/octet-stream.
func brandingMimeType(ext string) string {
	switch ext {
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}