package main

import (
	"awesomeProject18/final-project-Ppolyak/proto/user_service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

func (s *server) CheckHash(ctx context.Context, request *user_service.CheckHashRequest) (*user_service.CheckHashResponse, error) {
	var result bool
	payload := request.Payload
	err := s.db.db.QueryRow("SELECT * FROM HASHES WHERE hash = $1", payload).Scan(result)
	if errors.Is(err, sql.ErrNoRows) {
		result = false
	} else {
		result = true
	}

	return &user_service.CheckHashResponse{
		Exists: result,
	}, nil
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
