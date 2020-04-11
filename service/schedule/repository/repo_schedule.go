package repository

import (
	"context"
	"database/sql"
	//"fmt"
	guuid "github.com/google/uuid"
	"github.com/service/schedule"

	"time"

	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type scheduleRepository struct {
	Conn *sql.DB
}



// NewpromoRepository will create an object that represent the article.Repository interface
func NewScheduleRepository(Conn *sql.DB) schedule.Repository {
	return &scheduleRepository{Conn}
}
func (s scheduleRepository) Insert(ctx context.Context, a models.Schedule) (*string, error) {
	a.Id = guuid.New().String()
	query := `INSERT schedules SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_id=?,arrival_time=?,departure_time=?,day=?,month=?,
				year=?,price=?,departure_timeoption_id=?,arrival_timeoption_id=?`
	stmt, err := s.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil,err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.TransId,a.ArrivalTime,
		a.DepartureTime,a.Day,a.Month,a.Year, a.Price,a.DepartureTimeoptionId,a.ArrivalTimeoptionId)
	if err != nil {
		return nil,err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id,nil
}
func (s scheduleRepository) DeleteByTransId(ctx context.Context, transId *string) error {
	query := "DELETE FROM schedules WHERE trans_id = ?"

	stmt, err := s.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, transId)
	if err != nil {

		return err
	}

	//rowsAfected, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//
	//if rowsAfected != 1 {
	//	err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
	//	return err
	//}

	return nil
}