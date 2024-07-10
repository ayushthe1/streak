package chatbot

import (
	"context"
	"fmt"
	"time"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/option"
)

type ChatRequest struct {
	QueryInput QueryInput `json:"queryInput"`
}

type QueryInput struct {
	TextInput TextInput `json:"textInput"`
}

type TextInput struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsFile(CredentialFile))
	if err != nil {
		return "", err
	}

	defer client.Close()

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	request := &dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         text,
					LanguageCode: languageCode,
				},
			},
		},
	}

	response, err := client.DetectIntent(ctx, request)
	if err != nil {
		return "", err
	}

	responseText := response.QueryResult.FulfillmentMessages[0].GetText().Text[0]
	return responseText, nil
}

func ChatbotHandler(query string, sessionID string) (responeFromBot string, err error) {
	var req ChatRequest
	req.QueryInput.TextInput.Text = query

	// if err := c.BodyParser(&req); err != nil {
	// 	c.Status(fiber.StatusBadRequest)
	// 	return c.JSON(fiber.Map{"error": err})

	// }

	// sessionID := generateSessionID()
	text := req.QueryInput.TextInput.Text
	languageCode := "en"

	responseText, err := DetectIntentText(ProjectID, sessionID, text, languageCode)
	if err != nil {
		return "", fmt.Errorf("error while detecting intent : %v", err)
	}

	return responseText, nil

}
