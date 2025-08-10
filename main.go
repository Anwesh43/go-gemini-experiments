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

const MANGA_TEMPLATE_PROMP = `
A single black-and-white manga page in traditional Japanese style, drawn with bold ink lines, screentones for shading, speed lines for action, and hand-drawn details. The page contains a dynamic number of panels (between 3 and 8) of varying sizes and shapes, arranged in a visually engaging layout.
Panels should include a mix of:

    Wide establishing shots for setting and mood

    Close-ups for emotion and detail

    Mid-shots for dialogue and interaction

    Dynamic action shots with motion lines and bold Japanese sound effects (e.g., “ドン!”, “ガキン!”, “ザッ!”)

    Flashback or emotional cutaways with lighter screentones

    Large dramatic panels for climactic moments
    Text should be integrated naturally into speech bubbles and narration boxes within the art.
    Style: High-contrast black-and-white, gritty Shōnen/Shōjo manga aesthetic, authentic hand-inked look, expressive faces, cinematic camera angles.
    Theme, characters, and setting should be customized per request, but layout and manga visual conventions should always be preserved.
`

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

func (g GenaiService) CreateImageContent(prompt string, output string) {
	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}
	response, err := g.client.Models.GenerateContent(g.ctx, IMAGE_MODEL_NAME, genai.Text(prompt), config)
	if err != nil {
		log.Fatal("Error generating image", err)
		return
	}
	for _, part := range response.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		}
		if part.InlineData != nil {
			buf := part.InlineData.Data
			fileName := output
			os.WriteFile(fileName, buf, 0644)
			fmt.Println("saved image file")
		}
	}
}

type MangaAgent struct {
	genaiService GenaiService
	topic        string
	script       string
	imagePrompt  string
}

func (ma *MangaAgent) GenerateScript() {
	ma.script = ma.genaiService.CreateTextContent(fmt.Sprintf("generate a manga script for this topic %s, just give the first page for first chapter", ma.topic))
	fmt.Println(ma.script)
	ma.imagePrompt = ma.genaiService.CreateTextContent(fmt.Sprintf("Generate an image prompt to generate first page as a image with 6-7 panesl and text bubbles for the script %s. Please give one option as prompt and no additional text. Please follow the %s", ma.script, MANGA_TEMPLATE_PROMP))
	fmt.Println(fmt.Sprintf("IMAGE_PROMPT %s", ma.imagePrompt))
}

func (ma *MangaAgent) GenerateMangaPage() {
	ma.genaiService.CreateImageContent(ma.imagePrompt, "manga.png")
}

func main() {
	godotenv.Load()
	genaiService := NewGenaiService()
	// text_chan := make(chan bool)
	// image_chan := make(chan bool)
	// if genaiService != nil {
	// 	go func() {
	// 		result := genaiService.CreateTextContent("Write python code to get factorial of a number")
	// 		fmt.Println("RESULT", result)
	// 		os.WriteFile("result.txt", []byte(result), 0777)
	// 		text_chan <- true
	// 	}()
	// 	go func() {
	// 		genaiService.CreateImageContent("A human playing football with a dog in cartoon format")
	// 		image_chan <- true
	// 	}()
	// }
	// k := 0
	// for k < 2 {

	// 	select {
	// 	case <-text_chan:
	// 		fmt.Println("Created text content")
	// 		k = k + 1
	// 	case <-image_chan:
	// 		fmt.Println("Created image content")
	// 		k = k + 1
	// 	}
	// }

	mangaAgent := MangaAgent{
		genaiService: *genaiService,
		topic:        "A boy getting powers and gets up from wheelchair",
	}
	mangaAgent.GenerateScript()
	mangaAgent.GenerateMangaPage()
}
