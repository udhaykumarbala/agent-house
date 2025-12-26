// Todo App - Simple CRUD operations with localStorage persistence

const STORAGE_KEY = 'todos';

// Get todos from localStorage
function getTodos() {
    const data = localStorage.getItem(STORAGE_KEY);
    return data ? JSON.parse(data) : [];
}

// Save todos to localStorage
function saveTodos(todos) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(todos));
}

// Generate unique ID
function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
}

// Add a new todo
function addTodo(text) {
    const todos = getTodos();
    const newTodo = {
        id: generateId(),
        text: text.trim(),
        completed: false,
        createdAt: new Date().toISOString()
    };
    todos.push(newTodo);
    saveTodos(todos);
    return newTodo;
}

// Delete a todo
function deleteTodo(id) {
    const todos = getTodos();
    const filtered = todos.filter(todo => todo.id !== id);
    saveTodos(filtered);
}

// Toggle todo completion
function toggleTodo(id) {
    const todos = getTodos();
    const todo = todos.find(t => t.id === id);
    if (todo) {
        todo.completed = !todo.completed;
        saveTodos(todos);
    }
}

// Render todos to the DOM
function renderTodos() {
    const list = document.getElementById('todo-list');
    const todos = getTodos();

    list.innerHTML = '';

    if (todos.length === 0) {
        list.innerHTML = '<li>No todos yet. Add one above!</li>';
        return;
    }

    todos.forEach(todo => {
        const li = document.createElement('li');

        // Checkbox for completion
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.checked = todo.completed;
        checkbox.addEventListener('change', () => {
            toggleTodo(todo.id);
            renderTodos();
        });

        // Todo text
        const span = document.createElement('span');
        span.textContent = todo.text;
        if (todo.completed) {
            span.style.textDecoration = 'line-through';
            span.style.opacity = '0.6';
        }

        // Delete button
        const deleteBtn = document.createElement('button');
        deleteBtn.textContent = 'Delete';
        deleteBtn.addEventListener('click', () => {
            deleteTodo(todo.id);
            renderTodos();
        });

        li.appendChild(checkbox);
        li.appendChild(span);
        li.appendChild(deleteBtn);
        list.appendChild(li);
    });
}

// Form submission handler
document.getElementById('todo-form').addEventListener('submit', (e) => {
    e.preventDefault();
    const input = document.getElementById('todo-input');
    const text = input.value.trim();

    if (text) {
        addTodo(text);
        input.value = '';
        renderTodos();
    }
});

// Initial render
renderTodos();
