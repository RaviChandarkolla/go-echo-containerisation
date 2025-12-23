package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

type Item struct {
	ID        string `json:"id" dynamodbav:"id"`       // ← CRITICAL
	Score     int    `json:"score" dynamodbav:"score"` // ← CRITICAL
	CreatedAt string `json:"createdAt,omitempty" dynamodbav:"createdAt"`
}

var ddbClient *dynamodb.Client
var tableName = os.Getenv("DYNAMODB_TABLE")

func initDynamoDB() {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://dynamodb:8000"
	}

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"DUMMYIDEXAMPLE123",
			"DUMMYEXAMPLEKEY123456789",
			"")),
		config.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
			},
		}),
	)
	ddbClient = dynamodb.NewFromConfig(cfg)
}

func createTableIfNotExists() {
	// Wait for DynamoDB Local to fully start (10s grace period)
	time.Sleep(10 * time.Second)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	// Ignore all errors - table may already exist or DynamoDB still starting
	_, err := ddbClient.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Table %s creation ignored: %v", tableName, err)
	}
	log.Printf("✅ DynamoDB %s ready (table auto-created if needed)", tableName)
}

func putItem(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 putItem() ENTERED - parsing JSON...")
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Printf("❌ JSON decode FAILED: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("✅ JSON decoded: %+v", item)

	// ✅ CRITICAL: Validate required primary key
	if item.ID == "" {
		log.Println("❌ ID empty after decode")
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	item.CreatedAt = time.Now().Format(time.RFC3339)

	log.Println("🔍 About to marshal for DynamoDB...")
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Println("error:coming here1")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = ddbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		log.Println("error:coming here2")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("everything ok:coming here3")
	json.NewEncoder(w).Encode(item)
	log.Println("everything ok:coming here4")
}

func getItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, _ := ddbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if result.Item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	var item Item
	attributevalue.UnmarshalMap(result.Item, &item)
	json.NewEncoder(w).Encode(item)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("📥 %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("📤 %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	log.Println("🚀 Go app starting on :8080")
	initDynamoDB()
	createTableIfNotExists()

	r := mux.NewRouter()

	// ✅ CORRECT: Apply middleware to router
	r.Use(loggingMiddleware)

	r.HandleFunc("/items", putItem).Methods("POST")
	r.HandleFunc("/items/{id}", getItem).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
