package controller

import (
	// "encoding/json"
	// "log"

	"github.com/authnull0/database-service/src/models/dto"
	"github.com/authnull0/database-service/src/service"
	"github.com/gin-gonic/gin"
)

type DbController struct{}

var dbService = service.DbService{}

func (s *DbController) DbSync(ctx *gin.Context) {
	var dbSyncRequest dto.DbSyncRequest
	if err := ctx.ShouldBindJSON(&dbSyncRequest); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := dbService.DbSync(dbSyncRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)

}
func (s *DbController) DbUser(ctx *gin.Context) {
	var dbUserRequest dto.DbUserRequest
	if err := ctx.ShouldBindJSON(&dbUserRequest); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := dbService.DbUser(dbUserRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)

}

// func (s *DbController) DbPrivilege(ctx *gin.Context) {
// 	var dbPrivilegeRequest dto.DbPrivilegeRequest
// 	if err := ctx.ShouldBindJSON(&dbPrivilegeRequest); err != nil {
// 		ctx.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	resp, err := dbService.DbPrivilege(dbPrivilegeRequest)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(200, resp)

// }

func (s *DbController) ListDatabase(ctx *gin.Context) {
	var listDbRequest dto.ListDbRequest
	if err := ctx.ShouldBindJSON(&listDbRequest); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := dbService.ListDatabase(listDbRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)

}
func (s *DbController) ListUser(ctx *gin.Context) {
	var listUserRequest dto.ListUserRequest
	if err := ctx.ShouldBindJSON(&listUserRequest); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := dbService.ListUser(listUserRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)

}

func (s *DbController) ListUserPrivilege(ctx *gin.Context) {
	var listUserPrivilegeRequest dto.ListUserPrivilegeRequest
	if err := ctx.ShouldBindJSON(&listUserPrivilegeRequest); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := dbService.ListUserPrivilege(listUserPrivilegeRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, resp)

}
