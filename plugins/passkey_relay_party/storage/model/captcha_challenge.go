package model

import "time"

type ChallengeType string

const (
	Email ChallengeType = "email"
)

type CaptchaChallenge struct {
	ID        int64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time     `gorm:"column:created_at" json:"created_at"`
	Type      ChallengeType `gorm:"type:varchar(255);column:type"`
	Object    string        `gorm:"type:varchar(255);column:object;index:idx_object"`
	Code      string        `gorm:"type:varchar(16);column:code"`
}

func (CaptchaChallenge) TableName() string {
	return "challenge_captchas"
}
