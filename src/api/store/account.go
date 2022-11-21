package store

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"database/sql"
)

type AccountStore struct {
	DB *sql.DB
}

func NewAccountStore(db *sql.DB) *AccountStore {
	return &AccountStore{DB: db}
}

const sqlFindAccountByUserName = `
SELECT
	accounts.username,
	accounts.hashed_password,
	accounts.google_id,
	accounts.created_at,
	accounts.updated_at
FROM accounts
WHERE
	accounts.username = ?
LIMIT 1;`

const sqlCreateAccount = `
INSERT INTO accounts (
	username,
	hashed_password
) VALUES (
	?,
  	?
);`

const sqlCreateAccountGoogle = `
INSERT INTO accounts (
	username,
	google_id
) VALUES (
	?,
  	?
);`

const sqlUpdatePasswordToAccount = `
UPDATE accounts 
SET hashed_password = ? 
WHERE username = ?;
`

const sqlUpdateGoogleIDToAccount = `
UPDATE accounts 
SET google_id = ? 
WHERE username = ?;
`

// FindAccountByUserName finds account by username
func (a *AccountStore) FindAccountByUserName(ctx context.Context, userName string) (*model.Account, error) {
	var account = &model.Account{}
	err := a.DB.QueryRowContext(ctx, sqlFindAccountByUserName, userName).
		Scan(&account.Username,
			&account.HashedPassword,
			&account.GoogleID,
			&account.CreateAt,
			&account.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		logger.Err(ctx, "unable to query find account by username", err)
		return nil, err
	}
	return account, nil
}

// CreateAccount creates a new account
func (a *AccountStore) CreateAccount(ctx context.Context, username, hashedPassword string) error {
	_, err := a.DB.ExecContext(ctx, sqlCreateAccount, username, hashedPassword)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}
	return nil
}

// CreateAccountGoogle creates a new account using username and goole id
func (a *AccountStore) CreateAccountGoogle(ctx context.Context, username, googleID string) error {
	_, err := a.DB.ExecContext(ctx, sqlCreateAccountGoogle, username, googleID)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}
	return nil
}

func (a *AccountStore) UpdatePasswordToAccount(ctx context.Context, username, hashedPassword string) error {
	logger.Info(ctx, "update password to account", username)
	_, err := a.DB.ExecContext(ctx, sqlUpdatePasswordToAccount, hashedPassword, username)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}
	return nil
}

func (a *AccountStore) UpdateGoogleIDToAccount(ctx context.Context, username, googleID string) error {
	logger.Info(ctx, "update google id to account", username, googleID)
	_, err := a.DB.ExecContext(ctx, sqlUpdateGoogleIDToAccount, googleID, username)
	if err != nil {
		logger.Err(ctx, err)
		return err
	}
	return nil
}
