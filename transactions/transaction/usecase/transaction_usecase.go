package usecase

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/auth/admin"
	"github.com/auth/merchant"
	"github.com/service/promo"

	"github.com/models"
	"github.com/service/exp_payment"
	"github.com/transactions/transaction"
)

type transactionUsecase struct {
	adminUsecase              admin.Usecase
	merchantUsecase           merchant.Usecase
	experiencePaymentTypeRepo exp_payment.Repository
	transactionRepo           transaction.Repository
	contextTimeout            time.Duration
	promoRepo                 promo.Repository
}

func NewTransactionUsecase(promoRepo promo.Repository, au admin.Usecase, mu merchant.Usecase, ep exp_payment.Repository, t transaction.Repository, timeout time.Duration) transaction.Usecase {
	return &transactionUsecase{
		promoRepo:                 promoRepo,
		adminUsecase:              au,
		merchantUsecase:           mu,
		experiencePaymentTypeRepo: ep,
		transactionRepo:           t,
		contextTimeout:            timeout,
	}
}

func (t transactionUsecase) GetDetailTransactionSchedule(ctx context.Context, date string, transId string, expId string, token string, status string) (*models.TransactionScheduleDto, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	currentMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}
	listTransactions, err := t.transactionRepo.GetTransactionByExpIdORTransId(ctx, date, expId, transId, currentMerchant.Id, status)
	if err != nil {
		return nil, err
	}
	if len(listTransactions) == 0 {
		return nil, nil
	}
	var expType []string
	if listTransactions[0].ExpType != "" {
		if !strings.Contains(listTransactions[0].ExpType, "]") {
			// Default type for Transportation
			expType = []string{"Transportation"}
		} else {
			if errUnmarshal := json.Unmarshal([]byte(listTransactions[0].ExpType), &expType); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}
	}
	var result models.TransactionScheduleDto
	if expType[0] == "Transportation" {
		result.TransId = &listTransactions[0].ExpId
		result.TransTo = listTransactions[0].CountryName
		result.TransFrom = listTransactions[0].ProvinceName
		result.ArrivalTime = listTransactions[0].ArrivalTime
		result.DepartureTime = listTransactions[0].DepartureTime
	} else {
		result.ExpId = &listTransactions[0].ExpId
		result.ExpTitle = &listTransactions[0].ExpTitle
	}
	result.Transactions = make([]models.TransactionBooked, len(listTransactions))
	for i, item := range listTransactions {
		var status string
		if item.TransactionStatus == 0 {
			status = "Pending"
		} else if item.TransactionStatus == 1 {
			status = "Waiting approval"
		} else if item.TransactionStatus == 2 &&
			(item.CheckInDate.Format("02 January 2006") <= time.Now().AddDate(0, 0, 14).Format("02 January 2006") &&
				item.CheckInDate.Format("02 January 2006") >= time.Now().Format("02 January 2006")) {
			status = "Upcoming"
		} else if item.TransactionStatus == 2 && (item.CheckInDate.Format("2006-01-02") >= time.Now().Format("2006-01-02")) {
			status = "Confirm"
		} else if item.TransactionStatus == 2 && (item.CheckInDate.Format("2006-01-02") < time.Now().Format("2006-01-02")) {
			status = "Finished"
		} else if item.TransactionStatus == 3 || item.TransactionStatus == 4 {
			status = "Failed"
		} else if item.TransactionStatus == 2 && item.BookingStatus == 3 {
			status = "Boarded"
		}
		var guestDesc []models.GuestDescObj
		if item.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}
		var bookedBy []models.BookedByObj
		if item.BookedBy != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.BookedBy), &bookedBy); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}
		result.Transactions[i] = models.TransactionBooked{
			BookedBy:          bookedBy[0].FullName,
			GuestCount:        len(guestDesc),
			BookingExpId:      item.BookingExpId,
			OrderId:           item.OrderId,
			TransactionId:     &item.TransactionId,
			TransactionStatus: status,
		}
	}
	return &result, nil
}

