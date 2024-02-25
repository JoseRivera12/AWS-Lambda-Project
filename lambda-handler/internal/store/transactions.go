package store

import (
	"context"
	"fmt"
	"transaction-mailer/internal/fileutils"

	"gorm.io/gorm"
)

const batchSize int = 500

func InitDB(conn *gorm.DB) error {
	return conn.AutoMigrate(&User{}, &Transaction{})
}

func GetOrCreateUser(db *gorm.DB, ctx context.Context, user *User) (*User, error) {
	result := db.WithContext(ctx).FirstOrCreate(user, user)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating or getting user: %w", result.Error)
	}
	return user, nil
}

func insertTransactionsInBatches(db *gorm.DB, transactions []*Transaction) error {
	result := db.CreateInBatches(transactions, len(transactions))
	if result.Error != nil {
		return fmt.Errorf("error creating transactions in bulk: %w", result.Error)
	}
	return nil
}

func CreateTransactions(db *gorm.DB, ctx context.Context, transactionRecords []*fileutils.Record, user *User) error {
	user, err := GetOrCreateUser(db, ctx, user)
	if err != nil {
		return err
	}

	tx := db.WithContext(ctx).Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			fmt.Println("Error during transaction:", err)
		}
	}()

	for i := 0; i < len(transactionRecords); i += batchSize {
		end := min(i+batchSize, len(transactionRecords))
		batch := make([]*Transaction, end-i)
		for index, record := range transactionRecords[i:end] {
			transaction := &Transaction{
				UserID:      user.ID,
				Amount:      record.Transaction,
				Description: record.Description,
				Currency:    record.Currency,
				CreatedAt:   record.Date.Time,
			}
			batch[index] = transaction
		}

		if err := insertTransactionsInBatches(tx, batch); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
