package gateway

import (
	"context"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func RecommendNutrition(nutrition string) ([]string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return []string{}, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text("Please give me more five"+nutrition+"abundant foods especially for breakfast"))
	if err != nil {
		return []string{}, err
	}
	meals := extractMealNames(resp)

	return meals, nil

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
	// delete the semicolon in each meal 
	for i := range meals {
		meals[i] = meals[i][:len(meals[i])-1]
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
