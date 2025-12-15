package domain

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromUnique string   `gorm:"index;not null"`
	ToUnique   string   `gorm:"index;not null"`
	From       Identity `gorm:"foreignKey:FromUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	To         Identity `gorm:"foreignKey:ToUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Text       string   `gorm:"type:text;not null"`
}
