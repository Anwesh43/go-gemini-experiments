package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

const (
	IMAGE_MODEL_NAME    = "gemini-2.0-flash-preview-image-generation"
	TEXT_GEN_MODEL_NAME = "gemini-2.5-flash"
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
	result, err := g.client.Models.GenerateContent(g.ctx, TEXT_GEN_MODEL_NAME, genai.Text(prompt), nil)
	if err != nil {
		log.Fatal("Error creating content")
		return ""
	}
	return result.Text()
}

func (g GenaiService) CreateImageContent(prompt string) {
	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}
	response, err := g.client.Models.GenerateContent(g.ctx, IMAGE_MODEL_NAME, genai.Text(prompt), config)
	if err != nil {
		log.Fatal("Error generating image")
		return
	}
	for _, part := range response.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		}
		if part.InlineData != nil {
			buf := part.InlineData.Data
			fileName := "image_gen_1.png"
			os.WriteFile(fileName, buf, 0644)
			fmt.Println("saved image file")
		}
	}
}

func main() {
	godotenv.Load()
	genaiService := NewGenaiService()
	text_chan := make(chan bool)
	image_chan := make(chan bool)
	if genaiService != nil {
		go func() {
			result := genaiService.CreateTextContent("Write python code to get factorial of a number")
			fmt.Println("RESULT", result)
			os.WriteFile("result.txt", []byte(result), 0777)
			text_chan <- true
		}()
		go func() {
			genaiService.CreateImageContent("A human playing football with a dog in cartoon format")
			image_chan <- true
		}()
	}
	k := 0
	for k < 2 {

		select {
		case <-text_chan:
			fmt.Println("Created text content")
			k = k + 1
		case <-image_chan:
			fmt.Println("Created image content")
			k = k + 1
		}
	}
}
