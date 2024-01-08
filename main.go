package main

import (
	"fmt"
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
	defaultAddr = ":8080"
	addr        string
	apiKey      string
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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)
	e.GET("/ping", ping)

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
