package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"not null"`
	Domain    string    `gorm:"unique;not null"`
	Type      string    `gorm:"not null"` // e.g., "ecommerce", "education", "healthcare"
	ApiKey    string    `gorm:"unique;not null"`
	Active    bool      `gorm:"default:true"`
	Settings  JSON      `gorm:"type:json"`
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
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	EntityID  uuid.UUID `gorm:"type:uuid"`
	Rating    int
	Content   string
	Verified  bool
	Metadata  JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Tenant    Tenant `gorm:"foreignKey:TenantID"`
	User      User   `gorm:"foreignKey:UserID"`
	Entity    Entity `gorm:"foreignKey:EntityID"`
}

type SocialProof struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	Type      string    // e.g., "purchase", "review", "view", "enrollment"
	EntityID  uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Content   string
	Metadata  JSON `gorm:"type:json"`
	CreatedAt time.Time
	Tenant    Tenant `gorm:"foreignKey:TenantID"`
	User      User   `gorm:"foreignKey:UserID"`
	Entity    Entity `gorm:"foreignKey:EntityID"`
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
