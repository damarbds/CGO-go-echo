package usecase

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/product/reviews"

	"github.com/third-party/xendit"

	"github.com/product/experience_add_ons"
	"github.com/service/exp_payment"

	"github.com/service/experience"
	"github.com/transactions/transaction"

	"github.com/third-party/paypal"

	"github.com/third-party/midtrans"

	"github.com/auth/identityserver"
	"github.com/auth/merchant"
	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/skip2/go-qrcode"
	"golang.org/x/net/context"
)

type bookingExpUsecase struct {
	reviewRepo                reviews.Repository
	adOnsRepo                 experience_add_ons.Repository
	experiencePaymentTypeRepo exp_payment.Repository
	bookingExpRepo            booking_exp.Repository
	userUsecase               user.Usecase
	merchantUsecase           merchant.Usecase
	isUsecase                 identityserver.Usecase
	expRepo                   experience.Repository
	transactionRepo           transaction.Repository
	contextTimeout            time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewbookingExpUsecase(reviewRepo reviews.Repository, adOnsRepo experience_add_ons.Repository, ept exp_payment.Repository, a booking_exp.Repository, u user.Usecase, m merchant.Usecase, is identityserver.Usecase, er experience.Repository, tr transaction.Repository, timeout time.Duration) booking_exp.Usecase {
	return &bookingExpUsecase{
		reviewRepo:                reviewRepo,
		adOnsRepo:                 adOnsRepo,
		experiencePaymentTypeRepo: ept,
		bookingExpRepo:            a,
		userUsecase:               u,
		merchantUsecase:           m,
		isUsecase:                 is,
		expRepo:                   er,
		transactionRepo:           tr,
		contextTimeout:            timeout,
	}
}

func (b bookingExpUsecase) RemainingPaymentNotification(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()
	list, err := b.transactionRepo.GetTransactionDownPaymentByDate(ctx)
	if err != nil {
		return err
	}
	for _, element := range list {
		var bookedBy []models.BookedByObj
		if element.BookedBy != "" {
			if errUnmarshal := json.Unmarshal([]byte(element.BookedBy), &bookedBy); errUnmarshal != nil {
				return err
			}
		}
		msg := "<h1>" + element.ExpTitle + "</h1>" +
			"<p>Trip Dates :" + element.BookingDate.Format("2006-01-01") + "</p>" +
			"<p>Price :" + strconv.FormatFloat(element.Price, 'f', 6, 64) + "</p>" +
			"<p>Remaining Payment Price :" + strconv.FormatFloat(element.TotalPrice, 'f', 6, 64) + "</p>"
		pushEmail := &models.SendingEmail{
			Subject:  "Remaining Payment",
			Message:  msg,
			From:     "CGO Indonesia",
			To:       element.BookedByEmail,
			FileName: "",
		}
		_, err = b.isUsecase.SendingEmail(pushEmail)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b bookingExpUsecase) XenPayment(ctx context.Context, amount float64, tokenId, authId, orderId, paymentType string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	var result map[string]interface{}

	xendit.XenditSetup()

	booking, err := b.bookingExpRepo.GetByID(ctx, orderId)
	if err != nil {
		return nil, err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, errUnmarshal
		}
	}

	if paymentType == "BRI" {
		va := &xendit.VirtualAccount{
			Client:     xendit.XenClient.VirtualAccount,
			ExternalID: orderId,
			BankCode:   paymentType,
			Name:       bookedBy[0].FullName,
			ExpireDate: booking.ExpiredDatePayment,
		}
		resVA, err := va.CreateFixedVA(ctx)
		if err != nil {
			return result, err
		}

		var bookingCode string
		if booking.ExpId != nil {
			bookingCode = booking.Id
		} else {
			bookingCode = booking.OrderId
		}
		if err := b.transactionRepo.UpdateAfterPayment(ctx, 0, resVA.AccountNumber, "", bookingCode); err != nil {
			return nil, err
		}

		result = structToMap(resVA)
	}

	if paymentType == "cc" || (authId != "" && tokenId != "") {
		cc := &xendit.CreditCard{
			Client:     xendit.XenClient.Card,
			TokenID:    tokenId,
			AuthID:     authId,
			ExternalID: orderId,
			Amount:     amount,
			IsCapture:  true,
		}
		resCC, err := cc.CreateCharge(ctx)
		if err != nil {
			return result, err
		}

		if err := b.SetAfterCCPayment(ctx, resCC.ExternalID, resCC.MaskedCardNumber, resCC.Status); err != nil {
			return result, err
		}

		result = structToMap(resCC)
	}

	return result, nil
}

func (b bookingExpUsecase) SetAfterCCPayment(ctx context.Context, externalId, accountNumber, status string) error {
	booking, err := b.bookingExpRepo.GetByID(ctx, externalId)
	if err != nil {
		return err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return errUnmarshal
		}
	}

	var transactionStatus int
	if status == "CAPTURED" {
		if booking.ExpId != nil {
			exp, err := b.expRepo.GetByID(ctx, *booking.ExpId)
			if err != nil {
				return err
			}
			bookingDetail, err := b.GetDetailBookingID(ctx, booking.Id, "")
			if err != nil {
				return err
			}
			if exp.ExpBookingType == "No Instant Booking" {
				transactionStatus = 1
			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
				transactionStatus = 5
			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Full Payment" {
				transactionStatus = 2
			}
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, accountNumber, "", booking.Id); err != nil {
				return err
			}
		} else {
			transactionStatus = 2
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, accountNumber, "", booking.OrderId); err != nil {
				return err
			}
		}
		msg := "<p>This is your order id " + booking.OrderId + " and your ticket QR code " + booking.TicketQRCode + "</p>"
		pushEmail := &models.SendingEmail{
			Subject:  "E-Ticket cGO",
			Message:  msg,
			From:     "CGO Indonesia",
			To:       bookedBy[0].Email,
			FileName: "Ticket.pdf",
		}
		if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
			return err
		}
	} else if status == "FAILED" {
		var bookingCode string
		if booking.ExpId != nil {
			bookingCode = booking.Id
		} else {
			bookingCode = booking.OrderId
		}
		transactionStatus = 3
		if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, accountNumber, "", bookingCode); err != nil {
			return err
		}
	}

	return nil
}

