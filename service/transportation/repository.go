package transportation

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetTransportationByBookingId(ctx context.Context,bookingIds string)(*models.TransportationJoinBooking,error)
	GetById(ctx context.Context,id string)(*models.Transportation,error)
	UpdateStatus(ctx context.Context, status int,id string,user string)error
	Insert(ctx context.Context, transportation models.Transportation) (*string, error)
	Update(ctx context.Context, transportation models.Transportation) (*string, error)
	FilterSearch(ctx context.Context, query string, limit, offset int,isMerchant bool ,qstatus string) ([]*models.TransSearch, error)
	CountFilterSearch(ctx context.Context, query string) (int, error)
	GetTransCount(ctx context.Context, merchantId string) (int, error)
	SelectIdGetByMerchantId(ctx context.Context,merchantId string,date string)([]*string,error)
	GetAllTransport(ctx context.Context)([]*models.Transportation,error)
}
