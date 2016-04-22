package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type Feedback struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var (
	GClient       *github.Client
	tc            *http.Client
	latestRelease string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file does not exist")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc = oauth2.NewClient(oauth2.NoContext, ts)
	GClient = github.NewClient(tc)

	latestRelease = os.Getenv("LATEST_RELEASE")
}

func sendIssue(title string, body string) error {
	i := &github.IssueRequest{
		Title: &title,
		Body:  &body,
	}
	_, _, err := GClient.Issues.Create("IBM-Bluemix", "bluemix-code", i)
	return err
}

func handleIndex(c *gin.Context) {
	c.String(200, "Nothing to see here")
}

func handleUpdate(c *gin.Context) {
	os := c.Param("os")
	quality := c.Param("quality")
	commitID := c.Param("commitID")

	if os == "darwin" && quality == "stable" && commitID != latestRelease {
		c.JSON(200, gin.H{
			"url":     "https://ibm.biz/bluemixcode",
			"version": latestRelease,
		})
	} else {
		c.JSON(200, gin.H{"message": "Up to date"})
	}
}

func httpSendFeedback(c *gin.Context) {
	var feedback Feedback
	if c.BindJSON(&feedback) == nil {
		err := sendIssue(feedback.Title, feedback.Body)
		if err != nil {
			fmt.Println(err)
			c.String(400, "Unable to create feedback")
			return
		} else {
			c.String(200, "Thanks For the Feedback")
		}
	} else {
		c.String(400, "Invalid JSON body")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	//fix for gin not serving HEAD
	router.HEAD("/", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/", handleIndex)
	router.POST("/api/feedback", httpSendFeedback)
	router.GET("/api/update/:os/:quality/:commitID", handleUpdate)

	router.Run(":" + port)
	fmt.Println("Server started on", port)
}
