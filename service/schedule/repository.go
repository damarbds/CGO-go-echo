package schedule

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context,a models.Schedule)(*string,error)
	DeleteByTransId(ctx context.Context,transId *string)error
}
