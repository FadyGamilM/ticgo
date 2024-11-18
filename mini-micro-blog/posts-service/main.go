package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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
	var event_bus_svc_port = os.Getenv("EVENT_BUS_SVC_PORT")
	posts := []Post{}
	idCounter := 0

	r := gin.Default()

	r.GET("/posts", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, posts)
	})

	r.POST("/posts", func(ctx *gin.Context) {
		postBOdy := &Post{}

		err := ctx.BindJSON(postBOdy)
		if err != nil {
			fmt.Println("error_binding_request_body")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		posts = append(posts, Post{
			Id:   idCounter,
			Body: postBOdy.Body,
		})

		idCounter++

		event := Event{
			Type: EventType(post_created_event),
			Body: PostCreatedEventBody{
				Id:   posts[len(posts)-1].Id,
				Body: posts[len(posts)-1].Body,
			},
		}
		jsonData, err := json.Marshal(&event)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error_marshling_event_%s", err.Error())})
			return
		}
		http.Post("http://eventbus-srv:"+event_bus_svc_port+"/events", "application/json", bytes.NewBuffer(jsonData))

		ctx.JSON(http.StatusCreated, gin.H{
			"postId": posts[len(posts)-1].Id,
		})
	})

	r.POST("/events", func(ctx *gin.Context) {
		eventBody := &Event{}
		_ = ctx.Bind(eventBody)
		log.Println(eventBody)
		fmt.Println("consumed_event_", eventBody.Type, "_body_", eventBody.Body)
		ctx.JSON(http.StatusAccepted, gin.H{})
	})

	fmt.Println("server_listening_on_port", os.Getenv("PORT"))
	err := r.Run("0.0.0.0:" + os.Getenv("PORT"))
	if err != nil {
		log.Panicln("error_Starting_posts_Server_", err)
	}
}
