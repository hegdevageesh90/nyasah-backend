package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Review struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	Rating    int
	Content   string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SocialProof struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Type      string    // e.g., "purchase", "review", "view"
	ProductID uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Content   string
	CreatedAt time.Time
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

func (s *SocialProof) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}