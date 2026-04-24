package normalize

import (
	"encoding/json"
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	"regexp"
	"strings"
)

type postContent struct {
	ZhCn postLanguage `json:"zh_cn"`
}

type postLanguage struct {
	Title   string          `json:"title"`
	Content [][]postElement `json:"content"`
}

type postElement struct {
	Tag      string `json:"tag"`
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
	ImageKey string `json:"image_key,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

// SimpleMarkdownToPost converts a basic markdown string to Lark Post JSON string.
// It supports paragraphs and links in the format [text](url).
func SimpleMarkdownToPost(title, markdown string, mentions []types.Mention) (string, error) {
	lines := strings.Split(markdown, "\n")
	content := make([][]postElement, 0, len(lines))

	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		var paragraph []postElement

		matches := linkRegex.FindAllStringSubmatchIndex(line, -1)

		lastIndex := 0
		for _, match := range matches {
			start, end := match[0], match[1]
			textStart, textEnd := match[2], match[3]
			hrefStart, hrefEnd := match[4], match[5]

			if start > lastIndex {
				paragraph = append(paragraph, postElement{
					Tag:  "text",
					Text: line[lastIndex:start],
				})
			}

			paragraph = append(paragraph, postElement{
				Tag:  "a",
				Text: line[textStart:textEnd],
				Href: line[hrefStart:hrefEnd],
			})

			lastIndex = end
		}

		if lastIndex < len(line) {
			paragraph = append(paragraph, postElement{
				Tag:  "text",
				Text: line[lastIndex:],
			})
		}

		// Keep empty lines as empty text elements to maintain paragraph spacing
		if len(paragraph) == 0 {
			paragraph = append(paragraph, postElement{
				Tag:  "text",
				Text: "",
			})
		}

		content = append(content, paragraph)
	}

	// Prepend mentions
	if len(mentions) > 0 {
		atElements := ComposePostMentionElements(mentions)
		if len(atElements) > 0 {
			var first []postElement
			for _, el := range atElements {
				first = append(first, el, postElement{Tag: "text", Text: " "})
			}
			content = append([][]postElement{first}, content...)
		}
	}

	post := postContent{
		ZhCn: postLanguage{
			Title:   title,
			Content: content,
		},
	}

	bytes, err := json.Marshal(post)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
