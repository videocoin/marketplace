package datastore

import (
	"context"
	"errors"
	"github.com/AlekSi/pointer"
	"github.com/gocraft/dbr/v2"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/pkg/dbrutil"
	"github.com/videocoin/marketplace/pkg/random"
	"time"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type UpdateAccountFields struct {
	Username *string
	Name     *string
	ImageURL *string
}

func (f *UpdateAccountFields) IsEmpty() bool {
	return f != nil && f.Username == nil && f.Name == nil && f.ImageURL == nil
}

type AccountDatastore struct {
	conn  *dbr.Connection
	table string
}

func NewAccountDatastore(ctx context.Context, conn *dbr.Connection) (*AccountDatastore, error) {
	return &AccountDatastore{
		conn:  conn,
		table: "accounts",
	}, nil
}

func (ds *AccountDatastore) Create(ctx context.Context, account *model.Account) error {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	if account.CreatedAt == nil || account.CreatedAt.IsZero() {
		account.CreatedAt = pointer.ToTime(time.Now())
	}

	if account.Nonce.String == "" {
		account.Nonce = dbr.NewNullString(random.RandomString(20))
	}

	cols := []string{"created_at", "address", "nonce"}
	err = tx.
		InsertInto(ds.table).
		Columns(cols...).
		Record(account).
		Returning("id").
		LoadContext(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (ds *AccountDatastore) List(ctx context.Context, fltr *AccountsFilter, limit *LimitOpts) ([]*model.Account, error) {
	var tx *dbr.Tx
	var err error

	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	accounts := make([]*model.Account, 0)
	selectStmt := tx.Select("*").From(ds.table)
	if fltr != nil {
		if fltr.Query != nil {
			likeQ := "%" + *fltr.Query + "%"
			selectStmt = selectStmt.Where(
				"name ILIKE ? OR username ILIKE ?",
				likeQ,
				likeQ,
			)
		}
		if fltr.Sort != nil && fltr.Sort.Field != "" {
			selectStmt = selectStmt.OrderDir(fltr.Sort.Field, fltr.Sort.IsAsc)
		}
	}

	if limit != nil {
		if limit.Offset != nil {
			selectStmt = selectStmt.Offset(*limit.Offset)
		}
		if limit.Limit != nil && *limit.Limit != 0 {
			selectStmt = selectStmt.Limit(*limit.Limit)
		}
	}

	_, err = selectStmt.LoadContext(ctx, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (ds *AccountDatastore) ListByIds(ctx context.Context, ids []int64) ([]*model.Account, error) {
	var tx *dbr.Tx
	var err error

	if len(ids) <= 0 {
		return []*model.Account{}, nil
	}

	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok || tx == nil {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	accounts := []*model.Account{}
	_, err = tx.
		Select("*").
		From(ds.table).
		Where("id IN ?", ids).
		LoadContext(ctx, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (ds *AccountDatastore) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	account := new(model.Account)
	err = tx.
		Select("*").
		From(ds.table).
		Where("id = ?", id).
		LoadOneContext(ctx, account)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (ds *AccountDatastore) GetByAddress(ctx context.Context, address string) (*model.Account, error) {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	account := new(model.Account)
	err = tx.
		Select("*").
		From(ds.table).
		Where("address = ?", address).
		LoadOneContext(ctx, account)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (ds *AccountDatastore) GetByUsername(ctx context.Context, username string) (*model.Account, error) {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return nil, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	account := new(model.Account)
	err = tx.
		Select("*").
		From(ds.table).
		Where("username = ?", username).
		LoadOneContext(ctx, account)
	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func (ds *AccountDatastore) RegenerateNonce(ctx context.Context, account *model.Account) error {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	nonce := dbr.NewNullString(random.RandomString(20))

	_, err = tx.
		Update(ds.table).
		Set("nonce", nonce).
		Where("address = ?", account.Address).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	account.Nonce = dbr.NewNullString(nonce)

	return nil
}

func (ds *AccountDatastore) UpdatePublicKey(ctx context.Context, account *model.Account, pk string) error {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	_, err = tx.
		Update(ds.table).
		Set("public_key", pk).
		Where("address = ?", account.Address).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	account.PublicKey = dbr.NewNullString(pk)

	return nil
}

func (ds *AccountDatastore) Update(ctx context.Context, account *model.Account, fields UpdateAccountFields) error {
	var err error
	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	stmt := tx.Update(ds.table)
	if fields.Username != nil {
		stmt.Set("username", dbr.NewNullString(*fields.Username))
		account.Username = dbr.NewNullString(*fields.Username)
	}

	if fields.Name != nil {
		stmt.Set("name", dbr.NewNullString(*fields.Name))
		account.Name = dbr.NewNullString(*fields.Name)
	}

	if fields.ImageURL != nil {
		stmt.Set("image_url", dbr.NewNullString(*fields.ImageURL))
		account.ImageURL = dbr.NewNullString(*fields.ImageURL)
	}

	_, err = stmt.Where("id = ?", account.ID).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ds *AccountDatastore) Count(ctx context.Context, fltr *AccountsFilter) (int64, error) {
	var tx *dbr.Tx
	var err error

	tx, ok := dbrutil.DbTxFromContext(ctx)
	if !ok {
		sess := ds.conn.NewSession(nil)
		tx, err = sess.Begin()
		if err != nil {
			return 0, err
		}

		defer func() {
			err = tx.Commit()
			tx.RollbackUnlessCommitted()
		}()
	}

	count := int64(0)
	selectStmt := tx.Select("COUNT(id)").From(ds.table)
	if fltr.Query != nil {
		likeQ := "%" + *fltr.Query + "%"
		selectStmt = selectStmt.Where(
			"name ILIKE ? OR username ILIKE ?",
			likeQ,
			likeQ,
		)
	}

	err = selectStmt.LoadOneContext(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
