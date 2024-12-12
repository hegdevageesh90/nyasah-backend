package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"encoding/json"
)

type Tenant struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key"`
	Name      string          `gorm:"not null"`
	Domain    string          `gorm:"unique;not null"`
	Type      string          `gorm:"not null"` // e.g., "ecommerce", "education", "healthcare"
	ApiKey    string          `gorm:"unique;not null"`
	Active    bool            `gorm:"default:true"`
	Settings  json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Name      string
	Role      string `gorm:"default:'user'"` // 'admin', 'user'
	CreatedAt time.Time
	UpdatedAt time.Time
	Tenant    Tenant `gorm:"foreignKey:TenantID"`
}

type Entity struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null"`
	Type        string    `gorm:"not null"` // e.g., "product", "course", "property"
	Name        string    `gorm:"not null"`
	Description string
	Metadata    JSON `gorm:"type:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tenant      Tenant `gorm:"foreignKey:TenantID"`
}

type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID   string    `gorm:"type:uuid;not null"`
	UserID     string
	EntityID   uuid.UUID `gorm:"type:uuid"`
	Rating     int
	Content    string
	Verified   bool
	Metadata   json.RawMessage `gorm:"type:json"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Tenant     Tenant           `gorm:"foreignKey:TenantID"`
	User       User             `gorm:"foreignKey:UserID"`
	Entity     Entity           `gorm:"foreignKey:EntityID"`
	Engagement ReviewEngagement `gorm:"foreignKey:ReviewID"`
	Sentiment  float64          // AI-analyzed sentiment score
	Keywords   []string         `gorm:"type:json"`
}

type ReviewEngagement struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	ReviewID  uuid.UUID `gorm:"type:uuid"`
	Views     int
	Likes     int
	Shares    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AIInsights  ProductInsights `gorm:"foreignKey:ProductID"`
}

type ProductInsights struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key"`
	ProductID          uuid.UUID `gorm:"type:uuid"`
	SentimentTrend     []float64 `gorm:"type:json"`
	TopKeywords        []string  `gorm:"type:json"`
	EngagementScore    float64
	RecommendedActions []string `gorm:"type:json"`
	AverageRating      float64
	EngagementRate     float64
	SentimentScore     float64
	LastUpdated        time.Time
}

type ProofPerformance struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	ProofID        uuid.UUID `gorm:"type:uuid"`
	Views          int
	Conversions    int
	EngagementRate float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AIQuery struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID  uuid.UUID `gorm:"type:uuid"`
	Query     string
	Response  string
	CreatedAt time.Time
}

type AIRecommendation struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID   uuid.UUID `gorm:"type:uuid"`
	Type       string    // "content", "timing", "placement"
	Suggestion string
	Confidence float64
	CreatedAt  time.Time
}

type SocialProof struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID    uuid.UUID `gorm:"type:uuid;not null"`
	Type        string    // e.g., "purchase", "review", "view", "enrollment"
	EntityID    uuid.UUID `gorm:"type:uuid"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	Content     string
	Metadata    JSON `gorm:"type:json"`
	CreatedAt   time.Time
	Tenant      Tenant           `gorm:"foreignKey:TenantID"`
	User        User             `gorm:"foreignKey:UserID"`
	Entity      Entity           `gorm:"foreignKey:EntityID"`
	Performance ProofPerformance `gorm:"foreignKey:ProofID"`
	MediaType   string           // "image", "video", "text"
}

// JSON is a custom type for handling JSON data
type JSON map[string]interface{}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

func (e *Entity) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

func (s *SocialProof) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}