func (t transactionUsecase) GetTransactionByDate(ctx context.Context, date string, isTransportation bool, isExperience bool, token string) ([]*models.TransactionByDateDto, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	currentMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	listTransactions, err := t.transactionRepo.GetTransactionByDate(ctx, date, isTransportation, isExperience, currentMerchant.Id)
	if err != nil {
		return nil, err
	}
	result := make([]*models.TransactionByDateDto, len(listTransactions))
	for i, item := range listTransactions {
		result[i] = &models.TransactionByDateDto{
			ExpId:         item.ExpId,
			ExpTitle:      item.ExpTitle,
			TransId:       item.TransId,
			DepartureTime: item.DepartureTime,
			ArrivalTime:   item.ArrivalTime,
			TransTo:       item.HarborsDest,
			TransFrom:     item.HarborsSource,
			MaxCapacity:   item.Capacity,
		}
		if result[i].ExpId != nil {
			getBooking, err := t.transactionRepo.GetCountByExpId(ctx, date, *result[i].ExpId, true)
			if err != nil {
				return nil, err
			}

			var guestDescs int

			if getBooking != nil {
				for _, booking := range getBooking {
					guestDesc := make([]models.GuestDescObj, 0)
					if errUnmarshal := json.Unmarshal([]byte(*booking), &guestDesc); errUnmarshal != nil {
						return nil, models.ErrInternalServerError
					}
					guestDescs = guestDescs + len(guestDesc)
				}
			}
			result[i].SeatTaken = guestDescs
		} else if result[i].TransId != nil {
			getbooking, _ := t.transactionRepo.GetCountByTransId(ctx, *result[i].TransId, true, date)

			var getbookingCount int
			if getbooking != nil {
				for _, booking := range getbooking {
					guestDesc := make([]models.GuestDescObj, 0)
					if errUnmarshal := json.Unmarshal([]byte(*booking), &guestDesc); errUnmarshal != nil {
						return nil, models.ErrInternalServerError
					}
					getbookingCount = getbookingCount + len(guestDesc)
				}
			}
			result[i].SeatTaken = getbookingCount
		}
	}
	return result, nil
}

func (t transactionUsecase) CountThisMonth(ctx context.Context) (*models.TotalTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	total, err := t.transactionRepo.CountThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	return total, nil
}

