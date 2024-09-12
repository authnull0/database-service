package dto

type DbSyncRequest struct {
	OrgID        int    `json:"orgId"`
	TenantID     int    `json:"tenantId"`
	Databasetype string `json:"databaseType"`
	DatabaseName string `json:"databaseName"`
	TableName    string `json:"tableName"`
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
