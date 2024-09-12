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
	ID       int `json:"id"`
	OrgId    int `json:"org_id"`
	TenantId int `json:"tenant_id"`

	DatabaseType string `json:"database_type"`
	DatabaseName string `json:"database_name"`
	Table        string `json:"table_name"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Status       string `json:"status"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (o *DbSynchronization) TableName() string {
	return "db_synchronization"
}

type DbUser struct {
	ID       int `json:"id"`
	OrgId    int `json:"org_id"`
	TenantId int `json:"tenant_id"`

	TableId   int    `json:"table_id"`
	UserName  string `json:"username"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (o *DbUser) TableName() string {
	return "db_user"
}

type DbPrivilege struct {
	ID        int    `json:"id"`
	OrgId     int    `json:"org_id"`
	TenantId  int    `json:"tenant_id"`
	UserId    int    `json:"user_id"`
	TableId   int    `json:"table_id"`
	Privilege string `json:"privilege"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (o *DbPrivilege) TableName() string {
	return "db_user"
}
