package main

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

var (
	defaultAddr  = ":8080"
	addr         string
	apiKey       string
	genaiContent context.Context
	genaiClient  *genai.Client
	genaiErr     error
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr = os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddr
	}
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Error genai api_key")
	}

	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)

	genaiContent = context.Background()
	genaiClient, genaiErr = genai.NewClient(genaiContent, option.WithAPIKey(apiKey))
	if genaiErr != nil {
		log.Fatal(genaiErr)
	}
	defer func(genaiClient *genai.Client) {
		err := genaiClient.Close()
		if err != nil {

		}
	}(genaiClient)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)
	e.GET("/ping", ping)
	e.GET("/text", text)

	e.Logger.Fatal(e.Start(addr))
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, &Message{
		Code:    200,
		Message: "welcome",
	})
}

func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, &Message{
		Code:    200,
		Message: "pong",
	})
}

func text(c echo.Context) error {
	text := c.FormValue("message")
	if text == "" {
		return c.JSON(http.StatusOK, &Message{
			Code:    500,
			Message: "欢迎使用AI CHAT,请输入内容",
		})
	}
	message, err := chat(text)
	if err != nil {
		return c.JSON(http.StatusOK, &Message{
			Code:    500,
			Message: err.Error(),
		})
	} else {
		return c.JSON(http.StatusOK, &Message{
			Code:    200,
			Message: message,
		})
	}
}

func chat(text string) (string, error) {
	model := genaiClient.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(genaiContent, genai.Text(text))
	if err != nil {
		return "", err
	}
	var result = ""
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				result = result + fmt.Sprintf("%v", part)
			}
		}
	}
	return result, nil
}
