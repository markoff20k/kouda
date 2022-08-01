package models

import (
	"time"

	"github.com/volatiletech/null/v9"
)

type MemberState string

var (
	MemberStatePending = MemberState("pending")
	MemberStateActive  = MemberState("active")
	MemberStateSuspend = MemberState("suspend")
)

var MemberStates = []MemberState{
	MemberStateActive,
	MemberStatePending,
	MemberStateSuspend,
}

type Member struct {
	ID        int64       `gorm:"primaryKey;autoIncrement"`
	UID       string      `gorm:"type:character varying(32);not null;uniqueIndex:index_members_on_uid"`
	Email     string      `gorm:"type:character varying(255);not null;uniqueIndex:index_members_on_email"`
	Username  null.String `gorm:"type:character varying(255);uniqueIndex:index_members_on_username"`
	Level     int64       `gorm:"type:integer;not null"`
	Role      string      `gorm:"type:character varying(16);not null"`
	State     MemberState `gorm:"type:character varying(16);not null;default:pending"`
	CreatedAt time.Time   `gorm:"type:timestamp;not null"`
	UpdatedAt time.Time   `gorm:"type:timestamp;not null"`
}

func (m Member) TableName() string {
	return "members"
}

type MemberJSON struct {
	ID        int64       `json:"id,omitempty"`
	UID       string      `json:"uid,omitempty"`
	Email     string      `json:"email,omitempty"`
	Username  null.String `json:"username,omitempty"`
	Level     int64       `json:"level,omitempty"`
	Role      string      `json:"role,omitempty"`
	State     MemberState `json:"state,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
}

func (m *Member) ToJSON() MemberJSON {
	return MemberJSON{
		ID:        m.ID,
		UID:       m.UID,
		Email:     m.Email,
		Username:  m.Username,
		Level:     m.Level,
		Role:      m.Role,
		State:     m.State,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}
