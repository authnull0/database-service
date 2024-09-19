package repository

import (
	"log"
	"strings"
	"time"

	"github.com/authnull0/database-service/src/db"
	"github.com/authnull0/database-service/src/models"
	"github.com/authnull0/database-service/src/models/dto"
	"github.com/authnull0/database-service/src/utils"
	"gorm.io/gorm"
)

type DbRepository struct{}

// Dbsync inserts the database sync details into the db_synchronization table
func (s *DbRepository) DbSync(req dto.DbSyncRequest) (dto.DbSyncResponse, error) {

	// Get the organization-specific database name dynamically
	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		log.Default().Println("Error while fetching organization database name:", err)
		return dto.DbSyncResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	// Establish a dynamic database connection
	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)
	var dbSync models.DbSynchronization

	err = orgDb.Table("did.db_synchronization").Where("db_name = ?", req.DatabaseName).First(&dbSync).Error
	if err == nil {
		// Database exists, update the status if it has changed
		if dbSync.Status != req.Status {
			dbSync.Status = req.Status
			dbSync.CreatedAt = time.Now().Unix()

			if err := orgDb.Table("did.db_synchronization").Save(&dbSync).Error; err != nil {
				log.Default().Println("Error while updating db_synchronization table:", err)
				return dto.DbSyncResponse{
					Code:    500,
					Status:  "Internal Server Error",
					Message: "Error while updating db_synchronization table",
				}, err
			}

			log.Default().Println("Database status updated successfully")
			return dto.DbSyncResponse{
				Code:    200,
				Status:  "Success",
				Message: "Database status updated successfully",
			}, nil
		}

		// No change in role, no update needed
		log.Default().Println("Database already exists with the same status. No update needed.")
		return dto.DbSyncResponse{
			Code:    200,
			Status:  "Success",
			Message: "Database already exists with the same status. No update needed.",
		}, nil
	}

	// Create the db_synchronization record to insert
	dbSync = models.DbSynchronization{
		OrgId:        req.OrgID,
		TenantId:     req.TenantID,
		DatabaseType: req.Databasetype,
		DatabaseName: req.DatabaseName,
		Host:         req.Host,
		Port:         req.Port,
		Status:       req.Status,
		Uuid:         req.Uuid,
		CreatedAt:    time.Now().Unix(),
	}

	// Insert the record into the db_synchronization table
	if err := orgDb.Table("did.db_synchronization").Create(&dbSync).Error; err != nil {
		log.Default().Println("Error while inserting into db_synchronization:", err)
		return dto.DbSyncResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while inserting into db_synchronization",
		}, err
	}

	// Return success response
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

	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	// Find db_synchronization entry
	var dbSync models.DbSynchronization
	err = orgDb.Table("did.db_synchronization").Where("org_id = ? AND tenant_id = ? AND db_name = ?",
		req.OrgID, req.TenantID, req.DatabaseName).First(&dbSync).Error
	if err != nil {
		log.Default().Println("Error while fetching from db_synchronization:", err)
		return dto.DbUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching from db_synchronization",
		}, err
	}

	cleanUserName := strings.Trim(req.UserName, "'")

	// Check if the user already exists in db_user
	var dbUser models.DbUser
	err = orgDb.Table("did.db_user").Where("org_id = ? AND tenant_id = ? AND user_name = ? AND db_id = ?",
		req.OrgID, req.TenantID, cleanUserName, dbSync.ID).First(&dbUser).Error

	if err == gorm.ErrRecordNotFound {
		// User doesn't exist, insert new user record
		dbUser = models.DbUser{
			OrgId:      req.OrgID,
			TenantId:   req.TenantID,
			DatabaseId: dbSync.ID,
			UserName:   cleanUserName,
			Role:       req.Role,
			Host:       req.Host,
			Status:     dbSync.Status,
			CreatedAt:  time.Now().Unix(),
		}

		if err := orgDb.Table("did.db_user").Create(&dbUser).Error; err != nil {
			log.Default().Println("Error while inserting into db_user:", err)
			return dto.DbUserResponse{
				Code:    500,
				Status:  "Internal Server Error",
				Message: "Error while inserting into db_user",
			}, err
		}

		log.Default().Printf("User %s inserted successfully into db_user\n", cleanUserName)

	} else if err == nil {
		// User exists, check if the role needs to be updated
		if dbUser.Role != req.Role {
			log.Default().Printf("Updating role for user %s from %s to %s\n", cleanUserName, dbUser.Role, req.Role)
			dbUser.Role = req.Role

			if err := orgDb.Table("did.db_user").Save(&dbUser).Error; err != nil {
				log.Default().Println("Error while updating role in db_user:", err)
				return dto.DbUserResponse{
					Code:    500,
					Status:  "Internal Server Error",
					Message: "Error while updating role in db_user",
				}, err
			}

			log.Default().Printf("Role for user %s updated successfully to %s\n", cleanUserName, req.Role)
		} else {
			log.Default().Printf("User %s already exists with the same role %s. No update needed.\n", cleanUserName, dbUser.Role)
		}

	} else {

		log.Default().Println("Error while checking user in db_user:", err)
		return dto.DbUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while checking user in db_user",
		}, err
	}

	// Split the privileges string by commas
	privileges := strings.Split(req.Privilege, ", ")

	for _, privilege := range privileges {
		cleanPrivilege := strings.Trim(privilege, `"`)

		// Check if the privilege already exists for this user
		var dbPrivilege models.DbPrivilege
		err = orgDb.Table("did.db_privilege").Where("user_id = ? AND db_id = ? AND privilege = ?",
			dbUser.ID, dbSync.ID, cleanPrivilege).First(&dbPrivilege).Error

		if err == nil {
			// Privilege already exists, skip this privilege
			log.Default().Printf("Privilege %s for user %s already exists. Skipping...\n", cleanPrivilege, cleanUserName)
			continue
		}

		// Insert a new privilege record if it doesn't exist
		dbPrivilege = models.DbPrivilege{
			OrgId:      req.OrgID,
			TenantId:   req.TenantID,
			DatabaseId: dbSync.ID,
			UserId:     dbUser.ID,
			Privilege:  cleanPrivilege,
			CreatedAt:  time.Now().Unix(),
		}

		if err := orgDb.Table("did.db_privilege").Create(&dbPrivilege).Error; err != nil {
			log.Default().Println("Error while inserting into db_privilege:", err)
			return dto.DbUserResponse{
				Code:    500,
				Status:  "Internal Server Error",
				Message: "Error while inserting into db_privilege",
			}, err
		}

		log.Default().Printf("Privilege %s for user %s inserted successfully into db_privilege\n", cleanPrivilege, cleanUserName)
	}

	return dto.DbUserResponse{
		Code:    200,
		Status:  "Success",
		Message: "User and privileges processed successfully",
	}, nil
}

