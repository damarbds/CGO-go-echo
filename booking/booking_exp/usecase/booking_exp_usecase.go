package usecase

import (
	"github.com/booking/booking_exp"
	"github.com/models"
	"golang.org/x/net/context"
	"math/rand"
	"time"
)

type bookingExpUsecase struct {
	bookingExpRepo    booking_exp.Repository
	contextTimeout time.Duration
}


// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewbookingExpUsecase(a booking_exp.Repository, timeout time.Duration) booking_exp.Usecase {
	return &bookingExpUsecase{
		bookingExpRepo:    a,
		contextTimeout: timeout,
	}
}
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
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

func (b bookingExpUsecase) Insert(c context.Context, booking *models.NewBookingExpCommand) (*models.NewBookingExpCommand,error,error) {

	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	if booking.ExpId == ""{
		return nil,models.ValidationExpId,nil
	}
	if booking.BookingDate == ""{
		return nil,models.ValidationBookedDate,nil
	}
	if booking.Status == ""{
		return nil,models.ValidationStatus,nil
	}
	if booking.BookedBy == ""{
		return  nil,models.ValidationBookedBy,nil
	}
	layoutFormat := "2006-01-02 15:04:05"
	bookngDate, errDate := time.Parse(layoutFormat,booking.BookingDate)
	if errDate != nil{
		return nil,errDate,nil
	}
	orderId, err := generateRandomString(12)
	if err != nil{
		return nil,models.ErrInternalServerError,nil
	}
	booking.OrderId = orderId
	bookingExp := models.BookingExp{
		Id:            "",
		CreatedBy:     "admin",
		CreatedDate:   time.Now(),
		ModifiedBy:    nil,
		ModifiedDate:  nil,
		DeletedBy:     nil,
		DeletedDate:   nil,
		IsDeleted:     0,
		IsActive:      1,
		ExpId:         booking.ExpId,
		OrderId:orderId,
		GuestDesc:     booking.GuestDesc,
		BookedBy:      booking.BookedBy,
		BookedByEmail: booking.BookedByEmail,
		BookingDate:   bookngDate,
		UserId:        booking.UserId,
		Status:        0,
		TicketCode:    booking.TicketCode,
		TicketQRCode:  booking.TicketQRCode,
		ExperienceAddOnId:booking.ExperienceAddOnId,
	}
	if *bookingExp.ExperienceAddOnId == ""{
		bookingExp.ExperienceAddOnId = nil
	}
	if *bookingExp.UserId == ""{
		bookingExp.UserId = nil
	}
	res,err := b.bookingExpRepo.Insert(ctx, &bookingExp)
	if err != nil {
		return nil,err,nil
	}
	booking.Id = res.Id
	return booking,nil,nil
}