func (b bookingExpUsecase) GetByGuestCount(ctx context.Context, expId string, date string, guest int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()
	getExperience, err := b.expRepo.GetByID(ctx, expId)
	if err != nil {
		return false, err
	}
	getBooking, err := b.transactionRepo.GetCountByExpId(ctx, date, expId)
	if err != nil {
		return false, err
	}
	guestDesc := make([]models.GuestDescObj, 0)
	if getBooking != nil && *getBooking != "" {
		if errUnmarshal := json.Unmarshal([]byte(*getBooking), &guestDesc); errUnmarshal != nil {
			return false, models.ErrInternalServerError
		}
	}
	var result = false
	currentAmountBooking := len(guestDesc)
	remainingSeat := getExperience.ExpMaxGuest - currentAmountBooking
	if guest > remainingSeat {
		result = true
	}
	return result, nil
}
func (b bookingExpUsecase) Verify(ctx context.Context, orderId, bookingCode string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	var result map[string]interface{}

	cfg := paypal.PaypalConfig{
		OAuthUrl: paypal.PaypalOauthUrl,
		OrderUrl: paypal.PaypalOrderUrl,
	}

	res, err := paypal.PaypalSetup(cfg, orderId)
	if err != nil {
		return nil, err
	}

	if orderId != res.ID {
		return nil, errors.New("Incorrect Paypal Order ID")
	}

	booking, err := b.bookingExpRepo.GetByID(ctx, bookingCode)
	if err != nil {
		return nil, err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, errUnmarshal
		}
	}

	var transactionStatus int
	if res.Status == "COMPLETED" {
		if booking.ExpId != nil {
			exp, err := b.expRepo.GetByID(ctx, *booking.ExpId)
			if err != nil {
				return nil, err
			}
			bookingDetail, err := b.GetDetailBookingID(ctx, booking.Id, "")
			if err != nil {
				return nil, err
			}
			if exp.ExpBookingType == "No Instant Booking" {
				transactionStatus = 1
				maxTime := time.Now().AddDate(0, 0, 1)
				msg := "<h1>" + bookingDetail.Experience[0].ExpTitle + "</h1>" +
					"<p>Trip Dates :" + bookingDetail.BookingDate.Format("2006-01-01") + "</p>" +
					"<p>Waiting for Approval Max Time:" + maxTime.Format("2006-01-02 15:04:05") + "</p>" +
					"<p>Price :" + strconv.FormatFloat(*bookingDetail.TotalPrice, 'f', 6, 64) + "</p>"
				pushEmail := &models.SendingEmail{
					Subject:  "Waiting Approval For Merchant",
					Message:  msg,
					From:     "CGO Indonesia",
					To:       bookedBy[0].Email,
					FileName: "",
				}
				if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil, err
				}
			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
				transactionStatus = 5
				//maxTime := time.Now().AddDate(0,0,1)
				//msg := "<h1>" + bookingDetail.Experience[0].ExpTitle + "</h1><p>Trip Dates :" + bookingDetail.BookingDate.Format("2006-01-01") + "</p><p>Waiting for Approval Max Time:" + maxTime.Format("2006-01-02 15:04:05")+"</p><p>Price :" + strconv.FormatFloat(*bookingDetail.TotalPrice, 'f', 6, 64) + "</p>"
				//pushEmail := &models.SendingEmail{
				//	Subject:  "Waiting Approval For Merchant",
				//	Message:  msg,
				//	From:     "CGO Indonesia",
				//	To:      bookedBy[0].Email,
				//	FileName: "",
				//}
				//if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
				//	return nil,err
				//}
			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Full Payment" {
				transactionStatus = 2
				msg := "<h1>" + bookingDetail.Experience[0].ExpTitle + "</h1>" +
					"<p>Trip Dates :" + bookingDetail.BookingDate.Format("2006-01-01") + "</p>" +
					"<p>Price :" + strconv.FormatFloat(*bookingDetail.TotalPrice, 'f', 6, 64) + "</p>"
				pushEmail := &models.SendingEmail{
					Subject:  "E-Ticket cGO",
					Message:  msg,
					From:     "CGO Indonesia",
					To:       bookedBy[0].Email,
					FileName: "Ticket.pdf",
				}
				if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil, nil
				}
			}
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, "", "", booking.OrderId); err != nil {
				return nil, err
			}
		} else {
			bookingDetail, err := b.GetDetailBookingID(ctx, booking.Id, "")
			if err != nil {
				return nil, err
			}
			msg := "<h1>" + bookingDetail.Experience[0].ExpTitle + "</h1>" +
				"<p>Trip Dates :" + bookingDetail.BookingDate.Format("2006-01-01") + "</p>" +
				"<p>Price :" + strconv.FormatFloat(*bookingDetail.TotalPrice, 'f', 6, 64) + "</p>"
			pushEmail := &models.SendingEmail{
				Subject:  "E-Ticket cGO",
				Message:  msg,
				From:     "CGO Indonesia",
				To:       bookedBy[0].Email,
				FileName: "Ticket.pdf",
			}
			if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil, nil
			}
			transactionStatus = 2
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, "", "", booking.OrderId); err != nil {
				return nil, err
			}
		}

	}

	var bookCode string
	if booking.ExpId != nil {
		bookCode = booking.Id
	} else {
		bookCode = booking.OrderId
	}
	if res.Status == "VOIDED" {
		transactionStatus = 3
		if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, "", "", bookCode); err != nil {
			return nil, err
		}
	}

	data, _ := json.Marshal(res)
	json.Unmarshal(data, &result)

	return result, nil
}

