package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

type GenaiService struct {
	client *genai.Client
	ctx    context.Context
}

func NewGenaiService() *GenaiService {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err == nil {
		fmt.Println("Successfully created client")
		return &GenaiService{
			client: client,
			ctx:    ctx,
		}
	}
	log.Fatal("Error initiating a client")
	return nil
}

func (g GenaiService) CreateTextContent(prompt string) string {
	result, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		log.Fatal("Error creating content")
		return ""
	}
	return result.Text()
}

func (g GenaiService) CreateImageContent(prompt string) {

}

func main() {
	godotenv.Load()
	genaiService := NewGenaiService()
	if genaiService != nil {
		result := genaiService.CreateTextContent("Write python code to get factorial of a number")
		fmt.Println("RESULT", result)
		os.WriteFile("result.txt", []byte(result), 0777)
	}
}
