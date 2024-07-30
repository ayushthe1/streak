package chatbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type DialogflowRequest struct {
	QueryResult struct {
		Action           string `json:"action"`
		QueryText        string `json:"queryText"`
		FullfillmentText string `json:"fulfillmentText"`
		Parameters       struct {
			Topic          string `json:"topic"`
			City           string `json:"city"`
			GeoCountryCode struct {
				Name   string `json:"name"`
				Alpha2 string `json:"alpha-2"`
				Alpha3 string `jaon:"alpha-3"`
			} `json:"geo-country-code"`
		} `json:"parameters"`
		Intent struct {
			DisplayName string `json:"displayName"`
		}
	} `json:"queryResult"`
}

type FulfillmentMessages struct {
	Text struct {
		Text []interface{} `json:"text"`
	} `json:"text"`
}

type DialogflowResponse struct {
	FulfillmentMessages []FulfillmentMessages `json:"fulfillmentMessages"`
}

func WebhookHandler(c *fiber.Ctx) error {
	var dfRequest DialogflowRequest
	if err := c.BodyParser(&dfRequest); err != nil {
		c.Status(http.StatusBadRequest)
		c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var res interface{}
	var err error
	var statusCode int

	switch intent := dfRequest.QueryResult.Intent.DisplayName; intent {
	// case "get-country-news":
	// 	countryCode := dfRequest.QueryResult.Parameters.GeoCountryCode.Alpha2
	// 	res, statusCode, err = getNews(countryCode)
	case "get-weather":
		city := dfRequest.QueryResult.Parameters.City
		res, statusCode, err = getWeather(city, dfRequest.QueryResult.QueryText)
	case "Default Fallback Intent":
		text := dfRequest.QueryResult.QueryText
		res, statusCode, err = getResponseFromGemini(text)

	default:
		log.Printf("No valid intent provided : %s", intent)
		err = fmt.Errorf("no valid intent provided : %s", intent)
		statusCode = http.StatusBadRequest

	}

	if err != nil {
		c.Status(statusCode)
		dR := DialogflowResponse{
			FulfillmentMessages: []FulfillmentMessages{
				{
					Text: struct {
						Text []interface{} `json:"text"`
					}{
						Text: []interface{}{err.Error()},
					},
				},
			},
		}
		return c.JSON(dR)
	}

	dR := DialogflowResponse{
		FulfillmentMessages: []FulfillmentMessages{
			{
				Text: struct {
					Text []interface{} `json:"text"`
				}{
					Text: []interface{}{res},
				},
			},
		},
	}
	c.Status(statusCode)
	return c.JSON(dR)

}

func getWeather(city string, question string) (interface{}, int, error) {
	if city == "" {
		return "", http.StatusBadRequest, fmt.Errorf("no city provided : %v", city)
	}

	weatherUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, WeatherAPIkey)

	req, err := http.NewRequest(http.MethodGet, weatherUrl, nil)
	if err != nil {
		log.Println("couldn't create request")
		return nil, http.StatusInternalServerError, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error making get request :%s", err)
		return nil, http.StatusInternalServerError, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error ,couldnt read response body :%s", err)
		return "", http.StatusInternalServerError, err
	}

	var response interface{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		log.Printf("error while unmarshalling the response : %s", err.Error())
		return "", http.StatusInternalServerError, err
	}

	// give response to gemini and get a summary
	query := fmt.Sprintf("I'm giving you a weather query and the weather data of that place. Give a answer to the query in less than 20 words. \n Query : %s, \n Data : %v \n ", question, response)
	return getResponseFromGemini(query)

}

func getNews(countryCode string) (interface{}, int, error) {
	if countryCode == "" {
		return "", http.StatusBadRequest, fmt.Errorf("no country code provided : %v", countryCode)
	}

	newsUrl := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=%s&apiKey=%s", countryCode, NewsAPIkey)

	req, err := http.NewRequest(http.MethodGet, newsUrl, nil)
	if err != nil {
		log.Println("couldn't create request")
		return nil, http.StatusInternalServerError, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error making get request :%s", err)
		return nil, http.StatusInternalServerError, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error ,couldnt read response body :%s", err)
		return "", http.StatusInternalServerError, err
	}

	var response interface{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		log.Printf("error while unmarshalling the response : %s", err.Error())
		return "", http.StatusInternalServerError, err
	}

	return response, http.StatusOK, nil

}

// function to get a response from a Ollama model
func getResponseFromOllama(text string) (interface{}, int, error) {

	log.Printf("Inside gRFO. text received is : %s", text)

	text = text + fmt.Sprint(".Answer in less than 20 words")

	type llmRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	if text == "" {
		return "", http.StatusBadRequest, fmt.Errorf("No text provided : %v", text)
	}

	ollamaServerURL := "http://host.docker.internal:11434/api/generate"

	llmreq := llmRequest{
		Model:  "phi3:mini",
		Prompt: text,
		Stream: false,
	}

	body, err := json.Marshal(llmreq)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, ollamaServerURL, bodyReader)
	if err != nil {
		log.Println("couldn't create request")
		return nil, http.StatusInternalServerError, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error making get request :%s", err)
		return nil, http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error ,couldnt read response body :%s", err)
		return "", http.StatusInternalServerError, err
	}

	// log.Printf("raw response body: %s", resBody)

	var response struct {
		Response string `json:"response"`
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		log.Printf("error while unmarshalling the response : %s", err.Error())
		return "", http.StatusInternalServerError, err
	}

	log.Println("response received from Ollama : ", response.Response)

	return response.Response, http.StatusOK, nil

}
