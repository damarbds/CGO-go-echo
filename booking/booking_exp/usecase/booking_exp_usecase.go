package usecase

import (
	"encoding/json"
	"github.com/auth/identityserver"
	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/skip2/go-qrcode"
	"golang.org/x/net/context"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type bookingExpUsecase struct {
	bookingExpRepo booking_exp.Repository
	userUsecase    user.Usecase
	isUsecase		identityserver.Usecase
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewbookingExpUsecase(a booking_exp.Repository, u user.Usecase, is identityserver.Usecase,timeout time.Duration) booking_exp.Usecase {
	return &bookingExpUsecase{
		bookingExpRepo: a,
		userUsecase:    u,
		isUsecase:is,
		contextTimeout: timeout,
	}
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

func generateQRCode(content string) (*string,error){
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil,err
	}
	name, err := generateRandomString(5)
	if err != nil {
		return nil,err
	}

	fileName := name + ".png"
	err = ioutil.WriteFile(fileName, png, 0700)
	copy , err:= os.Open(fileName)
	if err != nil {
		return nil,err
	}
	copy.Close()
	return &fileName,nil

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
func (b bookingExpUsecase) GetDetailBookingID(c context.Context, bookingId string) (*models.BookingExpDetailDto, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	getDetailBooking ,err := b.bookingExpRepo.GetDetailBookingID(ctx,bookingId)
	if err != nil {
		return nil,err
	}
	var bookedBy []models.BookedByObj
	var guestDesc []models.GuestDescObj
	var expType   []string
	if getDetailBooking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil,models.ErrInternalServerError
		}
	}
	if getDetailBooking.GuestDesc != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.GuestDesc), &guestDesc); errUnmarshal != nil {
			return nil,models.ErrInternalServerError
		}
	}
	if getDetailBooking.ExpType != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.ExpType), &expType); errUnmarshal != nil {
			return nil,models.ErrInternalServerError
		}
	}
	bookingExp := models.BookingExpDetailDto{
		Id:                getDetailBooking.Id,
		ExpId:             getDetailBooking.ExpId,
		OrderId:           getDetailBooking.OrderId,
		GuestDesc:         guestDesc,
		BookedBy:          bookedBy,
		BookedByEmail:     getDetailBooking.BookedByEmail,
		BookingDate:       getDetailBooking.BookingDate,
		UserId:            getDetailBooking.UserId,
		Status:            getDetailBooking.Status,
		//TicketCode:        getDetailBooking.TicketCode,
		TicketQRCode:      getDetailBooking.TicketQRCode,
		ExperienceAddOnId: getDetailBooking.ExperienceAddOnId,
		ExpTitle:          getDetailBooking.ExpTitle,
		ExpType:expType,
		ExpPickupPlace:    getDetailBooking.ExpPickupPlace,
		ExpPickupTime:     getDetailBooking.ExpPickupTime,
		TotalPrice:        getDetailBooking.TotalPrice,
		PaymentType:       getDetailBooking.PaymentType,
	}
	return &bookingExp,nil

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
	fileNameQrCode,err := generateQRCode(orderId)
	if err != nil {
		return nil,models.ErrInternalServerError,nil
	}
	imagePath, _ := b.isUsecase.UploadFileToBlob(*fileNameQrCode, "TicketBookingQRCode")

	errRemove := os.Remove(*fileNameQrCode)
	if errRemove != nil {
		return nil,models.ErrInternalServerError,nil
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