func (t transactionUsecase) List(ctx context.Context, startDate, endDate, search, status string, page, limit, offset *int, token string, isAdmin bool, isTransportation bool, isExperience bool, isSchedule bool, tripType, paymentType, activityType string, confirmType string, class string, departureTimeStart string, departureTimeEnd string, arrivalTimeStart string, arrivalTimeEnd string, transactionId string) (*models.TransactionWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	var merchantId string

	list := make([]*models.TransactionOut, 0)

	if token != "" && isAdmin == true {
		currentMerchant, err := t.adminUsecase.ValidateTokenAdmin(ctx, token)
		if err != nil {
			return nil, models.ErrUnAuthorize
		}
		merchantId = currentMerchant.Id

		list, err = t.transactionRepo.List(ctx, startDate, endDate, search, status, limit, offset, "", isTransportation, isExperience, isSchedule, tripType, paymentType, activityType, confirmType, class, departureTimeStart, departureTimeEnd, arrivalTimeStart, arrivalTimeEnd, transactionId)

		if err != nil {
			return nil, err
		}
	} else if token != "" {
		currentMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
		if err != nil {
			return nil, models.ErrUnAuthorize
		}
		merchantId = currentMerchant.Id

		list, err = t.transactionRepo.List(ctx, startDate, endDate, search, status, limit, offset, merchantId, isTransportation, isExperience, isSchedule, tripType, paymentType, activityType, confirmType, class, departureTimeStart, departureTimeEnd, arrivalTimeStart, arrivalTimeEnd, transactionId)

		if err != nil {
			return nil, err
		}
	}

	transactions := make([]*models.TransactionDto, len(list))
	for i, item := range list {
		var promo *models.PromoTransaction
		if item.PromoId != nil {
			getPromo, _ := t.promoRepo.GetById(ctx, *item.PromoId)
			promo = &models.PromoTransaction{
				PromoValue: getPromo.PromoValue,
				PromoType:  getPromo.PromoType,
			}
		}
		var expType []string
		if item.ExpType != "" {
			if !strings.Contains(item.ExpType, "]") {
				// Default type for Transportation
				expType = []string{"Transportation"}
			} else {
				if errUnmarshal := json.Unmarshal([]byte(item.ExpType), &expType); errUnmarshal != nil {
					return nil, errUnmarshal
				}
			}
		}
		var experiencePaymentType *models.ExperiencePaymentTypeDto
		if item.ExperiencePaymentId != nil {
			if *item.ExperiencePaymentId == "Economy" || *item.ExperiencePaymentId == "Executive" {
				// Default Payment Type for Transportation
				experiencePaymentType = &models.ExperiencePaymentTypeDto{
					Id:   "8a5e3eef-a6db-4584-a280-af5ab18a979b",
					Name: "Full Payment",
					Desc: "Full Payment",
				}
			} else {
				query, err := t.experiencePaymentTypeRepo.GetById(ctx, *item.ExperiencePaymentId)
				if err != nil {
					return nil, err
				}
				for _, element := range query {
					if element.Id == *item.ExperiencePaymentId {
						paymentType := models.ExperiencePaymentTypeDto{
							Id:   element.ExpPaymentTypeId,
							Name: element.ExpPaymentTypeName,
							Desc: element.ExpPaymentTypeDesc,
						}
						if paymentType.Name == "Down Payment" {
							if item.OriginalPrice != nil {
								paymentType.OriginalPrice = item.OriginalPrice
								remainingPayment := *item.OriginalPrice - item.TotalPrice
								paymentType.RemainingPayment = remainingPayment
							} else {
								remainingPayment := element.Price - item.TotalPrice
								paymentType.RemainingPayment = remainingPayment
							}
						}
						experiencePaymentType = &paymentType
					}
				}
			}
		}
		var guestDesc []models.GuestDescObj
		if item.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, errUnmarshal
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
		var bookedBy []models.BookedByObj
		if item.BookedBy != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.BookedBy), &bookedBy); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		var status string
		if item.TransactionStatus == 0 {
			status = "Pending"
		} else if item.TransactionStatus == 1 {
			status = "Waiting approval"
		} else if item.TransactionStatus == 2 &&
			(item.CheckInDate.Format("02 January 2006") <= time.Now().AddDate(0, 0, 14).Format("02 January 2006") &&
				item.CheckInDate.Format("02 January 2006") >= time.Now().Format("02 January 2006")) {
			status = "Upcoming"
		} else if item.TransactionStatus == 2 && (item.CheckInDate.Format("2006-01-02") >= time.Now().Format("2006-01-02")) {
			status = "Confirm"
		} else if item.TransactionStatus == 2 && (item.CheckInDate.Format("2006-01-02") < time.Now().Format("2006-01-02")) {
			status = "Finished"
		} else if item.TransactionStatus == 3 || item.TransactionStatus == 4 {
			status = "Failed"
		} else if item.TransactionStatus == 2 && item.BookingStatus == 3 {
			status = "Boarded"
		}
		transactions[i] = &models.TransactionDto{
			TransactionId:         item.TransactionId,
			ExpId:                 item.ExpId,
			ExpTitle:              item.ExpTitle,
			ExpType:               expType,
			BookingExpId:          item.BookingExpId,
			BookingCode:           item.BookingCode,
			BookingDate:           item.BookingDate,
			CheckInDate:           item.CheckInDate,
			BookedBy:              bookedBy,
			Guest:                 len(guestDesc),
			Email:                 item.Email,
			Status:                status,
			TotalPrice:            item.TotalPrice,
			ExperiencePaymentType: experiencePaymentType,
			Merchant:              item.MerchantName,
			OrderId:               item.OrderId,
			GuestCount:            expGuest,
			Promo:                 promo,
			Points:                item.Points,
			ArrivalTime:           item.ArrivalTime,
			DepartureTime:         item.DepartureTime,
		}
		if expType[0] != "Transportation" {
			transactions[i].ExpDuration = *item.ExpDuration
			transactions[i].ProvinceName = *item.ProvinceName
			transactions[i].CountryName = *item.CountryName
		} else {
			transactions[i].TransTo = *item.CountryName
			transactions[i].TransFrom = *item.ProvinceName
		}
	}
	var totalRecords int
	if token != "" && isAdmin == true {
		totalRecords, _ = t.transactionRepo.Count(ctx, startDate, endDate, search, status, "", isTransportation, isExperience, isSchedule, tripType, paymentType, activityType, confirmType, class, departureTimeStart, departureTimeEnd, arrivalTimeStart, arrivalTimeEnd, transactionId)
	} else {
		totalRecords, _ = t.transactionRepo.Count(ctx, startDate, endDate, search, status, merchantId, isTransportation, isExperience, isSchedule, tripType, paymentType, activityType, confirmType, class, departureTimeStart, departureTimeEnd, arrivalTimeStart, arrivalTimeEnd, transactionId)
	}

	if limit == nil {
		limit = &totalRecords
	}
	totalPage := int(math.Ceil(float64(totalRecords) / float64(*limit)))
	if page == nil {
		var number = 1
		page = &number
	}
	prev := page
	next := page

	if *page != 1 {
		*prev = *page - 1
	}

	if *page != totalPage {
		*next = *page + 1
	}
	meta := &models.MetaPagination{
		Page:          *page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          *prev,
		Next:          *next,
		RecordPerPage: len(list),
	}

	response := &models.TransactionWithPagination{
		Data: transactions,
		Meta: meta,
	}

	return response, nil
}

func (t transactionUsecase) CountSuccess(ctx context.Context) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	count, err := t.transactionRepo.CountSuccess(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}
