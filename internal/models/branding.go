package models

import (
	"time"

	"github.com/google/uuid"
)

// PlatformSettings is a singleton row (enforced in code) that holds the
// platform-wide branding shown to every tenant — name, logo, favicon, accent
// color, support contact, footer text, login tagline. Owned exclusively by
// the platform super admin; all other users read it but cannot change it.
type PlatformSettings struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BrandName     string    `gorm:"size:100;not null;default:'Whatomate'" json:"brand_name"`
	BrandTagline  string    `gorm:"size:255;default:'WhatsApp Platform'" json:"brand_tagline"`
	LogoURL       string    `gorm:"size:500" json:"logo_url"`
	FaviconURL    string    `gorm:"size:500" json:"favicon_url"`
	PrimaryColor  string    `gorm:"size:32;default:'142 71% 45%'" json:"primary_color"`   // HSL channel triplet (no commas) — emerald default
	AccentColor   string    `gorm:"size:32;default:'173 80% 40%'" json:"accent_color"`    // secondary accent (cyan default)
	SupportEmail  string    `gorm:"size:255" json:"support_email"`
	SupportURL    string    `gorm:"size:500" json:"support_url"`
	FooterText    string    `gorm:"size:500" json:"footer_text"`
	LoginTagline  string    `gorm:"size:500" json:"login_tagline"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PlatformSettings) TableName() string {
	return "platform_settings"
}

// SingletonPlatformSettingsID is the fixed ID used for the single
// platform_settings row. Using a stable UUID lets the migration upsert
// without having to enumerate rows first.
const SingletonPlatformSettingsID = "00000000-0000-0000-0000-000000000001"