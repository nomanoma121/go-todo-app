package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// username:password@(tcpname)(host:port)/dbname
	DSN         = "test:test@(db:3306)/test"
	createTable = `
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			task TEXT NOT NULL,
			completed INTEGER NOT NULL
		)
	`
	selectTodos    = "SELECT * FROM todos"
	insertTodo     = "INSERT INTO todos (task, completed) VALUES (?, ?)"
	editTodoById   = "UPDATE todos SET task = ?, completed = ? WHERE id = ?"
	deleteTodoById = "DELETE FROM todos WHERE id = ?"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed int    `json:"completed"`
}

func init() {
	db, err := sql.Open("mysql", DSN)
	// データベースに接続できない場合、5秒ごとに再接続を試みる
	if err != nil {
		fmt.Println("データベースに接続できませんでした")
	}
	defer db.Close()
	for {
		err := db.Ping()
		if err != nil {
			fmt.Println("データベースに接続できませんでした。再接続を試みます。")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("データベースを正常に作成しました")
	}
}

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		fmt.Println("データベースに接続できませんでした")
	}
	defer db.Close()
	for {
		err := db.Ping()
		if err != nil {
			fmt.Println("データベースに接続できませんでした。再接続を試みます。")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	//動的ルーティング
	http.HandleFunc("/api/todos/", HandleCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteTodos(w, r, db)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	//静的ルーティング
	http.HandleFunc("/api/todos", HandleCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodos(w, r, db)
		case http.MethodPost:
			createTodos(w, r, db)
		case http.MethodPut:
			editTodos(w, r, db)
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
	var todo Todo
	if err := decodeBody(r, &todo); err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(todo.ID, todo.Task, todo.Completed)
	_, err := db.Exec(editTodoById, todo.Task, todo.Completed, todo.ID)
	if err != nil {
		panic(err)
	}
	respondJSON(w, http.StatusOK, todo)
}

func deleteTodos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Println("this is running")
	idStr := r.URL.Path[len("/api/todos/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	_, err = db.Exec(deleteTodoById, id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully deleted.")
	respondJSON(w, http.StatusOK, id)
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h(w, r)
	}
}
