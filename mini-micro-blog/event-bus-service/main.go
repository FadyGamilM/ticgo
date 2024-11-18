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
	r := gin.Default()

	var posts_svc_port = os.Getenv("POSTS_SVC_PORT")
	var comments_svc_port = os.Getenv("COMMENTS_SVC_PORT")

	r.POST("/events", func(ctx *gin.Context) {
		event := Event{}
		err := ctx.BindJSON(&event)
		if err != nil {
			log.Println("failed to bind")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed_to_bind_event_body_" + err.Error(),
			})
			return
		}

		switch event.Type {
		case EventType(post_created_event):
			log.Println("post_has_been_created")
			// Convert the map to our struct
			bodyMap := event.Body.(map[string]interface{})
			postEvent := Event{
				Type: event.Type,
				Body: PostCreatedEventBody{
					Id:   int(bodyMap["id"].(int)),
					Body: bodyMap["body"].(string),
				},
			}

			jsonData, err := json.Marshal(postEvent)
			if err != nil {
				log.Println("failed_to_marshal_event_body:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			http.Post(fmt.Sprintf("http://posts-svc:%s/events", posts_svc_port), "application/json", bytes.NewBuffer(jsonData))

			http.Post(fmt.Sprintf("http://comments-svc:%s/events", comments_svc_port), "application/json", bytes.NewBuffer(jsonData))

		case EventType(comment_created_event):
			fmt.Println("comment_has_been_created")
			// Convert the map to our struct
			bodyMap := event.Body.(map[string]interface{})
			commentEvent := Event{
				Type: event.Type,
				Body: CommentCreatedEventBody{
					Id:     int(bodyMap["id"].(int)),
					Body:   bodyMap["body"].(string),
					PostId: bodyMap["post_id"].(int),
				},
			}

			jsonData, err := json.Marshal(commentEvent)
			if err != nil {
				log.Println("failed_to_marshal_event_body:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			http.Post(fmt.Sprintf("http://posts-svc:%s/events", posts_svc_port), "application/json", bytes.NewBuffer(jsonData))

			http.Post(fmt.Sprintf("http://comments-svc:%s/events", comments_svc_port), "application/json", bytes.NewBuffer(jsonData))
		}

	})

	r.Run("0.0.0.0:" + os.Getenv("PORT"))

}
