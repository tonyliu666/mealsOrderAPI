package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func Init() error {
	// read the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := Init()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text("Please give me more five Carbonhydrate abundant foods especially for breakfast"))
	if err != nil {
		log.Fatal(err)
	}

	// printResponse(resp)
	meals := extractMealNames(resp)
	fmt.Println("Meals:")
	for _, meal := range meals {
		fmt.Println(meal)
	}
}

func extractMealNames(resp *genai.GenerateContentResponse) []string {
	var meals []string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					lines := strings.Split(string(text), "\n")
					for _, line := range lines {
						// Check if the line starts with a numbered list (e.g., "1.", "2.")
						line = strings.TrimSpace(line)
						if strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "2.") ||
							strings.HasPrefix(line, "3.") || strings.HasPrefix(line, "4.") || strings.HasPrefix(line, "5.") {
							// Extract the portion after the number
							parts := strings.SplitN(line, " ", 2)
							if len(parts) > 1 {
								spaceClear := strings.TrimSpace(parts[1])
								// Extract the name between ** markers
								mealName := extractBetweenMarkers(spaceClear, "**", "**")
								if mealName != "" {
									meals = append(meals, mealName)
								}
							}
						}
					}
				}
			}
		}
	}
	return meals
}

// Helper function to extract text between two markers
func extractBetweenMarkers(text, startMarker, endMarker string) string {
	start := strings.Index(text, startMarker)
	if start == -1 {
		return ""
	}
	start += len(startMarker)

	end := strings.Index(text[start:], endMarker)
	if end == -1 {
		return ""
	}

	return text[start : start+end]
}
