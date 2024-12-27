package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "db.sqlite3"
	createTable = `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task TEXT NOT NULL,
			completed INTEGER NOT NULL
		)
	`
	selectTodos = "SELECT * FROM todos"
	insertTodo = "INSERT INTO todos (task, completed) VALUES (?, ?)"
	editTodo = "UPDATE todos SET task = ?, completed = ? WHERE id = ?"
)

type Todo struct {
	ID int `json:"id"`
	Task string `json:"task"`
	Completed int `json:"completed"`
}

func init() {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("データベースを正常に作成しました")
	}
}

func main() {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	http.HandleFunc("/api/todos", HandleCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodos(w, r, db)
		case http.MethodPost:
			createTodos(w, r, db)
		case http.MethodPut:
			editTodos(w, r, db)
		case http.MethodDelete:
			deleteTodos(w, r, db)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	
	fmt.Println("http://localhost:8080でサーバーを起動します")
	http.ListenAndServe(":8080", nil)
}

func getTodos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query(selectTodos)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var todos = []Todo{}
	
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Task, &todo.Completed)
		if err != nil {
			panic(err)
		}
		todos = append(todos, todo)
	}

	respondJSON(w, http.StatusOK, todos)
}

func createTodos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var todo Todo
	if err := decodeBody(r, &todo); err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return 
	}
	fmt.Println(todo.Task, todo.Completed)
	result, err := db.Exec(insertTodo, todo.Task, todo.Completed)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	todo.ID = int(id)
	
	fmt.Println("Successfully Created!")
	respondJSON(w, http.StatusCreated, todo)
}

func editTodos(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func deleteTodos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
}

func decodeBody(r *http.Request, v interface{}) error { 
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		panic(err)
	}
}

func HandleCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h(w, r)
	}
}
