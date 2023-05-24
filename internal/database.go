package internal

import (
  "database/sql"
  "fmt"
  "os"

  _ "github.com/lib/pq"
)

const (
  conversationsInsertSql = "INSERT INTO conversations(id, title, import_time, create_time, update_time, current_node) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING"
  nodesInsertSql         = "INSERT INTO nodes(id, conversation_id, parent, import_time, children) VALUES($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING"
  messagesInsertSql      = "INSERT INTO messages(id, node_id, author_role, author_name, import_time, create_time, update_time, content_type, parts, end_turn, weight, recipient) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (id) DO NOTHING"
)

type DatabaseService struct{}

func NewDatabaseService() *DatabaseService {
  return &DatabaseService{}
}

func (ds *DatabaseService) GetDatabaseHandle() (*sql.DB, error) {
  host := os.Getenv("DB_HOST")
  port := os.Getenv("DB_PORT")
  user := os.Getenv("DB_USER")
  password := os.Getenv("DB_PASSWORD")
  dbname := os.Getenv("DB_NAME")

  connection := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    host,
    port,
    user,
    password,
    dbname,
  )

  db, e := sql.Open("postgres", connection)

  if e != nil {
    return nil, e
  }

  return db, nil
}
