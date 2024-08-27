package usecase

import (
	"loan-tracker-api/domain"
)

type LogUsecase struct {
    LogRepo domain.LogRepository
}

func NewLogUsecase(logRepo domain.LogRepository) domain.LogService {
    return &LogUsecase{LogRepo: logRepo}
}

func (lu *LogUsecase) GetSystemLogs(Type, limit, page string) ([]domain.SystemLog, domain.ErrorResponse) {
    res,err := lu.LogRepo.GetSystemLogs(Type, limit, page )
	if err != nil {
		return []domain.SystemLog{}, domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}
	return res, domain.ErrorResponse{}
}

func (lu *LogUsecase) CreateLog(log domain.SystemLog) domain.ErrorResponse {
    err := lu.LogRepo.CreateLog(log)
	if err != nil {
		return domain.ErrorResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
		}
	}
	return domain.ErrorResponse{}
}