func (s *DbRepository) ListDatabase(req dto.ListDbRequest) (dto.ListDbResponse, error) {

	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		return dto.ListDbResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	var listDbSync []models.DbSynchronization
	query := orgDb.Table("did.db_synchronization").Where("org_id = ? AND tenant_id = ?", req.OrgID, req.TenantID).Find(&listDbSync)

	for _, filter := range req.Filters {
		if filter.FilterType == "Database" {
			query = query.Where("db_name = ?", filter.FilterValue)
		}
	}

	for _, filter := range req.Filters {
		if filter.FilterType == "Status" {
			query = query.Where("status = ?", filter.FilterValue)
		}
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		log.Printf("%s", err)
	}

	offset := (req.PageId - 1) * req.Limit
	totalPages := (totalCount + int64(req.Limit) - 1) / int64(req.Limit)

	// Fetch logs with limit and offset
	if err := query.Offset(offset).Limit(req.Limit).Find(&listDbSync).Error; err != nil {
		log.Default().Println("Error while fetching list from db_synchronization:", err)
		return dto.ListDbResponse{
			Code:      500,
			Status:    "Internal Server Error",
			Message:   "Error while listing databases",
			RequestId: req.RequestId,
			Limit:     req.Limit,
			PageId:    req.PageId,
		}, err
	}

	log.Default().Printf("Total count %d", totalCount)

	return dto.ListDbResponse{
		Code:       200,
		Status:     "Success",
		Message:    "Database list fetched successfully",
		Data:       listDbSync,
		RequestId:  req.RequestId,
		Limit:      req.Limit,
		PageId:     req.PageId,
		TotalPages: int(totalPages),
		TotalCount: totalCount,
	}, nil
}

func (s *DbRepository) ListUser(req dto.ListUserRequest) (dto.ListUserResponse, error) {

	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		return dto.ListUserResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	var listUser []models.DbUser

	query := orgDb.Table("did.db_user AS u").
		Select("u.*, ds.db_name").
		Joins("JOIN did.db_synchronization AS ds ON u.db_id = ds.id").
		Where("u.org_id = ? AND u.tenant_id = ?", req.OrgID, req.TenantID)

	for _, filter := range req.Filters {
		if filter.FilterType == "Database" {
			query = query.Where("ds.db_name = ?", filter.FilterValue)
		}
	}

	for _, filter := range req.Filters {
		if filter.FilterType == "User" {
			query = query.Where("u.user_name = ?", filter.FilterValue)
		}
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		log.Printf("%s", err)
	}

	offset := (req.PageId - 1) * req.Limit
	totalPages := (totalCount + int64(req.Limit) - 1) / int64(req.Limit)

	// Fetch logs with limit and offset
	if err := query.Offset(offset).Limit(req.Limit).Find(&listUser).Error; err != nil {
		log.Default().Println("Error while fetching list from db_synchronization:", err)
		return dto.ListUserResponse{
			Code:      500,
			Status:    "Internal Server Error",
			Message:   "Error while listing user",
			RequestId: req.RequestId,
			Limit:     req.Limit,
			PageId:    req.PageId,
		}, err
	}

	log.Default().Printf("Total count %d", totalCount)

	return dto.ListUserResponse{
		Code:       200,
		Status:     "Success",
		Message:    "Database list fetched successfully",
		Data:       listUser,
		RequestId:  req.RequestId,
		Limit:      req.Limit,
		PageId:     req.PageId,
		TotalPages: int(totalPages),
		TotalCount: totalCount,
	}, nil
}

