package model

import (
	"github.com/gocraft/dbr/v2"
	"time"
)

type Account struct {
	ID         int64 `db:"id"`
	IsActive   bool  `db:"is_active"`
	IsVerified bool  `db:"is_verified"`
	CreatedAt  *time.Time
	Address    string
	Nonce      dbr.NullString
	Username   dbr.NullString
	Email      dbr.NullString
	Name       dbr.NullString
	PublicKey  dbr.NullString
	ImageURL   dbr.NullString `db:"image_url"`
	CoverURL   dbr.NullString `db:"cover_url"`
	CustomURL  dbr.NullString `db:"custom_url"`
	Bio        dbr.NullString `db:"bio"`
	YTUsername dbr.NullString `db:"yt_username"`
}

func (u *Account) Id() int64 {
	return u.ID
}
