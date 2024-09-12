package repository

import (
	"log"

	"github.com/authnull0/database-service/src/db"
	"github.com/authnull0/database-service/src/models"
	"github.com/authnull0/database-service/src/models/dto"
	"github.com/authnull0/database-service/src/utils"
)

type DbRepository struct{}

// Dbsync inserts the database sync details into the db_synchronization table
func (s *DbRepository) DbSync(req dto.DbSyncRequest) (dto.DbSyncResponse, error) {
	// Get the organization-specific database name dynamically
	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		log.Default().Println("Error while fetching organization:", err)
		return dto.DbSyncResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	// Establish a dynamic database connection
	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	// Create the db_synchronization record to insert
	dbSync := models.DbSynchronization{
		OrgId:        req.OrgID,
		TenantId:     req.TenantID,
		DatabaseType: req.Databasetype,
		DatabaseName: req.DatabaseName,
		Table:        req.TableName,
		Status:       req.Status,
	}

	// Insert the record into the db_synchronization table
	if err := orgDb.Table("db_synchronization").Create(&dbSync).Error; err != nil {
		log.Default().Println("Error while inserting into db_synchronization:", err)
		return dto.DbSyncResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while inserting into db_synchronization",
		}, err
	}

	return dto.DbSyncResponse{
		Code:    200,
		Status:  "Success",
		Message: "Database synchronization details inserted successfully",
	}, nil
}

func (s *DbRepository) DbUser(req dto.DbUserRequest) (dto.DbUserResponse, error) {
	// Get the organization-specific database name dynamically
	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		log.Default().Println("Error while fetching organization:", err)
		return dto.DbUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	// Step 1: Find id in db_synchronization
	var dbSync models.DbSynchronization
	err = orgDb.Table("db_synchronization").Where("org_id = ? AND tenant_id = ? AND database_name = ? AND table = ?",
		req.OrgID, req.TenantID, req.DatabaseName, req.TableName).First(&dbSync).Error
	if err != nil {
		log.Default().Println("Error while fetching from db_synchronization:", err)
		return dto.DbUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching from db_synchronization",
		}, err
	}

	dbUser := models.DbUser{
		OrgId:    req.OrgID,
		TenantId: req.TenantID,

		TableId:  dbSync.ID, // Using the ID from db_synchronization as table_id
		UserName: req.UserName,
	}

	if err := orgDb.Table("db_user").Create(&dbUser).Error; err != nil {
		log.Default().Println("Error while inserting into db_user:", err)
		return dto.DbUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while inserting into db_user",
		}, err
	}

	return dto.DbUserResponse{
		Code:    200,
		Status:  "Success",
		Message: "User details inserted successfully into db_user",
	}, nil
}
func (s *DbRepository) DbPrivilege(req dto.DbPrivilegeRequest) (dto.DbPrivilegeResponse, error) {

	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		log.Default().Println("Error while fetching organization:", err)
		return dto.DbPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	// Step 1: Find  table_id
	var dbSync models.DbSynchronization
	err = orgDb.Table("db_synchronization").Where("org_id = ? AND tenant_id = ? AND database_name = ? AND table = ?",
		req.OrgID, req.TenantID, req.DatabaseName, req.TabelName).First(&dbSync).Error
	if err != nil {
		log.Default().Println("Error while fetching from db_synchronization:", err)
		return dto.DbPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching from db_synchronization",
		}, err
	}

	// Step 2: Find the relevant entry in db_user (for user_id)
	var dbUser models.DbUser
	err = orgDb.Table("db_user").Where("org_id = ? AND tenant_id = ? AND database_name = ? AND user = ?",
		req.OrgID, req.TenantID, req.DatabaseName, req.UserName).First(&dbUser).Error
	if err != nil {
		log.Default().Println("Error while fetching from db_user:", err)
		return dto.DbPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching from db_user",
		}, err
	}

	userPrivilege := models.DbPrivilege{
		OrgId:     req.OrgID,
		TenantId:  req.TenantID,
		UserId:    dbUser.ID, // User ID from db_user
		TableId:   dbSync.ID, // Table ID from db_synchronization
		Privilege: req.Privilege,
	}

	if err := orgDb.Table("user_privilege").Create(&userPrivilege).Error; err != nil {
		log.Default().Println("Error while inserting into user_privilege:", err)
		return dto.DbPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while inserting into user_privilege",
		}, err
	}

	return dto.DbPrivilegeResponse{
		Code:    200,
		Status:  "Success",
		Message: "Privilege details inserted successfully into user_privilege",
	}, nil
}
