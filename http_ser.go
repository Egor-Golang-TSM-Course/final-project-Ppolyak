package main

import (
	"awesomeProject18/final-project-Ppolyak/proto/user_service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io"
	"net/http"
)

func checkHash(payload string) (*user_service.CheckHashResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer conn.Close()

	client := user_service.NewHashingClient(conn)

	response, err := client.CheckHash(context.Background(), &user_service.CheckHashRequest{Payload: payload})
	if err != nil {
		return nil, fmt.Errorf("error calling GetUser: %v", err)
	}

	return response, nil
}

func createHash(payload string) (*user_service.CreateHashResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer conn.Close()

	client := user_service.NewHashingClient(conn)

	response, err := client.CreateHash(context.Background(), &user_service.CreateHashRequest{Payload: payload})
	if err != nil {
		return nil, fmt.Errorf("error calling GetUser: %v", err)
	}

	return response, nil
}

func getHash(payload string) (*user_service.GetHashResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer conn.Close()

	client := user_service.NewHashingClient(conn)

	response, err := client.GetHash(context.Background(), &user_service.GetHashRequest{Payload: payload})
	if err != nil {
		return nil, fmt.Errorf("error calling GetUser: %v", err)
	}

	return response, nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/checkHash", CheckHash).Methods("POST")
	r.HandleFunc("/createHash", CreateHash).Methods("POST")
	r.HandleFunc("/getHash", GetHash).Methods("Get")
	port := ":6262"
	fmt.Printf("Сервер запущен на порту %s...\n", port)

	http.ListenAndServe(port, r)
}

func GetHash(w http.ResponseWriter, r *http.Request) {
	payload := r.URL.Query().Get("payload")
	if payload == "" {
		http.Error(w, "payload empty or not set", http.StatusBadRequest)
	}

	response, err := getHash(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response.Hash == "" {
		fmt.Fprintf(w, "Hash for payload %q not exist", payload)
	} else {
		fmt.Fprintf(w, "Hash for payload %q :%q", payload, response.Hash)
	}

	w.Header().Set("Content-Type", "application/json")

}

func CreateHash(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	defer r.Body.Close()
	var request struct {
		Payload string `json:"payload"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error unmarshall json", http.StatusInternalServerError)
		return
	}

	if request.Payload == "" {
		http.Error(w, "Payload cannot be empty", http.StatusBadRequest)
		return
	}

	response, err := createHash(request.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hash for payload %q : %q\n", request.Payload, response.Hash)

	w.Header().Set("Content-Type", "application/json")

}

func CheckHash(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var request struct {
		Payload string `json:"payload"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error unmarshall json", http.StatusInternalServerError)
		return
	}

	if request.Payload == "" {
		http.Error(w, "Payload cannot be empty", http.StatusBadRequest)
		return
	}

	response, err := checkHash(request.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response.Exists == true {
		fmt.Fprintf(w, "Hash for payload %q exists\n", request.Payload)
	} else {
		fmt.Fprintf(w, "Hash for payload %q does not exist", request.Payload)
	}

	w.Header().Set("Content-Type", "application/json")
}
