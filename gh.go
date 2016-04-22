package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"os"
)

var (
	gclient       *github.Client
	latestRelease string
	org string = "IBM-Bluemix"
	repo string = "bluemix-code"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file does not exist")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	gclient = github.NewClient(tc)

	latestRelease = os.Getenv("LATEST_RELEASE")
}

func sendIssue(issue *github.IssueRequest) error {
	_, _, err := gclient.Issues.Create(org, repo, issue)
	return err
}

func handleIndex(c *gin.Context) {
	c.String(200, "Nothing to see here")
}

func handleUpdate(c *gin.Context) {
	os := c.Param("os")
	quality := c.Param("quality")
	commitID := c.Param("commit_id")

	if os == "darwin" && quality == "stable" && commitID != latestRelease {
		c.JSON(200, gin.H{
			"url":          "https://ibm.biz/bluemixcode",
			"version":      latestRelease,
			"releaseNotes": "https://ibm.biz/bluemixcode-releasenotes",
		})
	} else {
		c.JSON(200, gin.H{"message": "Up to date"})
	}
}

func handleFeedback(c *gin.Context) {
	labels := []string{"feedback", "from_ide"}
	issue := &github.IssueRequest{
		Labels: &labels,
	}
	if c.BindJSON(&issue) == nil {
		err := sendIssue(issue)
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
	router.POST("/api/feedback", handleFeedback)
	router.GET("/api/update/:os/:quality/:commit_id", handleUpdate)

	router.Run(":" + port)
}
