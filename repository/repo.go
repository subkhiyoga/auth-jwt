package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/subkhiyoga/auth-jwt/database"
	"github.com/subkhiyoga/auth-jwt/model"
)

func ConnectDB() *sql.DB {
	dbHost := database.DotEnv("DB_HOST")
	dbPort := database.DotEnv("DB_PORT")
	dbUser := database.DotEnv("DB_USER")
	dbPassword := database.DotEnv("DB_PASSWORD")
	dbName := database.DotEnv("DB_NAME")
	sslMode := database.DotEnv("SSL_MODE")
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	} else {
		log.Fatalln("Database successfully connected!")
	}

	return db
}

type LoginRepository interface {
	GetByUnameAndPassword(username string, password string) (*model.Credentials, error)
}

type loginRepository struct {
	db *sql.DB
}

func (r *loginRepository) GetByUnameAndPassword(username string, password string) (*model.Credentials, error) {
	query := "SELECT c.username, c.password FROM credentials c JOIN mahasiswa m ON c.username = m.user_name WHERE c.username = $1 AND c.password = $2"
	row := r.db.QueryRow(query, username, password)

	c := &model.Credentials{}
	err := row.Scan(&c.Username, &c.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}

	return c, nil
}

func NewMahasiswaRepo(db *sql.DB) LoginRepository {
	repo := new(loginRepository)
	repo.db = db

	return repo
}