func (b bookingExpUsecase) GetDetailTransportBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	details, err := b.bookingExpRepo.GetDetailTransportBookingID(ctx, bookingId, bookingCode)
	if err != nil {
		return nil, err
	}

	transport := make([]models.BookingTransportationDetail, len(details))
	for i, detail := range details {
		var tripDuration string
		if detail.DepartureTime != nil && detail.ArrivalTime != nil {
			departureTime, _ := time.Parse("15:04:00", *detail.DepartureTime)
			arrivalTime, _ := time.Parse("15:04:00", *detail.ArrivalTime)

			tripHour := arrivalTime.Hour() - departureTime.Hour()
			tripMinute := arrivalTime.Minute() - departureTime.Minute()
			tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
		}
		transport[i] = models.BookingTransportationDetail{
			TransID:          *detail.TransId,
			TransName:        *detail.TransName,
			TransTitle:       *detail.TransTitle,
			TransStatus:      *detail.TransStatus,
			TransClass:       *detail.TransClass,
			DepartureDate:    *detail.DepartureDate,
			DepartureTime:    *detail.DepartureTime,
			ArrivalTime:      *detail.ArrivalTime,
			TripDuration:     tripDuration,
			HarborSourceName: *detail.HarborSourceName,
			HarborDestName:   *detail.HarborDestName,
			MerchantName:     detail.MerchantName.String,
			MerchantPhone:    detail.MerchantPhone.String,
			MerchantPicture:  detail.MerchantPicture.String,
		}
	}

	var bookedBy []models.BookedByObj
	var guestDesc []models.GuestDescObj
	var accountBank models.AccountDesc
	if details[0].BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(details[0].BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if details[0].GuestDesc != "" {
		if errUnmarshal := json.Unmarshal([]byte(details[0].GuestDesc), &guestDesc); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if details[0].AccountBank != nil {
		if errUnmarshal := json.Unmarshal([]byte(*details[0].AccountBank), &accountBank); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	var currency string
	if details[0].Currency == 1 {
		currency = "USD"
	} else {
		currency = "IDR"
	}
	transport[0].TotalGuest = len(guestDesc)
	if len(transport) > 1 {
		transport[1].TotalGuest = len(guestDesc)
	}
	vaNumber := ""
	if details[0].VaNumber != nil {
		vaNumber = *details[0].VaNumber
	}

	results := &models.BookingExpDetailDto{
		Id:                     details[0].Id,
		GuestDesc:              guestDesc,
		BookedBy:               bookedBy,
		BookedByEmail:          details[0].BookedByEmail,
		BookingDate:            details[0].BookingDate,
		ExpiredDatePayment:     details[0].ExpiredDatePayment,
		CreatedDateTransaction: details[0].CreatedDateTransaction,
		UserId:                 details[0].UserId,
		Status:                 details[0].Status,
		TransactionStatus:      details[0].TransactionStatus,
		OrderId:                details[0].OrderId,
		TicketQRCode:           details[0].TicketQRCode,
		ExperienceAddOnId:      details[0].ExperienceAddOnId,
		TotalPrice:             details[0].TotalPrice,
		Currency:               currency,
		PaymentType:            details[0].PaymentType,
		AccountNumber:          vaNumber,
		AccountHolder:          accountBank.AccHolder,
		BankIcon:               details[0].Icon,
		ExperiencePaymentId:    details[0].ExperiencePaymentId,
		Transportation:         transport,
		MidtransUrl:            details[0].PaymentUrl,
	}

	return results, nil
}

func (b bookingExpUsecase) SendCharge(ctx context.Context, bookingId, paymentType string) (map[string]interface{}, error) {
	var data map[string]interface{}

	midtrans.SetupMidtrans()
	client := &http.Client{}

	booking, err := b.bookingExpRepo.GetByID(ctx, bookingId)
	if err != nil {
		return nil, err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, err
		}
	}

	fullName := bookedBy[0].FullName
	email := bookedBy[0].Email

	var phoneNumber string
	if phoneStr, ok := bookedBy[0].PhoneNumber.(string); ok {
		phoneNumber = phoneStr
	} else if phoneInt, ok := bookedBy[0].PhoneNumber.(int); ok {
		phoneNumber = strconv.Itoa(phoneInt)
	}

	name := strings.Split(fullName, " ")
	var first, last string
	if len(name) < 2 {
		first = fullName
		last = fullName
	} else {
		first = name[0]
		last = name[1]
	}

	var charge midtrans.MidtransCharge
	charge.CustomerDetail = midtrans.CustomerDetail{
		FirstName: first,
		LastName:  last,
		Phone:     phoneNumber,
		Email:     email,
	}

	charge.TransactionDetails.GrossAmount = math.Round(booking.TotalPrice)
	charge.TransactionDetails.OrderID = booking.OrderId

	charge.EnablePayment = []string{paymentType}
	charge.OptionColorTheme = midtrans.OptionColorTheme{
		Primary:     "#c51f1f",
		PrimaryDark: "#1a4794",
		Secondary:   "#1fce38",
	}
	j, _ := json.Marshal(charge)
	fmt.Println(string(j))
	AUTH_STRING := b64.StdEncoding.EncodeToString([]byte(midtrans.Midclient.ServerKey + ":"))
	req, _ := http.NewRequest("POST", midtrans.TransactionEndpoint, bytes.NewBuffer(j))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+AUTH_STRING)

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return data, err
		}
		return data, nil
	} else {
		err := errors.New("MIDTRANS ERROR : " + resp.Status)
		return data, err
	}
}

