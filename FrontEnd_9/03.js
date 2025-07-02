import { createContext, useContext, useReducer, useEffect, useMemo, useRef } from 'react';

// 创建TodoContext
const TodoContext = createContext();

// 初始状态
const initialState = {
    todos: [],
};

// 定义action类型
const ADD_TODO = 'ADD_TODO';
const DELETE_TODO = 'DELETE_TODO';
const TOGGLE_TODO = 'TOGGLE_TODO';

// reducer函数
function todoReducer(state, action) {
    switch (action.type) {
        case ADD_TODO:
            return {
                ...state,
                todos: [
                    ...state.todos,
                    {
                        id: Date.now(),
                        text: action.text,
                        completed: false,
                    },
                ],
            };
        case DELETE_TODO:
            return {
                ...state,
                todos: state.todos.filter(todo => todo.id !== action.id),
            };
        case TOGGLE_TODO:
            return {
                ...state,
                todos: state.todos.map(todo =>
                    todo.id === action.id ? { ...todo, completed: !todo.completed } : todo
                ),
            };
        default:
            return state;
    }
}

// 自定义Hook useTodos
function useTodos() {
    const [state, dispatch] = useReducer(todoReducer, initialState);

    const addTodo = text => {
        dispatch({ type: ADD_TODO, text });
    };

    const deleteTodo = id => {
        dispatch({ type: DELETE_TODO, id });
    };

    const toggleTodo = id => {
        dispatch({ type: TOGGLE_TODO, id });
    };

    return { todos: state.todos, addTodo, deleteTodo, toggleTodo };
}

function TodoList() {
    const { todos, addTodo, deleteTodo, toggleTodo } = useContext(TodoContext);
    const inputRef = useRef(null);

    // 使用useEffect在组件挂载时打印消息
    useEffect(() => {
        console.log('Todo List已加载');
    }, []);

    // 使用useMemo优化列表渲染
    const memoizedTodos = useMemo(() => {
        return todos.map(todo => (
            <TodoItem
                key={todo.id}
                todo={todo}
                onDelete={deleteTodo}
                onToggle={toggleTodo}
            />
        ));
    }, [todos, deleteTodo, toggleTodo]);

    const handleSubmit = e => {
        e.preventDefault();
        const text = inputRef.current.value.trim();
        if (text) {
            addTodo(text);
            inputRef.current.value = '';
        }
    };

    return (
        <div style={styles.container}>
            <h1 style={styles.title}>Todo List</h1>
            <form onSubmit={handleSubmit} style={styles.form}>
                <input
                    type="text"
                    ref={inputRef}
                    placeholder="输入任务内容..."
                    style={styles.input}
                />
                <button type="submit" style={styles.addButton}>添加</button>
            </form>
            <ul style={styles.list}>
                {memoizedTodos}
            </ul>
        </div>
    );
}

function TodoItem({ todo, onDelete, onToggle }) {
    const itemStyle = {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '10px',
        margin: '5px 0',
        backgroundColor: '#f9f9f9',
        borderRadius: '4px',
        textDecoration: todo.completed ? 'line-through' : 'none',
        color: todo.completed ? '#888' : '#333',
    };

    return (
        <li style={itemStyle}>
            <div>
                <input
                    type="checkbox"
                    checked={todo.completed}
                    onChange={() => onToggle(todo.id)}
                    style={{ marginRight: '10px' }}
                />
                {todo.text}
            </div>
            <button
                onClick={() => onDelete(todo.id)}
                style={styles.deleteButton}
            >
                删除
            </button>
        </li>
    );
}

function TodoApp() {
    const todos = useTodos();

    return (
        <TodoContext.Provider value={todos}>
            <TodoList />
        </TodoContext.Provider>
    );
}

// CSS样式
const styles = {
    container: {
        maxWidth: '600px',
        margin: '0 auto',
        padding: '20px',
        fontFamily: 'Arial, sans-serif',
    },
    title: {
        textAlign: 'center',
        color: '#333',
    },
    form: {
        display: 'flex',
        marginBottom: '20px',
    },
    input: {
        flex: 1,
        padding: '10px',
        fontSize: '16px',
        border: '1px solid #ddd',
        borderRadius: '4px 0 0 4px',
    },
    addButton: {
        padding: '10px 20px',
        backgroundColor: '#4CAF50',
        color: 'white',
        border: 'none',
        borderRadius: '0 4px 4px 0',
        cursor: 'pointer',
    },
    list: {
        listStyle: 'none',
        padding: 0,
    },
    deleteButton: {
        padding: '5px 10px',
        backgroundColor: '#f44336',
        color: 'white',
        border: 'none',
        borderRadius: '4px',
        cursor: 'pointer',
    },
};

export default TodoApp;