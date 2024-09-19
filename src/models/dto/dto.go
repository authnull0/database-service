package dto

type DbSyncRequest struct {
	OrgID        int    `json:"orgId"`
	TenantID     int    `json:"tenantId"`
	Databasetype string `json:"databaseType"`
	DatabaseName string `json:"databaseName"`
	TableName    string `json:"tableName"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Status       string `json:"status"`
	Uuid         string `json:"uuid"`
}

type DbSyncResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type DbUserRequest struct {
	OrgID        int    `json:"orgId"`
	TenantID     int    `json:"tenantId"`
	Databasetype string `json:"databaseType"`
	DatabaseName string `json:"databaseName"`
	Status       string `json:"status"`
	UserName     string `json:"userName"`
	Host         string `json:"host"`
	Role         string `json:"role"`
}

type DbUserResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DbPrivilegeRequest struct {
	OrgID        int    `json:"orgId"`
	TenantID     int    `json:"tenantId"`
	Databasetype string `json:"databaseType"`
	DatabaseName string `json:"databaseName"`
	TabelName    string `json:"tableName"`
	UserName     string `json:"userName"`
	Host         string `json:"host"`
	Privilege    string `json:"privilegeType"`
}

type DbPrivilegeResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ListDbRequest struct {
	OrgID     int      `json:"orgId"`
	TenantID  int      `json:"tenantId"`
	Filters   []Filter `json:"filters"`
	RequestId string   `json:"requestId"`
	Limit     int      `json:"limit"`
	PageId    int      `json:"page_id"`
}

type ListDbResponse struct {
	Code       int         `json:"code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	RequestId  string      `json:"requestId"`
	Limit      int         `json:"limit"`
	PageId     int         `json:"page_id"`
	TotalPages int         `json:"total_pages"`
	TotalCount int64       `json:"total_count"`
}

type ListUserRequest struct {
	OrgID     int      `json:"orgId"`
	TenantID  int      `json:"tenantId"`
	Filters   []Filter `json:"filters"`
	RequestId string   `json:"requestId"`
	Limit     int      `json:"limit"`
	PageId    int      `json:"page_id"`
}

type ListUserResponse struct {
	Code       int         `json:"code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	DbName     string      `json:db_name`
	RequestId  string      `json:"requestId"`
	Limit      int         `json:"limit"`
	PageId     int         `json:"page_id"`
	TotalPages int         `json:"total_pages"`
	TotalCount int64       `json:"total_count"`
}
type ListUserPrivilegeRequest struct {
	OrgID     int      `json:"orgId"`
	TenantID  int      `json:"tenantId"`
	Filters   []Filter `json:"filters"`
	RequestId string   `json:"requestId"`
	Limit     int      `json:"limit"`
	PageId    int      `json:"page_id"`
}

type ListUserPrivilegeResponse struct {
	Code       int         `json:"code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	RequestId  string      `json:"requestId"`
	Limit      int         `json:"limit"`
	PageId     int         `json:"page_id"`
	TotalPages int         `json:"total_pages"`
	TotalCount int64       `json:"total_count"`
}

type Filter struct {
	FilterType  string `json:"filterParameter"`
	FilterValue string `json:"filterValue"`
}
type DbUserPrivilegeResponse struct {
	ID        int    `json:"id"`
	OrgID     int    `json:"org_id"`
	TenantID  int    `json:"tenant_id"`
	DbName    string `json:"db_name"`
	UserName  string `json:"user_name"`
	Host      string `json:"host"`
	Status    string `json:"status"`
	Role      string `json:"role"`
	Privilege string `json:"privilege"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
type AuthnzResponseDTO struct {
	Validation bool
	Code       int
	Message    string
	Status     string
	User       string
}
