package repositories

import (
	"context"

	"github.com/devrapture/pod-events/internal/database"
	"gorm.io/gorm"
)

func dbFromCtx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx := database.TxFromContext(ctx); tx != nil {
		return tx
	}
	return defaultDB
}
