package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// User represents a user of the portal
type User struct {
	ID         string      `json:"id"`
	SecretCode string      `json:"secretCode"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints"`
}

// Complaint represents a complaint submitted by a user
type Complaint struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Rating   int    `json:"rating"`
	Resolved bool   `json:"resolved"`
}

var usersDB map[string]User = make(map[string]User)

func ReturnJsonResponse(res http.ResponseWriter, resMessage []byte) {
	res.Header().Set("content-type", "application/json")
	res.Write(resMessage)
}

func loginUser(secretCode string) (User, error) {
	// Loop through the usersDB to find a user with the provided secret code
	for _, user := range usersDB {
		if user.SecretCode == secretCode {
			return user, nil
		}
	}

	// If no user with the provided secret code is found, return an error
	return User{}, fmt.Errorf("invalid secret code")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" :"check your HTTP method : Invalid HTTP method executed",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	// Parse the request body to get the user's secret code
	var requestBody struct {
		SecretCode string `json:"secretCode"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		HandlerMessage := []byte(`{
			"success" : false,
			"message" : "Error parsing the req body data",
		}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	// Call the loginUser function to authenticate the user
	user, err := loginUser(requestBody.SecretCode)

	if err != nil {
		HandlerMessage := []byte(`{
			"success":false,
			"message":"Wrong secret code/user not found",
 }`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	// Convert the user to JSON format
	userJSON, err := json.Marshal(user)

	if err != nil {
		HandlerMessage := []byte(`{
   "success":false,
   "message":"Error parsing the user data",
}`)
		ReturnJsonResponse(w, HandlerMessage)
		return
	}

	// Set the Content-Type header and respond with the user details
	HandlerMessage := []byte(`{
		"success" : true,
		"message" : "user sign-in successfully",
	}`)
	ReturnJsonResponse(w, HandlerMessage)
	ReturnJsonResponse(w, userJSON)
	return
}

// Function to generate a unique ID
func generateUniqueID() string {
	// Get the current timestamp in nanoseconds
	timestamp := time.Now().UnixNano()

	// Generate a random number between 0 and 99999
	// This is just a simple way to introduce randomness, you may want to use a more robust method in production.
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(100000)

	// Combine timestamp and random number to create the unique ID
	uniqueID := fmt.Sprintf("%d%d", timestamp, randomNum)

	return uniqueID
}

// Function to generate a unique secret code
func generateUniqueSecretCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const secretCodeLength = 10

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random secret code by selecting characters from the charset
	secretCode := make([]byte, secretCodeLength)
	for i := range secretCode {
		secretCode[i] = charset[rand.Intn(len(charset))]
	}

	return string(secretCode)
}

// Function to register a new user
func registerUser(name, email string) User {
	userID := generateUniqueID()             // Implement a function to generate a unique ID
	secretCode := generateUniqueSecretCode() // Implement a function to generate a unique secret code
	newUser := User{
		ID:         userID,
		SecretCode: secretCode,
		Name:       name,
		Email:      email,
		Complaints: []Complaint{},
	}
	usersDB[userID] = newUser
	return newUser
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Expected POST.", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the name and email of the new user
	var newUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the registerUser function to create a new user
	user := registerUser(newUser.Name, newUser.Email)

	// Convert the user to JSON format
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and respond with the user details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}

// Function to submit a complaint
func submitComplaint(userID, title, summary string, rating int) Complaint {
	complaintID := generateUniqueID() // Implement a function to generate a unique complaint ID
	newComplaint := Complaint{
		ID:       complaintID,
		Title:    title,
		Summary:  summary,
		Rating:   rating,
		Resolved: false,
	}
	user, ok := usersDB[userID]
	if !ok {
		// Handle user not found error
	}
	user.Complaints = append(user.Complaints, newComplaint)
	usersDB[userID] = user
	return newComplaint
}

func submitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Expected POST.", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the complaint details
	var newComplaint struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Rating  int    `json:"rating"`
	}
	err := json.NewDecoder(r.Body).Decode(&newComplaint)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get the user ID from the request context (assuming you have implemented authentication)
	userID := "user123" // Replace with the actual user ID from the request context

	// Call the submitComplaint function to create a new complaint
	complaint := submitComplaint(userID, newComplaint.Title, newComplaint.Summary, newComplaint.Rating)

	// Convert the complaint to JSON format
	complaintJSON, err := json.Marshal(complaint)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and respond with the complaint details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(complaintJSON)
}

