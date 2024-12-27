import React, { useState } from "react";
import "./App.css";

interface Todo {
  id: number;
  text: string;
}

const App: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [input, setInput] = useState<string>("");
  const [editId, setEditId] = useState<number | null>(null);

  const handleAddTodo = () => {
    if (input.trim() === "") return;

    if (editId !== null) {
      // 編集モード
      setTodos((prevTodos) =>
        prevTodos.map((todo) =>
          todo.id === editId ? { ...todo, text: input } : todo
        )
      );
      setEditId(null);
    } else {
      // 新規作成
      setTodos([...todos, { id: Date.now(), text: input }]);
    }
    setInput("");
  };

  const handleDeleteTodo = (id: number) => {
    setTodos((prevTodos) => prevTodos.filter((todo) => todo.id !== id));
  };

  const handleEditTodo = (id: number, text: string) => {
    setEditId(id);
    setInput(text);
  };

  return (
    <div className="App">
      <h1>React Todo App</h1>
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
              <span>{todo.text}</span>
              <button onClick={() => handleEditTodo(todo.id, todo.text)}>
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
