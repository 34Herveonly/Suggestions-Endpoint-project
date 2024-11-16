package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variable for MongoDB client
var client *mongo.Client

func main() {
	// Initialize the database connection
	err := connectToMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer disconnectFromMongoDB()

	// Start HTTP server
	http.HandleFunc("/suggestion", PostSuggestion)
	log.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// Function to connect to MongoDB Atlas
func connectToMongoDB() error {
	const uri = "mongodb+srv://coffee:JnZl556iMCJEtyAM@digitalsuggestionboxpro.seedu.mongodb.net/?retryWrites=true&w=majority&appName=DigitalSuggestionBoxProject"

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a new MongoDB client
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("error creating MongoDB client: %v", err)
	}

	// Test the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB Atlas!")
	return nil
}

// Disconnect from MongoDB
func disconnectFromMongoDB() {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}

// HTTP handler to post a suggestion
func PostSuggestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var suggestion Suggestion
	err := json.NewDecoder(r.Body).Decode(&suggestion)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection := client.Database("suggestionsDB").Collection("suggestions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, suggestion)
	if err != nil {
		http.Error(w, "Failed to insert suggestion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Suggestion added successfully")
}

// Suggestion structure
type Suggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
