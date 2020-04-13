package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/models"
	"github.com/service/exp_photos"
)

type exp_photosUsecase struct {
	exp_photosRepo exp_photos.Repository
	contextTimeout time.Duration
}

// Newexp_photosUsecase will create new an exp_photosUsecase object representation of exp_photos.Usecase interface
func Newexp_photosUsecase(a exp_photos.Repository, timeout time.Duration) exp_photos.Usecase {
	return &exp_photosUsecase{
		exp_photosRepo: a,
		contextTimeout: timeout,
	}
}
func (m exp_photosUsecase) GetByExperienceID(c context.Context, id string) ([]models.ExpPhotosDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err := m.exp_photosRepo.GetByExperienceID(ctx, id)
	if err != nil {
		return nil, err
	}
	var exp_photoss []models.ExpPhotosDto
	for _, element := range res {
		var expPhotoImage []models.ExpPhotoImageObject
		errObject := json.Unmarshal([]byte(element.ExpPhotoImage), &expPhotoImage)
		if errObject != nil {
			//fmt.Println("Error : ",err.Error())
			return nil, models.ErrInternalServerError
		}
		exp_photos := models.ExpPhotosDto{
			Id:             element.Id,
			ExpPhotoFolder: element.ExpPhotoFolder,
			ExpPhotoImage:  expPhotoImage,
			ExpId:          element.ExpId,
		}
		exp_photoss = append(exp_photoss, exp_photos)
	}

	return exp_photoss, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
