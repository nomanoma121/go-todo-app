import React, { useEffect, useState } from "react";
import "./App.css";

interface Todo {
  id: number;
  task: string;
  completed: number;
}

const App: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [input, setInput] = useState<string>("");
  const [editId, setEditId] = useState<number | null>(null);

  // Fetch todos from the backend
  useEffect(() => {
    fetch("http://localhost:8080/api/todos")
      .then((response) => response.json())
      .then((data) => setTodos(data))
      .catch((err) => console.error("Error fetching todos:", err));
  }, []);

  // Add or edit a todo
  const handleAddTodo = () => {
    if (input.trim() === "") return;

    if (editId !== null) {
      // Update existing todo
      const updatedTodo = { id: editId, task: input, completed: 0 };

      fetch("http://localhost:8080/api/todos", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedTodo),
      })
        .then((response) => response.json())
        .then((updated) => {
          setTodos((prevTodos) =>
            prevTodos.map((todo) =>
              todo.id === updated.id ? { ...todo, task: updated.task } : todo
            )
          );
          setEditId(null);
          setInput("");
        })
        .catch((err) => console.error("Error updating todo:", err));
    } else {
      // Create new todo
      const newTodo = { task: input, completed: 0 };

      fetch("http://localhost:8080/api/todos", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newTodo),
      })
        .then((response) => response.json())
        .then((created) => {
          setTodos((prevTodos) => [...prevTodos, created]);
          setInput("");
        })
        .catch((err) => console.error("Error creating todo:", err));
    }
  };

  // Delete a todo
  const handleDeleteTodo = (id: number) => {
    fetch(`http://localhost:8080/api/todos/${id}`, {
      method: "DELETE",
    })
      .then(() => {
        setTodos((prevTodos) => prevTodos.filter((todo) => todo.id !== id));
      })
      .catch((err) => console.error("Error deleting todo:", err));
  };

  // Prepare to edit a todo
  const handleEditTodo = (id: number, task: string) => {
    setEditId(id);
    setInput(task);
  };

  console.log(todos)

  return (
    <div className="App">
      <h1>Todo App</h1>
      <div className="todo-input">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Add a new task..."
        />
        <button onClick={handleAddTodo}>
          {editId !== null ? "Update" : "Add"}
        </button>
      </div>
      <div className="list-container">
        <ul className="todo-list">
          {todos.map((todo) => (
            <li key={todo.id}>
              <span>{todo.task}</span>
              <button onClick={() => handleEditTodo(todo.id, todo.task)}>
                Edit
              </button>
              <button onClick={() => handleDeleteTodo(todo.id)}>Delete</button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default App;
