package repository

import (
	"log"
	"strings"
	"time"

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
		//Table:        req.TableName,
		Host:      req.Host,
		Port:      req.Port,
		Status:    req.Status,
		CreatedAt: time.Now().Unix(),
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

	dbUser := models.DbUser{
		OrgId:    req.OrgID,
		TenantId: req.TenantID,

		DatabaseId: dbSync.ID, // Using the ID from db_synchronization as database_id
		UserName:   cleanUserName,
		Role:       req.Role,
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

	return dto.DbUserResponse{
		Code:    200,
		Status:  "Success",
		Message: "User details inserted successfully into db_user",
	}, nil
}
func (s *DbRepository) DbPrivilege(req dto.DbPrivilegeRequest) (dto.DbPrivilegeResponse, error) {
	log.Default().Println("Inside the DbPrivilege - Database service")
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

	cleanUserName := strings.Trim(req.UserName, "'")

	// Step 1: Find  table_id
	var dbSync models.DbSynchronization
	err = orgDb.Table("did.db_synchronization").Where("org_id = ? AND tenant_id = ? AND db_name = ?",
		req.OrgID, req.TenantID, req.DatabaseName).First(&dbSync).Error
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
	err = orgDb.Table("did.db_user").Where("org_id = ? AND tenant_id = ? AND user_name = ?",
		req.OrgID, req.TenantID, cleanUserName).First(&dbUser).Error
	if err != nil {
		log.Default().Println("Error while fetching from db_user:", err)
		return dto.DbPrivilegeResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error while fetching from db_user",
		}, err
	}

	cleanPrivilege := strings.Trim(req.Privilege, `"`)
	userPrivilege := models.DbPrivilege{
		OrgId:      req.OrgID,
		TenantId:   req.TenantID,
		UserId:     dbUser.ID, // User ID from db_user
		DatabaseId: dbSync.ID, // Table ID from db_synchronization
		Privilege:  cleanPrivilege,
		CreatedAt:  time.Now().Unix(),
	}

	if err := orgDb.Table("did.db_privilege").Create(&userPrivilege).Error; err != nil {
		log.Default().Println("Error while inserting into db_privilege:", err)
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
			Host:      sync.Host,
			Status:    sync.Status,
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
