package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Comment struct {
	Id     int    `json:"id"`
	PostId int    `json:"post_id"`
	Body   string `json:"body"`
}

var post_created_event = "postCreated"
var comment_created_event = "commentCreated"

type EventType string

type PostCreatedEventBody struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type CommentCreatedEventBody struct {
	Id     int    `json:"id"`
	PostId int    `json:"post_id"`
	Body   string `json:"body"`
}

type Event struct {
	Type EventType   `json:"event_type"`
	Body interface{} `json:"event_body"`
}

func main() {
	commentsDB := map[int][]Comment{}

	commentIdCounter := 0

	r := gin.Default()

	r.GET("/posts/:id/comments", func(ctx *gin.Context) {
		postId := ctx.Param("id")
		validPostId, err := strconv.Atoi(postId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		comments, ok := commentsDB[validPostId]
		if ok {
			ctx.JSON(http.StatusOK, comments)
			return
		} else {
			ctx.JSON(http.StatusOK, []Comment{})
			return
		}
	})

	r.POST("/posts/:id/comments", func(ctx *gin.Context) {
		postId := ctx.Param("id")
		validPostId, err := strconv.Atoi(postId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		commentBody := Comment{}
		if err := ctx.BindJSON(&commentBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		commentBody.Id = commentIdCounter
		commentIdCounter++
		commentBody.PostId = validPostId

		comments, ok := commentsDB[validPostId]
		if ok {
			comments = append(comments, commentBody)
			commentsDB[validPostId] = comments
			ctx.JSON(http.StatusCreated, gin.H{})
			return
		}

		commentsDB[validPostId] = []Comment{
			commentBody,
		}
		ctx.JSON(http.StatusCreated, gin.H{})
	})

	r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
