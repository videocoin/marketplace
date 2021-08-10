package model

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"time"
)

type Account struct {
	ID                  int64 `db:"id"`
	IsActive            bool  `db:"is_active"`
	IsVerified          bool  `db:"is_verified"`
	CreatedAt           *time.Time
	Address             string
	Nonce               dbr.NullString
	Username            dbr.NullString
	Email               dbr.NullString
	Name                dbr.NullString
	PublicKey           dbr.NullString
	EncryptionPublicKey dbr.NullString `db:"enc_public_key"`
	ImageCID            dbr.NullString `db:"image_cid"`
	CoverCID            dbr.NullString `db:"cover_cid"`
	CustomURL           dbr.NullString `db:"custom_url"`
	Bio                 dbr.NullString `db:"bio"`
	YTUsername          dbr.NullString `db:"yt_username"`
}

func (u *Account) Id() int64 {
	return u.ID
}

func (u *Account) GetImageURL() *string {
	if u.ImageCID.String != "" {
		return pointer.ToString(fmt.Sprintf(DwebIpfsGateway, u.ImageCID.String))
	}

	return nil
}

func (u *Account) GetCoverURL() *string {
	if u.CoverCID.String != "" {
		return pointer.ToString(fmt.Sprintf(DwebIpfsGateway, u.CoverCID.String))
	}

	return nil
}