func (b bookingExpUsecase) CountThisMonth(ctx context.Context) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	count, err := b.bookingExpRepo.CountThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (b bookingExpUsecase) GetGrowthByMerchantID(ctx context.Context, token string) ([]*models.BookingGrowthDto, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	currentMerchant, err := b.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	growth, err := b.bookingExpRepo.GetGrowthByMerchantID(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}

	results := make([]*models.BookingGrowthDto, len(growth))
	for i, g := range growth {
		results[i] = &models.BookingGrowthDto{
			Date:  g.Date.Format("2006-01-02"),
			Count: g.Count,
		}
	}

	return results, nil
}

func (b bookingExpUsecase) GetByUserID(ctx context.Context, status string, token string, page, limit, offset int) (*models.MyBookingWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	currentUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, err
	}
	bookingIds, err := b.bookingExpRepo.GetBookingIdByUserID(ctx, status, currentUser.Id, limit, offset)
	if err != nil {
		return nil, err
	}

	expList, err := b.bookingExpRepo.GetBookingExpByUserID(ctx, bookingIds)
	if err != nil {
		return nil, err
	}

	myBooking := make([]*models.MyBooking, len(expList))
	for i, b := range expList {
		var guestDesc []models.GuestDescObj
		if b.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(b.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, err
			}
		}
		expType := make([]string, 0)
		if b.ExpType != nil {
			if errUnmarshal := json.Unmarshal([]byte(*b.ExpType), &expType); errUnmarshal != nil {
				return nil, err
			}
		}
		var expGuest models.TotalGuestTransportation
		if len(guestDesc) > 0 {
			for _, guest := range guestDesc {
				if guest.Type == "Adult" {
					expGuest.Adult = expGuest.Adult + 1
				} else if guest.Type == "Children" {
					expGuest.Children = expGuest.Children + 1
				}
			}
		}
		myBooking[i] = &models.MyBooking{
			OrderId:     b.OrderId,
			ExpType:     expType,
			ExpId:       *b.ExpId,
			ExpTitle:    *b.ExpTitle,
			BookingDate: b.BookingDate,
			ExpDuration: *b.ExpDuration,
			TotalGuest:  len(guestDesc),
			ExpGuest:    expGuest,
			City:        b.City,
			Province:    b.Province,
			Country:     b.Country,
		}
	}

	transList, err := b.bookingExpRepo.GetBookingTransByUserID(ctx, bookingIds)
	if err != nil {
		return nil, err
	}
	for _, b := range transList {
		var guestDesc []models.GuestDescObj
		if b.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(b.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, err
			}
		}
		var transGuest models.TotalGuestTransportation
		if len(guestDesc) > 0 {
			for _, guest := range guestDesc {
				if guest.Type == "Adult" {
					transGuest.Adult = transGuest.Adult + 1
				} else if guest.Type == "Children" {
					transGuest.Children = transGuest.Children + 1
				}
			}
		}
		var tripDuration string
		if b.DepartureTime != nil && b.ArrivalTime != nil {
			departureTime, _ := time.Parse("15:04:00", *b.DepartureTime)
			arrivalTime, _ := time.Parse("15:04:00", *b.ArrivalTime)

			tripHour := arrivalTime.Hour() - departureTime.Hour()
			tripMinute := arrivalTime.Minute() - departureTime.Minute()
			tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
		}
		booking := models.MyBooking{
			OrderId:            b.OrderId,
			ExpId:              "",
			ExpTitle:           "",
			TransId:            *b.TransId,
			TransName:          *b.TransName,
			TransFrom:          *b.HarborDestName,
			TransTo:            *b.HarborSourceName,
			TransDepartureTime: b.DepartureTime,
			TransArrivalTime:   b.ArrivalTime,
			TripDuration:       tripDuration,
			TransClass:         *b.TransClass,
			TransGuest:         transGuest,
			BookingDate:        b.BookingDate,
			ExpDuration:        0,
			TotalGuest:         len(guestDesc),
			City:               b.City,
			Province:           b.Province,
			Country:            b.Country,
		}

		myBooking = append(myBooking, &booking)
	}

	totalRecords, _ := b.bookingExpRepo.GetBookingCountByUserID(ctx, status, currentUser.Id)

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
		RecordPerPage: len(myBooking),
	}

	response := &models.MyBookingWithPagination{
		Data: myBooking,
		Meta: meta,
	}
	return response, nil
}

