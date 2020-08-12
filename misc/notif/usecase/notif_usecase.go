package usecase

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/booking/booking_exp"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/auth/merchant"
	"github.com/misc/notif"
	"github.com/models"
)

type notifUsecase struct {
	keyFCM 			string
	notifRepo       notif.Repository
	merchantUsecase merchant.Usecase
	contextTimeout  time.Duration
	bookingRepo booking_exp.Repository
}





func NewNotifUsecase(keyFCM string,bookingRepo booking_exp.Repository,n notif.Repository, u merchant.Usecase, timeout time.Duration) notif.Usecase {
	return &notifUsecase{
		keyFCM:keyFCM,
		bookingRepo:bookingRepo,
		notifRepo:       n,
		merchantUsecase: u,
		contextTimeout:  timeout,
	}
}
func (x notifUsecase) DeleteNotificationByIds(ctx context.Context, token string, ids string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, x.contextTimeout)
	defer cancel()

	currentUserMerchant,err := x.merchantUsecase.ValidateTokenMerchant(ctx ,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}
	err = x.notifRepo.DeleteNotificationByIds(ctx,currentUserMerchant.Id,ids,currentUserMerchant.MerchantEmail,time.Now())
	if err != nil{
		return nil,errors.New(err.Error())
	}
	result := models.ResponseDelete{
		Id:      currentUserMerchant.Id,
		Message: "Success Deleted Notification",
	}
	return &result,nil
}
func (x notifUsecase) UpdateStatusNotif(ctx context.Context, notif models.NotificationRead, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, x.contextTimeout)
	defer cancel()

	currentUserMerchant,err := x.merchantUsecase.ValidateTokenMerchant(ctx ,token)
	if err != nil{
		return nil,models.ErrUnAuthorize
	}
	err = x.notifRepo.UpdateStatusNotif(ctx,notif,currentUserMerchant.MerchantEmail,time.Now())
	if err != nil{
		return nil,errors.New(err.Error())
	}
	result := models.ResponseDelete{
		Id:      currentUserMerchant.Id,
		Message: "Success Update Read Notification",
	}
	return &result,nil
}
func (t notifUsecase) FCMPushNotification(ctx context.Context, ar models.FCMPushNotif) (*models.ResponseDelete, error) {
	data, _ := json.Marshal(ar)

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(data))
	//os.Exit(1)
	req.Header.Set("Authorization", "Bearer " + t.keyFCM)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil,errors.New("Error Push Notif FCM")
	}
	//user := models.NewCommandSchedule{}
	//json.NewDecoder(resp.Body).Decode(&user)
	result := models.ResponseDelete{
		Id:      ar.To,
		Message: "Success Push Notif",
	}
	return &result,nil
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
	countDeleteIndex := 0
	notifs := make([]*models.NotifDto, 0)
	for _, n := range res {
		notifType := "General"
		if n.Type == 2 {
			notifType = "Corporation"
		}
		notif := &models.NotifDto{
			Id:    n.Id,
			Type:  notifType,
			Title: n.Title,
			Desc:  n.Desc,
			Date:  n.CreatedDate,
		}
		if n.ScheduleId != nil && n.BookingExpId != nil{
			getDetailBooking,err := x.bookingRepo.GetDetailTransportBookingID(ctx,*n.BookingExpId,*n.BookingExpId,nil)
			if err != nil {
				countDeleteIndex = countDeleteIndex + 1
				//return nil,err
			}else {
				notif.TransId = getDetailBooking[0].TransId
				notif.DepartureTime = getDetailBooking[0].DepartureTime
				notif.ArrivalTime = getDetailBooking[0].ArrivalTime
				notif.TransName = getDetailBooking[0].TransId
				notif.HarborDestName = getDetailBooking[0].HarborDestName
				notif.HarborSourceName = getDetailBooking[0].HarborSourceName
				notifs = append(notifs,notif)
			}
		}else if n.ExpId != nil && n.BookingExpId != nil{
			getDetailBooking,err := x.bookingRepo.GetDetailBookingID(ctx,*n.BookingExpId,*n.BookingExpId)
			if err != nil {
				countDeleteIndex = countDeleteIndex + 1
				//return nil,err
			}else {
				notif.ExpId = getDetailBooking.ExpId
				notif.ExpTitle = getDetailBooking.ExpTitle
				notifs = append(notifs,notif)
			}
		}

	}


	totalRecords, _ := x.notifRepo.GetCountByMerchantID(ctx, currentMerchant.Id)
	totalRecords = totalRecords - countDeleteIndex

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
		RecordPerPage: len(notifs),
	}

	response := &models.NotifWithPagination{
		Data: notifs,
		Meta: meta,
	}

	return response, nil
}
