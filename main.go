package main

import (
	"BetterKroger/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
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
}

func Heartbeat(cx *gin.Context) {
	cx.JSON(http.StatusOK, gin.H{"status": "Active"})
}

func GetItems(cx *gin.Context) {
	items, err := models.GetItems()

	Verify(err)

	if items == nil {
		cx.JSON(http.StatusBadRequest, gin.H{"error": "No Items Found"})
		return
	} else {
		cx.JSON(http.StatusOK, gin.H{"data": items})
	}
}

func GetItemByCode(cx *gin.Context) {
	code := cx.Param("code")

	item, err := models.GetItemByCode(code)

	Verify(err)

	if len(item.Name) == 0 {
		cx.JSON(http.StatusBadRequest, gin.H{"error": "No Item Found By Requested Id"})
		return
	} else {
		cx.JSON(http.StatusOK, gin.H{"data": item})
	}
}

func AddItem(cx *gin.Context) {
	var json models.Item

	if err := cx.ShouldBindJSON(&json); err != nil {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddItem(json)

	if success {
		cx.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func UpdatePrice(cx *gin.Context) {
	code := cx.Param("code")
	quantity := cx.Param("quantity")

	success, err := models.UpdatePrice(code, quantity)

	if success {
		cx.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func RemoveItem(cx *gin.Context) {
	id, err := strconv.Atoi(cx.Param("id"))

	if err != nil {
		cx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.RemoveItem(id)

	if success {
		cx.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		cx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func Verify(err error) {
	// log
}
