package usecase

import (
	"context"
	"github.com/booking/booking_exp"
	"math"
	"time"

	"github.com/auth/merchant"
	"github.com/misc/notif"
	"github.com/models"
)

type notifUsecase struct {
	notifRepo       notif.Repository
	merchantUsecase merchant.Usecase
	contextTimeout  time.Duration
	bookingRepo booking_exp.Repository
}

func NewNotifUsecase(bookingRepo booking_exp.Repository,n notif.Repository, u merchant.Usecase, timeout time.Duration) notif.Usecase {
	return &notifUsecase{
		bookingRepo:bookingRepo,
		notifRepo:       n,
		merchantUsecase: u,
		contextTimeout:  timeout,
	}
}

func (x notifUsecase) GetByMerchantID(ctx context.Context, token string,page, limit, offset int) (*models.NotifWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, x.contextTimeout)
	defer cancel()

	currentMerchant, err := x.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	res, err := x.notifRepo.GetByMerchantID(ctx, currentMerchant.Id,limit,offset)
	if err != nil {
		return nil, err
	}

	notifs := make([]*models.NotifDto, len(res))
	for i, n := range res {
		notifType := "General"
		if n.Type == 2 {
			notifType = "Corporation"
		}
		notifs[i] = &models.NotifDto{
			Id:    n.Id,
			Type:  notifType,
			Title: n.Title,
			Desc:  n.Desc,
			Date:  n.CreatedDate,
		}
		if n.ScheduleId != nil && n.BookingExpId != nil{
			getDetailBooking,err := x.bookingRepo.GetDetailTransportBookingID(ctx,*n.BookingExpId,*n.BookingExpId,nil)
			if err != nil {
				return nil,err
			}
			notifs[i].TransId = getDetailBooking[0].TransId
			notifs[i].DepartureTime = getDetailBooking[0].DepartureTime
			notifs[i].ArrivalTime = getDetailBooking[0].ArrivalTime
			notifs[i].TransName = getDetailBooking[0].TransId
			notifs[i].HarborDestName = getDetailBooking[0].HarborDestName
			notifs[i].HarborSourceName = getDetailBooking[0].HarborSourceName
		}else if n.ExpId != nil && n.BookingExpId != nil{
			getDetailBooking,err := x.bookingRepo.GetDetailBookingID(ctx,*n.BookingExpId,*n.BookingExpId)
			if err != nil {
				return nil,err
			}
			notifs[i].ExpId = getDetailBooking.ExpId
			notifs[i].ExpTitle = getDetailBooking.ExpTitle
		}

	}


		totalRecords, _ := x.notifRepo.GetCountByMerchantID(ctx, currentMerchant.Id)


	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}
	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(res),
	}

	response := &models.NotifWithPagination{
		Data: notifs,
		Meta: meta,
	}

	return response, nil
}