func (b bookingExpUsecase) GetDetailBookingID(c context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	getDetailBooking, err := b.bookingExpRepo.GetDetailBookingID(ctx, bookingId, bookingCode)
	if err != nil {
		return nil, err
	}
	var bookedBy []models.BookedByObj
	var guestDesc []models.GuestDescObj
	var accountBank models.AccountDesc
	var expType []string
	if getDetailBooking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if getDetailBooking.GuestDesc != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.GuestDesc), &guestDesc); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if getDetailBooking.ExpType != nil {
		if errUnmarshal := json.Unmarshal([]byte(*getDetailBooking.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if getDetailBooking.AccountBank != nil {
		if errUnmarshal := json.Unmarshal([]byte(*getDetailBooking.AccountBank), &accountBank); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	var currency string
	if getDetailBooking.Currency == 1 {
		currency = "USD"
	} else {
		currency = "IDR"
	}
	var experiencePaymentType *models.ExperiencePaymentTypeDto
	if getDetailBooking.ExperiencePaymentId != "" {
		query, err := b.experiencePaymentTypeRepo.GetByExpID(ctx, *getDetailBooking.ExpId)
		if err != nil {

		}
		for _, element := range query {
			if element.Id == getDetailBooking.ExperiencePaymentId {
				paymentType := models.ExperiencePaymentTypeDto{
					Id:   element.ExpPaymentTypeId,
					Name: element.ExpPaymentTypeName,
					Desc: element.ExpPaymentTypeDesc,
				}
				if paymentType.Name == "Down Payment" {
					remainingPayment := element.Price - *getDetailBooking.TotalPrice
					paymentType.RemainingPayment = remainingPayment
				} else {
					paymentType.RemainingPayment = 0
				}
				experiencePaymentType = &paymentType
			}
		}
	}
	expAddOns := make([]models.ExperienceAddOnObj, 0)
	expAddOnsQuery, errorQuery := b.adOnsRepo.GetByExpId(ctx, *getDetailBooking.ExpId)
	if errorQuery != nil {
		return nil, err
	}
	if expAddOnsQuery != nil {
		for _, element := range expAddOnsQuery {
			var currency string
			if element.Currency == 1 {
				currency = "USD"
			} else {
				currency = "IDR"
			}
			addOns := models.ExperienceAddOnObj{
				Id:       element.Id,
				Name:     element.Name,
				Desc:     element.Desc,
				Currency: currency,
				Amount:   element.Amount,
			}
			expAddOns = append(expAddOns, addOns)
		}
	}

	vaNumber := ""
	if getDetailBooking.VaNumber != nil {
		vaNumber = *getDetailBooking.VaNumber
	}

	expDetail := make([]models.BookingExpDetail, 1)
	expDetail[0] = models.BookingExpDetail{
		ExpId:           *getDetailBooking.ExpId,
		ExpTitle:        *getDetailBooking.ExpTitle,
		ExpType:         expType,
		ExpPickupPlace:  *getDetailBooking.ExpPickupPlace,
		ExpPickupTime:   *getDetailBooking.ExpPickupTime,
		MerchantName:    getDetailBooking.MerchantName.String,
		MerchantPhone:   getDetailBooking.MerchantPhone.String,
		MerchantPicture: getDetailBooking.MerchantPicture.String,
		TotalGuest:      len(guestDesc),
		City:            getDetailBooking.City,
		ProvinceName:    getDetailBooking.Province,
		ExpDuration:     *getDetailBooking.ExpDuration,
		HarborsName:     *getDetailBooking.HarborsName,
		ExperienceAddOn: expAddOns,
		CountryName:getDetailBooking.Country,
	}
	if getDetailBooking.UserId == nil {
		getDetailBooking.UserId = new(string)
	}
	if getDetailBooking.ExpId == nil {
		getDetailBooking.ExpId = new(string)
	}
	reviews, _ := b.reviewRepo.GetByExpId(ctx, *getDetailBooking.ExpId, "", 0, 1, 0, *getDetailBooking.UserId)
	if err != nil {
		return nil, err
	}
	var isReview = false
	bookingExp := models.BookingExpDetailDto{
		Id:                     getDetailBooking.Id,
		OrderId:                getDetailBooking.OrderId,
		GuestDesc:              guestDesc,
		BookedBy:               bookedBy,
		BookedByEmail:          getDetailBooking.BookedByEmail,
		BookingDate:            getDetailBooking.BookingDate,
		ExpiredDatePayment:     getDetailBooking.ExpiredDatePayment,
		CreatedDateTransaction: getDetailBooking.CreatedDateTransaction,
		UserId:                 getDetailBooking.UserId,
		Status:                 getDetailBooking.Status,
		TransactionStatus:      getDetailBooking.TransactionStatus,
		//TicketCode:        getDetailBooking.TicketCode,
		TicketQRCode:          getDetailBooking.TicketQRCode,
		ExperienceAddOnId:     getDetailBooking.ExperienceAddOnId,
		TotalPrice:            getDetailBooking.TotalPrice,
		Currency:              currency,
		PaymentType:           getDetailBooking.PaymentType,
		AccountNumber:         vaNumber,
		AccountHolder:         accountBank.AccHolder,
		BankIcon:              getDetailBooking.Icon,
		ExperiencePaymentId:   getDetailBooking.ExperiencePaymentId,
		Experience:            expDetail,
		ExperiencePaymentType: experiencePaymentType,
		IsReview:              isReview,
		MidtransUrl:           getDetailBooking.PaymentUrl,
	}
	if len(reviews) != 0 {
		desc := models.ReviewDtoObject{}
		if reviews[0].Desc != "" {
			if errUnmarshal := json.Unmarshal([]byte(reviews[0].Desc), &desc); errUnmarshal != nil {
				return nil, models.ErrInternalServerError
			}
		}
		bookingExp.ReviewDesc = &desc.Desc
		bookingExp.IsReview = true
		bookingExp.GuideReview = reviews[0].GuideReview
		bookingExp.ActivitiesReview = reviews[0].ActivitiesReview
		bookingExp.ServiceReview = reviews[0].ServiceReview
		bookingExp.CleanlinessReview = reviews[0].CleanlinessReview
		bookingExp.ValueReview = reviews[0].ValueReview
	}
	return &bookingExp, nil

}

func (b bookingExpUsecase) Insert(c context.Context, booking *models.NewBookingExpCommand, transReturnId, scheduleReturnId, token string) ([]*models.NewBookingExpCommand, error, error) {

	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	if booking.ExpId == "" && booking.TransId == nil && booking.ScheduleId != nil {
		return nil, models.ValidationExpId, nil
	}
	if booking.BookingDate == "" {
		return nil, models.ValidationBookedDate, nil
	}
	if booking.Status == "" {
		return nil, models.ValidationStatus, nil
	}
	if booking.BookedBy == "" {
		return nil, models.ValidationBookedBy, nil
	}
	layoutFormat := "2006-01-02 15:04:05"
	bookingDate, errDate := time.Parse(layoutFormat, booking.BookingDate)
	if errDate != nil {
		return nil, errDate, nil
	}
	orderId, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}

	// re-generate if duplicate order id
	if b.bookingExpRepo.CheckBookingCode(ctx, orderId) {
		orderId, err = generateRandomString(12)
		if err != nil {
			return nil, models.ErrInternalServerError, nil
		}
	}

	ticketCode, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}
	var createdBy string
	if token != "" {
		currentUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
		if err != nil {
			return nil, err, nil
		}
		createdBy = currentUser.UserEmail
	} else {
		createdBy = booking.BookedByEmail
	}
	booking.OrderId = orderId
	booking.TicketCode = ticketCode
	fileNameQrCode, err := generateQRCode(orderId)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}
	imagePath, _ := b.isUsecase.UploadFileToBlob(*fileNameQrCode, "TicketBookingQRCode")

	errRemove := os.Remove(*fileNameQrCode)
	if errRemove != nil {
		return nil, models.ErrInternalServerError, nil
	}
	booking.TicketQRCode = imagePath

	reqBooking := make([]*models.BookingExp, 0)

	bookingExp := models.BookingExp{
		Id:                "",
		CreatedBy:         createdBy,
		CreatedDate:       time.Now(),
		ModifiedBy:        nil,
		ModifiedDate:      nil,
		DeletedBy:         nil,
		DeletedDate:       nil,
		IsDeleted:         0,
		IsActive:          1,
		ExpId:             &booking.ExpId,
		OrderId:           orderId,
		GuestDesc:         booking.GuestDesc,
		BookedBy:          booking.BookedBy,
		BookedByEmail:     booking.BookedByEmail,
		BookingDate:       bookingDate,
		UserId:            booking.UserId,
		Status:            0,
		TicketCode:        ticketCode,
		TicketQRCode:      imagePath,
		ExperienceAddOnId: booking.ExperienceAddOnId,
		TransId:           booking.TransId,
		ScheduleId:        booking.ScheduleId,
	}
	if *bookingExp.ExperienceAddOnId == "" {
		bookingExp.ExperienceAddOnId = nil
	}
	if *bookingExp.UserId == "" {
		bookingExp.UserId = nil
	}
	if *bookingExp.TransId == "" {
		bookingExp.TransId = nil
	}
	if *bookingExp.ExpId == "" {
		bookingExp.ExpId = nil
	}
	if *bookingExp.ScheduleId == "" {
		bookingExp.ScheduleId = nil
	}

	reqBooking = append(reqBooking, &bookingExp)

	if transReturnId != "" && scheduleReturnId != "" {
		bookingReturn := models.BookingExp{
			Id:                "",
			CreatedBy:         createdBy,
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			ExpId:             &booking.ExpId,
			OrderId:           orderId,
			GuestDesc:         booking.GuestDesc,
			BookedBy:          booking.BookedBy,
			BookedByEmail:     booking.BookedByEmail,
			BookingDate:       bookingDate,
			UserId:            booking.UserId,
			Status:            0,
			TicketCode:        ticketCode,
			TicketQRCode:      imagePath,
			ExperienceAddOnId: booking.ExperienceAddOnId,
			TransId:           &transReturnId,
			ScheduleId:        &scheduleReturnId,
		}
		if *bookingReturn.ExperienceAddOnId == "" {
			bookingReturn.ExperienceAddOnId = nil
		}
		if *bookingReturn.UserId == "" {
			bookingReturn.UserId = nil
		}
		if *bookingReturn.TransId == "" {
			bookingReturn.TransId = nil
		}
		if *bookingReturn.ExpId == "" {
			bookingReturn.ExpId = nil
		}
		if *bookingExp.ScheduleId == "" {
			bookingExp.ScheduleId = nil
		}

		reqBooking = append(reqBooking, &bookingReturn)
	}

	resBooking := make([]*models.NewBookingExpCommand, len(reqBooking))
	for i, req := range reqBooking {
		res, err := b.bookingExpRepo.Insert(ctx, req)
		if err != nil {
			return nil, err, nil
		}
		reqBooking[i].Id = res.Id
		resBooking[i] = &models.NewBookingExpCommand{
			Id:                res.Id,
			ExpId:             booking.ExpId,
			GuestDesc:         res.GuestDesc,
			BookedBy:          res.BookedBy,
			BookedByEmail:     res.BookedByEmail,
			BookingDate:       res.BookingDate.String(),
			UserId:            res.UserId,
			Status:            strconv.Itoa(res.Status),
			OrderId:           res.OrderId,
			TicketCode:        res.TicketCode,
			TicketQRCode:      res.TicketQRCode,
			ExperienceAddOnId: res.ExperienceAddOnId,
			TransId:           res.TransId,
			ScheduleId:        res.ScheduleId,
		}
	}

	return resBooking, nil, nil
}

