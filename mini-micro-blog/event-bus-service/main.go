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

	r.POST("/event/", func(ctx *gin.Context) {
		event := Event{}
		err := ctx.BindJSON(&event)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed_to_bind_event_body_" + err.Error(),
			})
			return
		}

		switch event.Type {
		case EventType(post_created_event):
			log.Println("post_has_been_created")
			jsonData, err := json.Marshal(event.Body.(PostCreatedEventBody))
			if err != nil {
				log.Println("failed_to_marshal_event_body")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			res, err := http.Post(fmt.Sprintf("http://localhost:%s/events", posts_svc_port), "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusAccepted {
				log.Println("failure_in_receiving request")
				ctx.JSON(res.StatusCode, gin.H{
					"error": "request_not_accepted_by_posts_service",
				})
				return
			}

		case EventType(comment_created_event):
			log.Println("comment_has_been_created")
			jsonData, err := json.Marshal(event.Body.(CommentCreatedEventBody))
			if err != nil {
				log.Println("failed_to_marshal_event_body")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			res, err := http.Post(fmt.Sprintf("http://localhost:%s/events", comments_svc_port), "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusAccepted {
				log.Println("failure_in_receiving request")
				ctx.JSON(res.StatusCode, gin.H{
					"error": "request_not_accepted_by_comments_service",
				})
				return
			}
		}

		// we decode the event, and we send it to all subscribed consumers (currently all svcs) and the svc who is interested in should consume it and process it
	})

	// if err := r.Run("0.0.0.0:", os.Getenv("PORT")); err != nil {
	// 	log.Println("error_starting_event_bus_server_", err.Error())
	// }
	r.Run("0.0.0.0:" + os.Getenv("PORT"))

}
