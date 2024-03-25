package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"grpc/proto/user_service"
	"log"
	"net"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

// Реализация сервера
type server struct {
	*user_service.UnimplementedHashingServer
	db *PostgresDB
}

type PostgresDB struct {
	db *sqlx.DB
}

func generateHash(payload string) string {
	hash := sha256.New()
	hash.Write([]byte(payload))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *server) GetHash(ctx context.Context, request *user_service.GetHashRequest) (*user_service.GetHashResponse, error) {
	var hash string
	payload := request.Payload
	err := s.db.db.QueryRow("SELECT hash FROM hashes WHERE payload = $1", payload).Scan(&hash)
	if err != nil {
		log.Println("Error while querying", err)
	}
	return &user_service.GetHashResponse{Hash: hash}, nil
}

func (s *server) CheckHash(ctx context.Context, request *user_service.CheckHashRequest) (*user_service.CheckHashResponse, error) {
	var result bool
	payload := request.Payload
	err := s.db.db.QueryRow("SELECT * FROM HASHES WHERE payload = $1", payload).Scan(result) // exec и пример запроса
	if !errors.Is(err, sql.ErrNoRows) {
		result = true
	}

	log.Println("RES", result)
	return &user_service.CheckHashResponse{
		Exists: result,
	}, nil
}

func (s *server) CreateHash(ctx context.Context, request *user_service.CreateHashRequest) (*user_service.CreateHashResponse, error) {
	var hash string
	payload := request.Payload
	hash = generateHash(payload)
	res, err := s.db.db.Exec("INSERT INTO hashes(payload, hash) SELECT $1,$2 WHERE NOT EXISTS (SELECT 1 FROM hashes WHERE payload = $1)", payload, hash)
	if err != nil {
		log.Println(err)
		log.Println(res)
	}

	return &user_service.CreateHashResponse{Hash: hash}, nil
}

func main() {
	postgres, err := NewPostgres()
	if err != nil {
		log.Fatalf("failed to initialize PostgreSQL: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	user_service.RegisterHashingServer(s, &server{db: postgres})
	log.Println("gRPC server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewPostgres() (*PostgresDB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	p := &PostgresDB{db: db}
	err = p.MigrateDb()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return p, nil
}

func (p *PostgresDB) MigrateDb() error {
	queryDrop := `DROP TABLE IF EXISTS hashes`

	queriesCreate := []string{
		`CREATE TABLE IF NOT EXISTS hashes (
			id SERIAL PRIMARY KEY,
			payload TEXT UNIQUE,
			hash TEXT UNIQUE                  
		)`,
	}

	_, err := p.db.Exec(queryDrop)
	if err != nil {
		return fmt.Errorf("error while dropping tables: %w", err)
	}

	for _, q := range queriesCreate {
		_, err := p.db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Migration finished")

	return nil
}