func (b bookingExpUsecase) GetHistoryBookingByUserId(c context.Context, token string, monthType string, page, limit, offset int) (*models.BookingHistoryDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	var currentUserId string
	if token != "" {
		validateUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
		if err != nil {
			return nil, err
		}
		currentUserId = validateUser.Id
	}
	var guestDesc []models.GuestDescObj
	var result []*models.BookingHistoryDto
	if monthType == "past-30-days" {
		bookingIds, err := b.bookingExpRepo.QuerySelectIdHistoryByUserId(ctx, currentUserId, "", limit, offset)
		if err != nil {
			return nil, err
		}
		query, err := b.bookingExpRepo.QueryHistoryPer30DaysExpByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		historyDto := models.BookingHistoryDto{
			Category: "past-30-days",
			Items:    nil,
		}
		for _, element := range query {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}

			var expGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						expGuest.Adult = expGuest.Adult + 1
					} else if guest.Type == "Children" {
						expGuest.Children = expGuest.Children + 1
					}
				}
			}
			totalGuest := len(guestDesc)
			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if element.StatusTransaction == 0 || element.StatusTransaction == 3{
					status = "Payment Expired"
				} else if element.StatusTransaction == 1 || element.StatusTransaction == 4 || element.StatusTransaction == 5 {
					status = "Cancelled"
				} else if element.StatusTransaction == 2 {
					status = "Success"
				}
			}else {
				if element.StatusTransaction == 0 && time.Now().Add(7 * time.Hour).After(element.ExpiredDatePayment.Add(7 * time.Hour)){
					status = "Payment Expired"
				}else if element.StatusTransaction == 3 || element.StatusTransaction == 4{
					status = "Cancelled"
				}
			}

			if element.UserId == nil {
				element.UserId = new(string)
			}
			checkReview, _ := b.reviewRepo.GetByExpId(ctx, element.ExpId, "", 0, 1, 0, *element.UserId)
			if err != nil {
				return nil, err
			}
			var isReview = false
			if len(checkReview) != 0 {
				isReview = true
			}

			itemDto := models.ItemsHistoryDto{
				OrderId:        element.OrderId,
				ExpId:          element.ExpId,
				ExpTitle:       element.ExpTitle,
				ExpType:        expType,
				ExpBookingDate: element.BookingDate,
				ExpDuration:    element.ExpDuration,
				TotalGuest:     totalGuest,
				ExpGuest:       expGuest,
				City:           element.CityName,
				Province:       element.ProvinceName,
				Country:        element.CountryName,
				Status:         status,
				IsReview:       isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}

		queryTrans, err := b.bookingExpRepo.QueryHistoryPer30DaysTransByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		for _, element := range queryTrans {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			var transGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						transGuest.Adult = transGuest.Adult + 1
					} else if guest.Type == "Children" {
						transGuest.Children = transGuest.Children + 1
					}
				}
			}
			//totalGuest := len(guestDesc)
			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if *element.TransactionStatus == 0 || *element.TransactionStatus == 3{
					status = "Payment Expired"
				} else if *element.TransactionStatus == 1 || *element.TransactionStatus == 4 || *element.TransactionStatus == 5 {
					status = "Cancelled"
				} else if *element.TransactionStatus == 2 {
					status = "Success"
				}
			}else {
				if *element.TransactionStatus == 0 && time.Now().Add(7 * time.Hour).After(element.ExpiredDatePayment.Add(7 * time.Hour)){
					status = "Payment Expired"
				}else if *element.TransactionStatus == 3 || *element.TransactionStatus == 4{
					status = "Cancelled"
				}
			}
			var tripDuration string
			if element.DepartureTime != nil && element.ArrivalTime != nil {
				departureTime, _ := time.Parse("15:04:00", *element.DepartureTime)
				arrivalTime, _ := time.Parse("15:04:00", *element.ArrivalTime)

				tripHour := arrivalTime.Hour() - departureTime.Hour()
				tripMinute := arrivalTime.Minute() - departureTime.Minute()
				tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
			}
			//if element.UserId == nil {
			//	*element.UserId = ""
			//}
			//if element.ExpId == nil {
			//	*element.ExpId = ""
			//}
			//checkReview ,_:= b.reviewRepo.GetByExpId(ctx , *element.ExpId,"",0,0,1,*element.UserId)
			//if err != nil {
			//	return nil,err
			//}
			var isReview = false
			//if len(checkReview) != 0 {
			//	isReview = true
			//}
			itemDto := models.ItemsHistoryDto{
				OrderId:            element.OrderId,
				ExpId:              "",
				ExpTitle:           "",
				ExpType:            nil,
				TransId:            *element.TransId,
				TransName:          *element.TransName,
				TransFrom:          *element.HarborDestName,
				TransTo:            *element.HarborSourceName,
				TransDepartureTime: element.DepartureTime,
				TransArrivalTime:   element.ArrivalTime,
				TripDuration:       tripDuration,
				TransClass:         *element.TransClass,
				TransGuest:         transGuest,
				ExpBookingDate:     element.BookingDate,
				ExpDuration:        0,
				TotalGuest:         0,
				City:               element.City,
				Province:           element.Province,
				Country:            element.Country,
				Status:             status,
				IsReview:           isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}
		result = append(result, &historyDto)
	} else {
		bookingIds, err := b.bookingExpRepo.QuerySelectIdHistoryByUserId(ctx, currentUserId, monthType, limit, offset)
		if err != nil {
			return nil, err
		}
		queryExp, err := b.bookingExpRepo.QueryHistoryPerMonthExpByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		monthType = monthType + "-" + "01" + " 00:00:00"
		layoutFormat := "2006-01-02 15:04:05"
		dt, _ := time.Parse(layoutFormat, monthType)
		dtstr2 := dt.Format("Jan '06")
		historyDto := models.BookingHistoryDto{
			Category: dtstr2,
			Items:    nil,
		}
		for _, element := range queryExp {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			var expGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						expGuest.Adult = expGuest.Adult + 1
					} else if guest.Type == "Children" {
						expGuest.Children = expGuest.Children + 1
					}
				}
			}
			totalGuest := len(guestDesc)

			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if element.StatusTransaction == 0 || element.StatusTransaction == 3{
					status = "Payment Expired"
				} else if element.StatusTransaction == 1 || element.StatusTransaction == 4 || element.StatusTransaction == 5 {
					status = "Cancelled"
				} else if element.StatusTransaction == 2 {
					status = "Success"
				}
			}else {
				if element.StatusTransaction == 0 && time.Now().Add(7 * time.Hour).After(element.ExpiredDatePayment.Add(7 * time.Hour)){
					status = "Payment Expired"
				}else if element.StatusTransaction == 3 || element.StatusTransaction == 4{
					status = "Cancelled"
				}
			}
			if element.UserId == nil {
				element.UserId = new(string)
			}
			checkReview, err := b.reviewRepo.GetByExpId(ctx, element.ExpId, "", 0, 1, 0, *element.UserId)
			if err != nil {
				return nil, err
			}
			var isReview = false
			if len(checkReview) != 0 {
				isReview = true
			}
			itemDto := models.ItemsHistoryDto{
				OrderId:        element.OrderId,
				ExpId:          element.ExpId,
				ExpTitle:       element.ExpTitle,
				ExpType:        expType,
				ExpBookingDate: element.BookingDate,
				ExpDuration:    element.ExpDuration,
				TotalGuest:     totalGuest,
				ExpGuest:       expGuest,
				City:           element.CityName,
				Province:       element.ProvinceName,
				Country:        element.CountryName,
				Status:         status,
				IsReview:       isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}

		queryTrans, err := b.bookingExpRepo.QueryHistoryPerMonthTransByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		for _, element := range queryTrans {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			var transGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						transGuest.Adult = transGuest.Adult + 1
					} else if guest.Type == "Children" {
						transGuest.Children = transGuest.Children + 1
					}
				}
			}
			//totalGuest := len(guestDesc)
			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if *element.TransactionStatus == 0 || *element.TransactionStatus == 3{
					status = "Payment Expired"
				} else if *element.TransactionStatus == 1 || *element.TransactionStatus == 4 || *element.TransactionStatus == 5 {
					status = "Cancelled"
				} else if *element.TransactionStatus == 2 {
					status = "Success"
				}
			}else {
				if *element.TransactionStatus == 0 && time.Now().Add(7 * time.Hour).After(element.ExpiredDatePayment.Add(7 * time.Hour)){
					status = "Payment Expired"
				}else if *element.TransactionStatus == 3 || *element.TransactionStatus == 4{
					status = "Cancelled"
				}
			}
			var tripDuration string
			if element.DepartureTime != nil && element.ArrivalTime != nil {
				departureTime, _ := time.Parse("15:04:00", *element.DepartureTime)
				arrivalTime, _ := time.Parse("15:04:00", *element.ArrivalTime)

				tripHour := arrivalTime.Hour() - departureTime.Hour()
				tripMinute := arrivalTime.Minute() - departureTime.Minute()
				tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
			}
			//if element.UserId == nil {
			//	*element.UserId = ""
			//}
			//if element.ExpId == nil {
			//	*element.ExpId = ""
			//}
			//checkReview ,_:= b.reviewRepo.GetByExpId(ctx , *element.ExpId,"",0,0,1,*element.UserId)
			//if err != nil {
			//	return nil,err
			//}
			var isReview = false
			//if len(checkReview) != 0 {
			//	isReview = true
			//}
			itemDto := models.ItemsHistoryDto{
				OrderId:            element.OrderId,
				ExpId:              "",
				ExpTitle:           "",
				ExpType:            nil,
				TransId:            *element.TransId,
				TransName:          *element.TransName,
				TransFrom:          *element.HarborDestName,
				TransTo:            *element.HarborSourceName,
				TransDepartureTime: element.DepartureTime,
				TransArrivalTime:   element.ArrivalTime,
				TripDuration:       tripDuration,
				TransClass:         *element.TransClass,
				TransGuest:         transGuest,
				ExpBookingDate:     element.BookingDate,
				ExpDuration:        0,
				TotalGuest:         0,
				City:               element.City,
				Province:           element.Province,
				Country:            element.Country,
				Status:             status,
				IsReview:           isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}
		result = append(result, &historyDto)
	}
	var totalRecords int
	if monthType == "past-30-days" {
		totalRecords, _ = b.bookingExpRepo.QueryCountHistoryByUserId(ctx, currentUserId, "")
	} else {
		totalRecords, _ = b.bookingExpRepo.QueryCountHistoryByUserId(ctx, currentUserId, monthType)
	}
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
		RecordPerPage: len(result[0].Items),
	}

	response := &models.BookingHistoryDtoWithPagination{
		Data: result,
		Meta: meta,
	}
	return response, nil
}

func generateQRCode(content string) (*string, error) {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	name, err := generateRandomString(5)
	if err != nil {
		return nil, err
	}

	fileName := name + ".png"
	err = ioutil.WriteFile(fileName, png, 0700)
	copy, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	copy.Close()
	return &fileName, nil

	//err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr.png")

}
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
