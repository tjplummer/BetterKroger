package main

import (
	"BetterKroger/models"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	zap.L().Info("Logger Started")

	Verify(models.StartDB())

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("", Heartbeat)
		v1.GET("", GetItems)
		v1.GET("/:code", GetItemByCode)
		v1.POST("item", AddItem)
		v1.PUT("/:code/:price", UpdatePrice)
		v1.DELETE("/:code", RemoveItem)
	}

	r.Run()

	zap.L().Info("Better Kroger initialized")
}

func Heartbeat(cx *gin.Context) {
	cx.JSON(http.StatusOK, gin.H{"status": "Active"})
	zap.L().Info("Heartbeat check")
}

func GetItems(cx *gin.Context) {
	items, err := models.GetItems()

	if items == nil {
		cx.JSON(http.StatusBadRequest, gin.H{"error": "No Items Found"})
		zap.L().Error(fmt.Sprintf("GetItems: Error - %x", err))
	} else {
		cx.JSON(http.StatusOK, gin.H{"data": items})
		zap.L().Info("GetItems: Okay")
	}
}

func GetItemByCode(cx *gin.Context) {
	code := cx.Param("code")

	item, err := models.GetItemByCode(code)

	if len(item.Name) == 0 {
		cx.JSON(http.StatusBadRequest, gin.H{"error": "No Item Found By Requested Id"})
		zap.L().Error(fmt.Sprintf("GetItemsByCode: Error - %x", err))
	} else {
		cx.JSON(http.StatusOK, gin.H{"data": item})
		zap.L().Info("GetItemsByCode: Okay")
	}
}

func AddItem(cx *gin.Context) {
	var json models.Item

	if err := cx.ShouldBindJSON(&json); err != nil {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		zap.L().Error(fmt.Sprintf("AddItem: Error - %x", err))
	}

	success, err := models.AddItem(json)

	if success {
		cx.JSON(http.StatusOK, gin.H{"message": "Success"})
		zap.L().Info("AddItem: Okay")
	} else {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err})
		zap.L().Error(fmt.Sprintf("AddItem: Error - %x", err))
	}
}

func UpdatePrice(cx *gin.Context) {
	code := cx.Param("code")
	quantity := cx.Param("quantity")

	success, err := models.UpdatePrice(code, quantity)

	if success {
		cx.JSON(http.StatusOK, gin.H{"message": "Success"})
		zap.L().Info("UpdatePrice: Okay")
	} else {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err})
		zap.L().Error(fmt.Sprintf("UpdatePrice: Error - %x", err))
	}
}

func RemoveItem(cx *gin.Context) {
	id, err := strconv.Atoi(cx.Param("id"))

	if err != nil {
		cx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		zap.L().Error(fmt.Sprintf("RemoveItem: Error - %x", err))
	}

	success, err := models.RemoveItem(id)

	if success {
		cx.JSON(http.StatusOK, gin.H{"message": "Success"})
		zap.L().Info("RemoveItem: Okay")
	} else {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err})
		zap.L().Error(fmt.Sprintf("RemoveItem: Error - %x", err))
	}
}

func Verify(err error) {
	// log
}