// Function to get all complaints for a user
func getAllComplaintsForUser(userID string) ([]Complaint, error) {
	user, ok := usersDB[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user.Complaints, nil
}

func getAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method. Expected GET.", http.StatusMethodNotAllowed)
		return
	}

	// Get the user ID from the request context (assuming you have implemented authentication)
	userID := "user123" // Replace with the actual user ID from the request context

	// Call the getAllComplaintsForUser function to get all complaints for the user
	complaints, err := getAllComplaintsForUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Convert the complaints to JSON format
	complaintsJSON, err := json.Marshal(complaints)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and respond with the complaints
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(complaintsJSON)
}

// Function to get all complaints for admin
func getAllComplaintsForAdmin() []Complaint {
	var allComplaints []Complaint
	for _, user := range usersDB {
		allComplaints = append(allComplaints, user.Complaints...)
	}
	return allComplaints
}

func getAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method. Expected GET.", http.StatusMethodNotAllowed)
		return
	}

	// Call the getAllComplaintsForAdmin function to get all complaints for the admin
	complaints := getAllComplaintsForAdmin()

	// Convert the complaints to JSON format
	complaintsJSON, err := json.Marshal(complaints)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and respond with the complaints
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(complaintsJSON)
}

// Function to view complaint by ID for either user or admin
func viewComplaint(userID, complaintID string) (Complaint, error) {
	user, ok := usersDB[userID]
	if !ok {
		return Complaint{}, fmt.Errorf("user not found")
	}
	for _, complaint := range user.Complaints {
		if complaint.ID == complaintID {
			return complaint, nil
		}
	}
	return Complaint{}, fmt.Errorf("complaint not found")
}

func viewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method. Expected GET.", http.StatusMethodNotAllowed)
		return
	}

	// Get the user ID from the request context (assuming you have implemented authentication)
	userID := "user123" // Replace with the actual user ID from the request context

	// Get the complaint ID from the query parameters
	complaintID := r.URL.Query().Get("complaintID")

	// Call the viewComplaint function to get the complaint data
	complaint, err := viewComplaint(userID, complaintID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Convert the complaint to JSON format
	complaintJSON, err := json.Marshal(complaint)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and respond with the complaint data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(complaintJSON)
}

// Function to resolve a complaint (only for admins)
func resolveComplaint(complaintID string) error {
	for _, user := range usersDB {
		for idx, complaint := range user.Complaints {
			if complaint.ID == complaintID {
				user.Complaints[idx].Resolved = true
				usersDB[user.ID] = user
				return nil
			}
		}
	}
	return fmt.Errorf("complaint not found")
}

func resolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is PATCH (or PUT, depending on your API design)
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method. Expected PATCH.", http.StatusMethodNotAllowed)
		return
	}

	// Get the complaint ID from the query parameters
	complaintID := r.URL.Query().Get("complaintID")

	// Call the resolveComplaint function to mark the complaint as resolved
	err := resolveComplaint(complaintID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Complaint resolved successfully"))
}

func main() {

	log.Println("Complaint API")

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/submitComplaint", submitComplaintHandler)
	http.HandleFunc("/getAllComplaintsForUser", getAllComplaintsForUserHandler)
	http.HandleFunc("/getAllComplaintsForAdmin", getAllComplaintsForAdminHandler)
	http.HandleFunc("/viewComplaint", viewComplaintHandler)
	http.HandleFunc("/resolveComplaint", resolveComplaintHandler)

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
