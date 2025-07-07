import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
    const [todos, setTodos] = useState([]);
    const [newTodo, setNewTodo] = useState('');

    useEffect(() => {
        fetch('http://localhost:8092/todos')
            .then(response => response.json())
            .then(data => setTodos(data))
            .catch(error => console.error('Error:', error));
    }, []);

    const addTodo = () => {
        if (!newTodo.trim()) return;

        fetch('http://localhost:8092/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title: newTodo, completed: false }),
        })
            .then(response => response.json())
            .then(data => {
                setTodos([...todos, data]);
                setNewTodo('');
            });
    };

    const toggleTodo = (id) => {
        const todo = todos.find(t => t.id === id);

        fetch(`http://localhost:8092/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ ...todo, completed: !todo.completed }),
        })
            .then(response => response.json())
            .then(updatedTodo => {
                setTodos(todos.map(t => t.id === id ? updatedTodo : t));
            });
    };

    const deleteTodo = (id) => {
        fetch(`http://localhost:8092/todos/${id}`, {
            method: 'DELETE',
        })
            .then(() => {
                setTodos(todos.filter(t => t.id !== id));
            });
    };

    return (
        <div className="App">
            <h1>Todo List</h1>
            <div className="todo-form">
                <input
                    type="text"
                    value={newTodo}
                    onChange={(e) => setNewTodo(e.target.value)}
                    placeholder="Add new todo..."
                />
                <button onClick={addTodo}>Add</button>
            </div>
            <ul className="todo-list">
                {todos.map(todo => (
                    <li key={todo.id} className={todo.completed ? 'completed' : ''}>
                        <span onClick={() => toggleTodo(todo.id)}>{todo.title}</span>
                        <button onClick={() => deleteTodo(todo.id)}>Delete</button>
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default App;