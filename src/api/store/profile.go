package store

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/model"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"database/sql"
)

type ProfileStore struct {
	DB *sql.DB
}

func NewProfileStore(db *sql.DB) *ProfileStore {
	return &ProfileStore{DB: db}
}

const sqlFindProfileByAccountID = `
SELECT
    profiles.id,
	profiles.full_name,
	profiles.phone,
	profiles.email,
	profiles.account_id,
	profiles.created_at,
	profiles.updated_at
FROM profiles
WHERE
	profiles.account_id = ?
LIMIT 1;`

const sqlInsertOrUpdateProfile = `
INSERT INTO profiles (full_name, phone, email, account_id)
VALUES (?,?,?,?)
ON DUPLICATE KEY UPDATE 
                     full_name = ?,
                     phone = ?,
                     email = ?,
                     account_id = ?
;`

// FindProfileByUserName finds profile by username
func (a *ProfileStore) FindProfileByUserName(ctx context.Context, accountID string) (*model.Profile, error) {
	var profile = &model.Profile{}
	err := a.DB.QueryRowContext(ctx, sqlFindProfileByAccountID, accountID).
		Scan(&profile.ID,
			&profile.FullName,
			&profile.Phone,
			&profile.Email,
			&profile.AccountID,
			&profile.CreatedAt,
			&profile.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		logger.Err(ctx, "unable to query find profile by accountID", err)
		return nil, err
	}
	return profile, nil
}

// InsertOrUpdateProfile creates a new profile or update if exists
func (a *ProfileStore) InsertOrUpdateProfile(ctx context.Context, profile *model.Profile) (int64, error) {
	result, err := a.DB.ExecContext(ctx, sqlInsertOrUpdateProfile,
		profile.FullName, profile.Phone, profile.Email, profile.AccountID,
		profile.FullName, profile.Phone, profile.Email, profile.AccountID)
	if err != nil {
		logger.Err(ctx, "unable to insert or update profile", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Err(ctx, "unable to get last insert ID when insert or update profile", err)
		return 0, err
	}
	return id, nil
}
