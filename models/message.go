package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromUnique string `gorm:"index;not null"` // FK to Indentity.UniqueString
	ToUnique   string `gorm:"index;not null"` // FK to Indentity.UniqueString

	From Identity `gorm:"foreignKey:FromUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	To   Identity `gorm:"foreignKey:ToUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Text string `gorm:"type:text;not null"`
}
