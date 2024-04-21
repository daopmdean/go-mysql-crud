package main

import (
	"net/http"
	"strconv"

	"github.com/daopmdean/go-mysql-crud/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(model.Book{})

	r := gin.Default()

	r.GET("/books", func(ctx *gin.Context) {
		var books []*model.Book
		db.Find(&books)
		ctx.JSON(http.StatusOK, gin.H{
			"books": books,
		})
	})

	r.GET("/books/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"error": "invalid id",
			})
			return
		}

		var book model.Book
		db.First(&book, id)
		ctx.JSON(http.StatusOK, gin.H{
			"book": book,
		})
	})

	r.POST("/books", func(ctx *gin.Context) {
		book := model.Book{}
		if err := ctx.ShouldBindJSON(&book); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed binding json",
			})
			return
		}

		if book.Name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "book name required",
			})
			return
		}

		tx := db.Create(&book)
		if tx.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed create book",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":     book.ID,
			"create": "ok",
		})
	})

	r.Run()
}
