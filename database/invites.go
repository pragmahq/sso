package database

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type InviteCode struct {
	Id          string     `pg:"id,pk"`
	Code        string     `pg:"code,unique"`
	GeneratedBy string     `pg:"generated_by"`
	UsedBy      string     `pg:"used_by"`
	CreatedAt   time.Time  `pg:"created_at"`
	UsedAt      *time.Time `pg:"used_at"`
}

func (i *InviteCode) Create(db *DB) error {
	_, err := db.Model(i).Insert()
	return err
}

func GetInviteCodeByCode(db *DB, code string) (*InviteCode, error) {
	inviteCode := &InviteCode{}
	err := db.Model(inviteCode).
		Where("code = ?", code).
		Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return inviteCode, nil
}

func GenerateInviteCode(db *DB, generatedBy string) (*InviteCode, error) {
	code := uuid.New().String()
	inviteCode := &InviteCode{
		Id:          uuid.New().String(),
		Code:        code,
		GeneratedBy: generatedBy,
		CreatedAt:   time.Now(),
	}
	err := inviteCode.Create(db)
	if err != nil {
		return nil, err
	}
	return inviteCode, nil
}
