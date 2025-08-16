package repository

import (
	"arvanch/model"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type MessageRepository interface {
	InsertMessage(*model.Message) error

	InsertUserWithAccount(userID, name string) error

	GetUserMessages(userID string) ([]model.Message, error)

	GetUserProfile(userID string) (model.Profile, error)

	IncrementAccountBalance(accountID string, amount int64) error
}

type MessageRepo struct {
	db *gorm.DB
	MessageRepository
}

func NewMessageRepo(db *gorm.DB) MessageRepository {
	return &MessageRepo{db: db}
}

func (m *MessageRepo) InsertMessage(msg *model.Message) error {
	return m.db.Create(msg).Error
}

func (m *MessageRepo) GetUserMessages(userID string) ([]model.Message, error) {
	var messages []model.Message

	err := m.db.
		Where("user_id = ?", userID).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageRepo) GetUserProfile(userID string) (model.Profile, error) {
	var profile model.Profile

	err := m.db.Table("users").
		Select("users.id, users.name, users.account_id, accounts.id as account_id, accounts.balance").
		Joins("left join accounts on accounts.id = users.account_id").
		Where("users.id = ?", userID).
		Scan(&profile).Error

	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (m *MessageRepo) InsertUserWithAccount(userID, name string) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// Create account
		account := model.Account{
			ID: uuid.New().String(), // or any ID generator
		}
		if err := tx.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		// Create user linked to the account
		user := model.User{
			ID:        userID,
			Name:      name,
			AccountID: account.ID,
		}
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
}

func (m *MessageRepo) IncrementAccountBalance(accountID string, amount int64) error {
	result := m.db.Model(&model.Account{}).
		Where("id = ?", accountID).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient balance")
	}

	return nil
}
