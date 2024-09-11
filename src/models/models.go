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
