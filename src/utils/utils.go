package utils

import (
	"github.com/authnull0/database-service/src/db"
	"github.com/authnull0/database-service/src/models"
	"github.com/spf13/viper"
)

func GetOrganizationDatabaseName(orgid int) (string, error) {

	var organization models.Organization

	// Get the organization from the database

	db := db.GetConnectiontoDatabaseDynamically(viper.GetString(viper.GetString("env") + ".db.name"))
	err := db.Where("id = ?", orgid).First(&organization).Error
	if err != nil {
		return "", err
	}

	return organization.OrganizationName, nil

}
