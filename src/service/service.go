package service

import (
	"github.com/authnull0/database-service/src/models/dto"

	"github.com/authnull0/database-service/src/repository"
)

type DbService struct{}

var dbRepository = repository.DbRepository{}

func (d *DbService) DbSync(req dto.DbSyncRequest) (dto.DbSyncResponse, error) {

	response, err := dbRepository.DbSync(req)
	if err != nil {
		return dto.DbSyncResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while displaying LogEntry",
		}, nil
	}
	return response, nil
}

func (d *DbService) DbUser(req dto.DbUserRequest) (dto.DbUserResponse, error) {

	response, err := dbRepository.DbUser(req)
	if err != nil {
		return dto.DbUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while displaying LogEntry",
		}, nil
	}
	return response, nil
}
func (d *DbService) DbPrivilege(req dto.DbPrivilegeRequest) (dto.DbPrivilegeResponse, error) {

	response, err := dbRepository.DbPrivilege(req)
	if err != nil {
		return dto.DbPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while displaying LogEntry",
		}, nil
	}
	return response, nil
}
