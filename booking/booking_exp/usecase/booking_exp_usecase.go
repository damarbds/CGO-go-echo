package usecase

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

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
	contextTimeout  time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewbookingExpUsecase(a booking_exp.Repository, u user.Usecase, m merchant.Usecase, is identityserver.Usecase, timeout time.Duration) booking_exp.Usecase {
	return &bookingExpUsecase{
		bookingExpRepo:  a,
		userUsecase:     u,
		merchantUsecase: m,
		isUsecase:       is,
		contextTimeout:  timeout,
	}
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
			ExpId:       b.ExpId,
			ExpTitle:    b.ExpTitle,
			BookingDate: b.BookingDate,
			ExpDuration: b.ExpDuration,
			TotalGuest:  len(guestDesc),
			City:        b.City,
			Province:    b.Province,
			Country:     b.Country,
		}
	}

	return myBooking, nil
}

func (b bookingExpUsecase) GetDetailBookingID(c context.Context, bookingId string) (*models.BookingExpDetailDto, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	getDetailBooking, err := b.bookingExpRepo.GetDetailBookingID(ctx, bookingId)
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
	if getDetailBooking.ExpType != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.ExpType), &expType); errUnmarshal != nil {
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
	bookingExp := models.BookingExpDetailDto{
		Id:                     getDetailBooking.Id,
		ExpId:                  getDetailBooking.ExpId,
		OrderId:                getDetailBooking.OrderId,
		GuestDesc:              guestDesc,
		BookedBy:               bookedBy,
		BookedByEmail:          getDetailBooking.BookedByEmail,
		BookingDate:            getDetailBooking.BookingDate,
		ExpiredDatePayment:     getDetailBooking.ExpiredDatePayment,
		CreatedDateTransaction: getDetailBooking.CreatedDateTransaction,
		UserId:                 getDetailBooking.UserId,
		Status:                 getDetailBooking.Status,
		//TicketCode:        getDetailBooking.TicketCode,
		TicketQRCode:        getDetailBooking.TicketQRCode,
		ExperienceAddOnId:   getDetailBooking.ExperienceAddOnId,
		ExpTitle:            getDetailBooking.ExpTitle,
		ExpType:             expType,
		ExpPickupPlace:      getDetailBooking.ExpPickupPlace,
		ExpPickupTime:       getDetailBooking.ExpPickupTime,
		TotalPrice:          getDetailBooking.TotalPrice,
		Currency:            currency,
		PaymentType:         getDetailBooking.PaymentType,
		AccountNumber:       accountBank.AccNumber,
		AccountHolder:       accountBank.AccHolder,
		BankIcon:            getDetailBooking.Icon,
		ExperiencePaymentId: getDetailBooking.ExperiencePaymentId,
	}
	return &bookingExp, nil

}

func (b bookingExpUsecase) Insert(c context.Context, booking *models.NewBookingExpCommand, token string) (*models.NewBookingExpCommand, error, error) {

	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	if booking.ExpId == "" {
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
	bookngDate, errDate := time.Parse(layoutFormat, booking.BookingDate)
	if errDate != nil {
		return nil, errDate, nil
	}
	orderId, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
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
		ExpId:             booking.ExpId,
		OrderId:           orderId,
		GuestDesc:         booking.GuestDesc,
		BookedBy:          booking.BookedBy,
		BookedByEmail:     booking.BookedByEmail,
		BookingDate:       bookngDate,
		UserId:            booking.UserId,
		Status:            0,
		TicketCode:        ticketCode,
		TicketQRCode:      imagePath,
		ExperienceAddOnId: booking.ExperienceAddOnId,
	}
	if *bookingExp.ExperienceAddOnId == "" {
		bookingExp.ExperienceAddOnId = nil
	}
	if *bookingExp.UserId == "" {
		bookingExp.UserId = nil
	}
	res, err := b.bookingExpRepo.Insert(ctx, &bookingExp)
	if err != nil {
		return nil, err, nil
	}
	booking.Id = res.Id
	return booking, nil, nil
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
