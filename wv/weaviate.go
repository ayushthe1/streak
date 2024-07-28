package wv

import (
	"context"
	"log"
	"time"

	sModels "github.com/ayushthe1/streak/models"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

var Client *weaviate.Client

const ClassName = "Users"

func ConnectToWeaviate() {
	cfg := weaviate.Config{
		Host:    "weaviate:8080",
		Scheme:  "http",
		Headers: nil,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Println("error :couldn't connect to weaviate")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	live, err := client.Misc().LiveChecker().Do(ctx)
	if err != nil {
		log.Println("error : while checking weaviate live status")
		panic(err)
	}

	log.Println("**********Connected to Weaviate********** : ", live)

	Client = client

	err = createCollectionIfNotExist()
	if err != nil {
		log.Println("error while creating collection")
		panic(err)
	}

}

func createCollectionIfNotExist() error {
	log.Println("Inside createCollectionIfNotExist")
	// Check it the class already exist
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	class, err := Client.Schema().ClassGetter().
		WithClassName(ClassName).
		Do(ctx)
	// if err != nil {
	// 	log.Println("error : Error while getting class")
	// 	return err
	// }

	if class != nil {
		log.Printf("Class %s already exists, not creating it again", ClassName)
		return nil
	}

	// class don't exist, create it

	classObj := &models.Class{
		Class:      ClassName,
		Vectorizer: "text2vec-transformers",
		MultiTenancyConfig: &models.MultiTenancyConfig{
			Enabled: true,
		},
	}

	err = Client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		log.Println("error while creating class")
		return err
	}

	log.Println("Collection Created")
	return nil
}

// function to add a new user (tenant) to weaviate. Ensure that username is unique
func CreateNewTenant(username string) error {
	log.Println("Inside CreateNewTenant")
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := Client.Schema().TenantsCreator().
		WithClassName(ClassName).
		WithTenants(models.Tenant{Name: username}).
		Do(ctx)

	if err != nil {
		log.Println("error while creating tenant(user) :", username)
		return err
	}

	log.Printf("tenant %s successfully created in weaviate", username)
	return err

}

// function to add chat object belonging to a tenant(user) into weaviate
func AddNewChatIntoWeaviate(chat *sModels.Chat) error {
	log.Println("Inside AddNewChatIntoWeaviate")

	tenant := chat.From
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	object, err := Client.Data().Creator().
		WithClassName(ClassName). // The class to which the object will be added
		WithProperties(map[string]interface{}{
			"from":    chat.From,
			"to":      chat.To,
			"message": chat.Msg,
		}).
		WithTenant(tenant). // The tenant to which the object will be added
		Do(ctx)

	if err != nil {
		log.Println("error while inserting chat into weaviate")
		return err
	}

	log.Printf("Chat object successfully inserted into weaviate %v", object)

	return nil

}

func GetChatsRelatedToQuery(username string, query string) (*models.GraphQLResponse, error) {
	log.Println("Inside GetChatsRelatedToQuery")
	tenant := username

	certanity := float32(0.6)
	nearText := Client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{query}).
		WithCertainty(certanity)

	fields := []graphql.Field{
		{Name: "from"},
		{Name: "to"},
		{Name: "message"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := Client.GraphQL().Get().
		WithClassName(ClassName).
		WithFields(fields...).
		WithNearText(nearText).
		WithTenant(tenant).
		Do(ctx)

	if err != nil {
		log.Println("error while querying data from wv")
		return nil, err
	}

	return result, nil

}
