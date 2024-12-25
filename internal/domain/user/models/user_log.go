package models

import "time"

type UserLog struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;not null"`
	CreatedDate time.Time `json:"created_date,omitempty" gorm:"type:timestamptz;not null;default:current_timestamp"`
	Action      string    `json:"password,omitempty" gorm:"type:varchar(255);not null"`
	IP          string    `json:"ip,omitempty" gorm:"type:varchar(255);not null"`
	UserAgent   string    `json:"user_agent,omitempty" gorm:"type:varchar(255)"`
	UserID      string    `json:"user_id,omitempty" gorm:"type:varchar(50);unique"`
}
