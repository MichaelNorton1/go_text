package hobbies

import (
	"fmt"
	"strings"

	"go_text/internal/db" // Adjust to match your go.mod
)

// FormatAddResponse creates the response text for #add
func FormatAddResponse(hobbyName string, created bool) string {
	if created {
		return fmt.Sprintf("✅ Added **%s** to your tracked hobbies!", strings.Title(hobbyName))
	}
	return fmt.Sprintf("⚠️ **%s** is already in your hobbies list.", strings.Title(hobbyName))
}

// FormatResultsResponse creates the weekly report message
func FormatResultsResponse(results []db.HobbyResult) string {
	if len(results) == 0 {
		return "📊 You don't have any hobbies set up yet! Use `#add <hobby>` to start tracking."
	}

	var sb strings.Builder
	sb.WriteString("📊 **Weekly Hobby Summary (Last 7 Days)**\n\n")

	for _, res := range results {
		sb.WriteString(fmt.Sprintf("• **%s**: %d completed\n", strings.Title(res.Name), res.Count))
	}

	return sb.String()
}
