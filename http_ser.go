package main

import (
	pb "awesomeProject18/proto/user_service"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io"
	"net/http"
)

func getHash(payload string) string {
	hash := sha256.New()
	hash.Write([]byte(payload))
	return hex.EncodeToString(hash.Sum(nil))
}

func checkHash(payload string) (*pb.CheckHashResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHashingClient(conn)

	response, err := client.CheckHash(context.Background(), &pb.CheckHashRequest{Payload: payload})
	if err != nil {
		return nil, fmt.Errorf("error calling GetUser: %v", err)
	}

	return response, nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", CheckHash).Methods("POST")
	port := ":6262"
	fmt.Printf("Сервер запущен на порту %s...\n", port)

	http.ListenAndServe(port, r)
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
