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
	TableName    string `json:"tableName"`
	UserName     string `json:"userName"`
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
