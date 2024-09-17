package models

type Organization struct {
	ID                   int64  `json:"id"`
	OrganizationName     string `json:"organizationName"`
	AdminEmail           string `json:"adminEmail"`
	SiteURL              string `json:"siteURL"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
	Status               string `json:"status"`
	AuthenticationMethod string `json:"authenticationMethod"`
	DatabaseStatus       string `json:"databaseStatus"`
}

func (o *Organization) TableName() string {
	return "did.organizations"
}

type DbSynchronization struct {
	ID       int `gorm:"column:id" json:"id"`
	OrgId    int `gorm:"column:org_id" json:"org_id"`
	TenantId int `gorm:"column:tenant_id" json:"tenant_id"`

	DatabaseType string `gorm:"column:db_type" json:"db_type"`
	DatabaseName string `gorm:"column:db_name" json:"db_name"`
	//Table        string `json:"table_name"`
	Host        string `gorm:"column:host" json:"host"`
	Port        string `gorm:"column:port" json:"port"`
	AgentStatus string `gorm:"column:agent_status" json:"agent_status"`
	Uuid        string `gorm:"column:uuid" json:"uuid"`
	Status      string `gorm:"column:status" json:"status"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (o *DbSynchronization) TableName() string {
	return "did.db_synchronization"
}

type DbUser struct {
	ID       int `gorm:"column:id" json:"id"`
	OrgId    int `gorm:"column:org_id" json:"org_id"`
	TenantId int `gorm:"column:tenant_id" json:"tenant_id"`

	DatabaseId int    `gorm:"column:db_id" json:"db_id"`
	UserName   string `gorm:"column:user_name" json:"username"`
	Role       string `gorm:"column:role" json:"role"`
	Host       string `gorm:"column:host" json:"host"`
	Status     string `gorm:"column:status" json:"status"`
	CreatedAt  int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (o *DbUser) TableName() string {
	return "did.db_user"
}

type DbPrivilege struct {
	ID         int    `gorm:"column:id" json:"id"`
	OrgId      int    `gorm:"column:org_id" json:"org_id"`
	TenantId   int    `gorm:"column:tenant_id" json:"tenant_id"`
	UserId     int    `gorm:"column:user_id" json:"user_id"`
	DatabaseId int    `gorm:"column:db_id" json:"db_id"`
	Privilege  string `gorm:"column:privilege" json:"privilege"`
	CreatedAt  int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (o *DbPrivilege) TableName() string {
	return "did.db_privilege"
}

type DbUserPrivilege struct {
	OrgId        int    `json:"org_id"`
	TenantId     int    `json:"tenant_id"`
	DatabaseName string `json:"db_name"`
	UserName     string `json:"user_name"`
	Host         string `json:"host"`
	Status       string `json:"status"`
	Privilege    string `json:"privilege"`
	CreatedAt    int64  `json:"created_at"`
}