func (s *DbRepository) ListUserPrivilege(req dto.ListUserPrivilegeRequest) (dto.ListUserPrivilegeResponse, error) {

	dbName, err := utils.GetOrganizationDatabaseName(req.OrgID)
	if err != nil {
		return dto.ListUserPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching organization",
		}, err
	}
	log.Default().Println("DB Name: ", dbName)

	orgDb := db.GetConnectiontoDatabaseDynamically(dbName)

	var dbUsers []models.DbUser
	err = orgDb.Where("org_id = ? AND tenant_id = ?", req.OrgID, req.TenantID).Find(&dbUsers).Error
	if err != nil {
		return dto.ListUserPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching users",
		}, err
	}

	var privileges []models.DbPrivilege
	for _, user := range dbUsers {
		var userPrivileges []models.DbPrivilege
		err = orgDb.Where("user_id = ?", user.ID).Find(&userPrivileges).Error
		if err != nil {
			log.Printf("Error while fetching privileges for user ID %d: %s", user.ID, err)
			continue
		}
		privileges = append(privileges, userPrivileges...)
	}

	var dbSyncs []models.DbSynchronization
	for _, privilege := range privileges {
		var sync models.DbSynchronization
		err = orgDb.Where("id = ?", privilege.DatabaseId).First(&sync).Error
		if err != nil {
			log.Printf("Error while fetching sync info for DB ID %d: %s", privilege.DatabaseId, err)
			continue
		}
		dbSyncs = append(dbSyncs, sync)
	}

	var listUserPrivilege []dto.DbUserPrivilegeResponse
	for _, privilege := range privileges {
		// Find the associated user and synchronization info for each privilege
		var user models.DbUser
		for _, u := range dbUsers {
			if u.ID == privilege.UserId {
				user = u
				break
			}
		}

		var sync models.DbSynchronization
		for _, s := range dbSyncs {
			if s.ID == privilege.DatabaseId {
				sync = s
				break
			}
		}

		// Combine the data into a single response struct
		listUserPrivilege = append(listUserPrivilege, dto.DbUserPrivilegeResponse{
			ID:        user.ID,
			OrgID:     sync.OrgId,
			TenantID:  sync.TenantId,
			DbName:    sync.DatabaseName,
			UserName:  user.UserName,
			Host:      user.Host,
			Status:    sync.Status,
			Role:      user.Role,
			Privilege: privilege.Privilege,
			CreatedAt: sync.CreatedAt,
		})
	}

	// Apply filters after fetching the results
	filteredResults := listUserPrivilege
	for _, filter := range req.Filters {
		var temp []dto.DbUserPrivilegeResponse
		for _, item := range filteredResults {
			switch filter.FilterType {
			case "Database":
				if item.DbName == filter.FilterValue {
					temp = append(temp, item)
				}
			case "User":
				if item.UserName == filter.FilterValue {
					temp = append(temp, item)
				}
			case "Privilege":
				if item.Privilege == filter.FilterValue {
					temp = append(temp, item)
				}
			case "Status":
				if item.Status == filter.FilterValue {
					temp = append(temp, item)
				}
			}
		}
		filteredResults = temp
	}

	totalCount := len(filteredResults)
	totalPages := (totalCount + req.Limit - 1) / req.Limit // Total pages calculation
	offset := (req.PageId - 1) * req.Limit

	// Handle out of bounds for pagination
	if offset >= totalCount {
		return dto.ListUserPrivilegeResponse{
			Code:       200,
			Status:     "Success",
			Message:    "Users and privileges list fetched successfully",
			Data:       []dto.DbUserPrivilegeResponse{}, // Return empty if no results for that page
			RequestId:  req.RequestId,
			Limit:      req.Limit,
			PageId:     req.PageId,
			TotalPages: totalPages,
			TotalCount: int64(totalCount),
		}, nil
	}

	// Get the paginated result
	end := offset + req.Limit
	if end > totalCount {
		end = totalCount
	}
	paginatedResult := filteredResults[offset:end]

	return dto.ListUserPrivilegeResponse{
		Code:       200,
		Status:     "Success",
		Message:    "Users and privileges list fetched successfully",
		Data:       paginatedResult,
		RequestId:  req.RequestId,
		Limit:      req.Limit,
		PageId:     req.PageId,
		TotalPages: totalPages,
		TotalCount: int64(totalCount),
	}, nil
}
