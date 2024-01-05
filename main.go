package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

type config struct {
	ADDR   string `env:"ADDR" envDefault:":8088"`
	ApiKey string `env:"API_KEY"`
}

var (
	cfg          config
	genaiContent context.Context
	genaiClient  *genai.Client
	genaiErr     error
	ginErr       error
)

type Talk struct {
	Message string `form:"message"`
}

func main() {
	cfg = config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	r := gin.Default()
	genaiContent = context.Background()
	genaiClient, genaiErr = genai.NewClient(genaiContent, option.WithAPIKey(cfg.ApiKey))
	if genaiErr != nil {
		log.Fatal(genaiErr)
	}
	defer func(genaiClient *genai.Client) {
		err := genaiClient.Close()
		if err != nil {

		}
	}(genaiClient)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/text", func(c *gin.Context) {
		var talk Talk
		if c.ShouldBind(&talk) == nil {
			fmt.Println(talk.Message)
			message, err := chat(talk.Message)
			if err != nil {
				data := map[string]interface{}{
					"code":    500,
					"message": err.Error(),
				}
				c.AsciiJSON(http.StatusOK, data)
			} else {
				data := map[string]interface{}{
					"code":    200,
					"message": message,
				}
				c.AsciiJSON(http.StatusOK, data)
			}
		} else {
			data := map[string]interface{}{
				"code":    500,
				"message": "欢迎使用AI CHAT",
			}
			c.AsciiJSON(http.StatusOK, data)
			c.AsciiJSON(http.StatusOK, data)
		}
	})

	ginErr = r.Run(cfg.ADDR)
	if ginErr != nil {
		fmt.Println(ginErr.Error())
		return
	}
}

func chat(text string) (string, error) {
	model := genaiClient.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(genaiContent, genai.Text(text))
	if err != nil {
		return "", err
	}
	var result string = ""
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				result = result + fmt.Sprintf("%v", part)
			}
		}
	}
	return result, nil
}
