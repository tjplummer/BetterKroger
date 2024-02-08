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
		v1.GET(GetItems)
		v1.GET("/:id", GetItemById)
		v1.POST("item", AddItem)
		v1.PUT("/:id/:quantity", UpdateQuantity)
		v1.DELETE("/:id", RemoveItem)
	}

	r.Run()
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

func GetItemById(cx *gin.Context) {
	id := cx.Param("id")

	item, err := models.GetItemById(id)

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

func UpdateQuantity(cx *gin.Context) {
	id := cx.Param("id")
	quantity := cx.Param("quantity")

	success, err := models.UpdateQuantity(id, quantity)

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
