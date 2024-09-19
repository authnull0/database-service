package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/authnull0/database-service/src/db"
	"github.com/authnull0/database-service/src/models"
	"github.com/authnull0/database-service/src/models/dto"
	"github.com/gin-gonic/gin"
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
func Authnz(token string, ctx *gin.Context, orgid int) bool {
	client := &http.Client{}
	//set token in header Authorization

	env := viper.GetString("env")

	authnzurl := viper.GetString(env + ".url.authnzurl")

	req, err := http.NewRequest("GET", authnzurl, nil)
	if err != nil {
		log.Default().Println(err)
	}
	req.Header.Set("X-Authorization", token)
	req.Header.Set("X-Request-Path", ctx.Request.URL.Path)
	req.Header.Set("X-Org-Id", strconv.Itoa(orgid))

	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println(err)
		return false
	}

	// Read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Default().Println(err)
		return false
	}

	// Print the response body as a string
	fmt.Println(string(body))

	var response dto.AuthnzResponseDTO

	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Default().Println(err)
		return false
	}

	return response.Validation
}
