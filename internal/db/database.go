package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(path string) (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := conn.AutoMigrate(&UsageEvent{}, &ModelPricing{}, &Budget{}, &AgentKey{}); err != nil {
		return nil, err
	}

	if err := seedDefaultPricing(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// seedDefaultPricing はモデル価格テーブルが空の場合のみデフォルト値を投入する。
// ユーザーが編集した値を起動のたびに上書きしないよう、既存行がある場合は何もしない。
func seedDefaultPricing(conn *gorm.DB) error {
	var count int64
	if err := conn.Model(&ModelPricing{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return conn.Create(DefaultPricing()).Error
}
