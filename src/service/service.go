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
			Message: "Error while inserting database",
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
			Message: "Error while inserting user",
		}, nil
	}
	return response, nil
}

// func (d *DbService) DbPrivilege(req dto.DbPrivilegeRequest) (dto.DbPrivilegeResponse, error) {

//		response, err := dbRepository.DbPrivilege(req)
//		if err != nil {
//			return dto.DbPrivilegeResponse{
//				Code:    500,
//				Status:  "Internal Server Error",
//				Message: "Error while inserting user privilege",
//			}, nil
//		}
//		return response, nil
//	}
func (d *DbService) ListDatabase(req dto.ListDbRequest) (dto.ListDbResponse, error) {

	response, err := dbRepository.ListDatabase(req)
	if err != nil {
		return dto.ListDbResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while displaying Database",
		}, nil
	}
	return response, nil
}

func (d *DbService) ListUser(req dto.ListUserRequest) (dto.ListUserResponse, error) {

	response, err := dbRepository.ListUser(req)
	if err != nil {
		return dto.ListUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while displaying User",
		}, nil
	}
	return response, nil
}
func (d *DbService) ListUserPrivilege(req dto.ListUserPrivilegeRequest) (dto.ListUserPrivilegeResponse, error) {

	response, err := dbRepository.ListUserPrivilege(req)
	if err != nil {
		return dto.ListUserPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while displaying User and privilege",
		}, nil
	}
	return response, nil
}
