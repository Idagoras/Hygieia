// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"encoding/json"
	"time"
)

type EegDatum struct {
	ID            uint64          `json:"id"`
	SessionID     uint64          `json:"session_id"`
	Offset        uint32          `json:"offset"`
	UserID        uint64          `json:"user_id"`
	CollectedAt   time.Time       `json:"collected_at"`
	Attention     uint32          `json:"attention"`
	Meditation    uint32          `json:"meditation"`
	BlinkStrength uint32          `json:"blink_strength"`
	Alpha1        uint32          `json:"alpha1"`
	Alpha2        uint32          `json:"alpha2"`
	Beta1         uint32          `json:"beta1"`
	Beta2         uint32          `json:"beta2"`
	Gamma1        uint32          `json:"gamma1"`
	Gamma2        uint32          `json:"gamma2"`
	Delta         uint32          `json:"delta"`
	Theta         uint32          `json:"theta"`
	Raw           json.RawMessage `json:"raw"`
}

type EegSession struct {
	ID          uint64       `json:"id"`
	Uid         uint64       `json:"uid"`
	Begin       time.Time    `json:"begin"`
	End         time.Time    `json:"end"`
	DataCount   uint32       `json:"data_count"`
	ExpiredAt   time.Time    `json:"expired_at"`
	Finish      sql.NullBool `json:"finish"`
	DataBitsMap string       `json:"data_bits_map"`
}

type Message struct {
	ID         uint64    `json:"id"`
	SendUid    uint64    `json:"send_uid"`
	RcvUid     uint64    `json:"rcv_uid"`
	CreatedAt  time.Time `json:"created_at"`
	HasRead    bool      `json:"has_read"`
	Type       uint32    `json:"type"`
	Text       string    `json:"text"`
	Subtype    uint32    `json:"subtype"`
	Attachment string    `json:"attachment"`
}

type Subscriber struct {
	ID            uint64    `json:"id"`
	Uid           uint64    `json:"uid"`
	SbsID         uint64    `json:"sbs_id"`
	SubscribeTime time.Time `json:"subscribe_time"`
	IsCancel      bool      `json:"is_cancel"`
	CancelTime    time.Time `json:"cancel_time"`
}

type User struct {
	ID                uint64    `json:"id"`
	Username          string    `json:"username"`
	Mobile            string    `json:"mobile"`
	Email             string    `json:"email"`
	Avatar            string    `json:"avatar"`
	HashedPassword    string    `json:"hashed_password"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedUsernameAt time.Time `json:"updated_username_at"`
	UpdatedAvatarAt   time.Time `json:"updated_avatar_at"`
	UpdatedMobileAt   time.Time `json:"updated_mobile_at"`
	UpdatedEmailAt    time.Time `json:"updated_email_at"`
}
