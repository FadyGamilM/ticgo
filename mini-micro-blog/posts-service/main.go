package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func main() {
	posts := []Post{}
	idCounter := 0

	r := gin.Default()

	r.GET("/posts", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, posts)
		return
	})

	r.POST("/posts", func(ctx *gin.Context) {
		postBOdy := &Post{}

		err := ctx.BindJSON(postBOdy)
		if err != nil {
			log.Println("error_binding_request_body")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		posts = append(posts, Post{
			Id:   idCounter,
			Body: postBOdy.Body,
		})

		ctx.JSON(http.StatusCreated, gin.H{
			"postId": posts[len(posts)-1].Id,
		})
		return
	})

	r.Run(":" + os.Getenv("PORT"))
}
