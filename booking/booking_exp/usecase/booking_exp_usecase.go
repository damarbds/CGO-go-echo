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
	"strconv"
	"strings"
	"time"

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
	bookingExpRepo  booking_exp.Repository
	userUsecase     user.Usecase
	merchantUsecase merchant.Usecase
	isUsecase       identityserver.Usecase
	expRepo         experience.Repository
	transactionRepo transaction.Repository
	contextTimeout  time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewbookingExpUsecase(a booking_exp.Repository, u user.Usecase, m merchant.Usecase, is identityserver.Usecase, er experience.Repository, tr transaction.Repository, timeout time.Duration) booking_exp.Usecase {
	return &bookingExpUsecase{
		bookingExpRepo:  a,
		userUsecase:     u,
		merchantUsecase: m,
		isUsecase:       is,
		expRepo:         er,
		transactionRepo: tr,
		contextTimeout:  timeout,
	}
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
			if exp.ExpBookingType == "No Instant Booking" {
				transactionStatus = 1
			} else {
				transactionStatus = 2
			}
			if err := b.transactionRepo.UpdateStatus(ctx, transactionStatus, "", booking.OrderId); err != nil {
				return nil, err
			}
		} else {
			transactionStatus = 2
			if err := b.transactionRepo.UpdateStatus(ctx, transactionStatus, "", booking.OrderId); err != nil {
				return nil, err
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
			return nil, nil
		}
	}

	if res.Status == "VOIDED" {
		transactionStatus = 3
		if err := b.transactionRepo.UpdateStatus(ctx, transactionStatus, "", booking.Id); err != nil {
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
			departureTime, _ := time.Parse("15:04", *detail.DepartureTime)
			arrivalTime, _ := time.Parse("15:04", *detail.ArrivalTime)

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
	if details[0].AccountBank != "" {
		if errUnmarshal := json.Unmarshal([]byte(details[0].AccountBank), &accountBank); errUnmarshal != nil {
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
		AccountNumber:          accountBank.AccNumber,
		AccountHolder:          accountBank.AccHolder,
		BankIcon:               details[0].Icon,
		ExperiencePaymentId:    details[0].ExperiencePaymentId,
		Transportation:         transport,
	}

	return results, nil
}

func (b bookingExpUsecase) SendCharge(ctx context.Context, bookingId, paymentType string) (map[string]interface{}, error) {
	var data map[string]interface{}

	midtrans.SetupMidtrans()
	client := &http.Client{}

	booking, err := b.bookingExpRepo.GetByID(ctx, bookingId)
	if err != nil {
		fmt.Println("errGet", err.Error())
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

func (b bookingExpUsecase) GetByUserID(ctx context.Context, transactionStatus, bookingStatus int, token string) ([]*models.MyBooking, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	currentUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, err
	}

	bList, err := b.bookingExpRepo.GetByUserID(ctx, transactionStatus, bookingStatus, currentUser.Id)
	if err != nil {
		return nil, err
	}

	myBooking := make([]*models.MyBooking, len(bList))
	for i, b := range bList {
		var guestDesc []models.GuestDescObj
		if b.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(b.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, err
			}
		}

		myBooking[i] = &models.MyBooking{
			ExpId:       *b.ExpId,
			ExpTitle:    *b.ExpTitle,
			BookingDate: b.BookingDate,
			ExpDuration: *b.ExpDuration,
			TotalGuest:  len(guestDesc),
			City:        b.City,
			Province:    b.Province,
			Country:     b.Country,
		}
	}

	return myBooking, nil
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
	if getDetailBooking.AccountBank != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.AccountBank), &accountBank); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	var currency string
	if getDetailBooking.Currency == 1 {
		currency = "USD"
	} else {
		currency = "IDR"
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
	}
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
		TicketQRCode:        getDetailBooking.TicketQRCode,
		ExperienceAddOnId:   getDetailBooking.ExperienceAddOnId,
		TotalPrice:          getDetailBooking.TotalPrice,
		Currency:            currency,
		PaymentType:         getDetailBooking.PaymentType,
		AccountNumber:       accountBank.AccNumber,
		AccountHolder:       accountBank.AccHolder,
		BankIcon:            getDetailBooking.Icon,
		ExperiencePaymentId: getDetailBooking.ExperiencePaymentId,
		Experience:          expDetail,
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

func (b bookingExpUsecase) GetHistoryBookingByUserId(c context.Context, token string, monthType string) ([]*models.BookingHistoryDto, error) {
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
		query, err := b.bookingExpRepo.QueryHistoryPer30DaysByUserId(ctx, currentUserId)
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
			totalGuest := len(guestDesc)
			itemDto := models.ItemsHistoryDto{
				ExpId:          element.ExpId,
				ExpTitle:       element.ExpTitle,
				ExpType:        expType,
				ExpBookingDate: element.BookingDate,
				ExpDuration:    element.ExpDuration,
				TotalGuest:     totalGuest,
				City:           element.CityName,
				Province:       element.ProvinceName,
				Country:        element.CountryName,
				Status:         element.StatusTransaction,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}
		result = append(result, &historyDto)
	} else {
		//test := "2006-05"
		//year := string(monthType[0] + monthType[1] + monthType[2] + monthType[3])
		//month := string (monthType[5] + monthType[6])
		query, err := b.bookingExpRepo.QueryHistoryPerMonthByUserId(ctx, currentUserId, monthType)
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
			totalGuest := len(guestDesc)
			itemDto := models.ItemsHistoryDto{
				ExpId:          element.ExpId,
				ExpTitle:       element.ExpTitle,
				ExpType:        expType,
				ExpBookingDate: element.BookingDate,
				ExpDuration:    element.ExpDuration,
				TotalGuest:     totalGuest,
				City:           element.CityName,
				Province:       element.ProvinceName,
				Country:        element.CountryName,
				Status:         element.StatusTransaction,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}
		result = append(result, &historyDto)
	}
	return result, nil
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
