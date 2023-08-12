import { router, store, events } from './framework/framework';
import { todoService, Todo } from './todo_mvc/todo';

// Add state to handle the current filter ("all", "active", "completed")
store.setState({ filter: 'all' });

// Add a utility function to get filtered todos based on the current filter
function getFilteredTodos(): Todo[] {
  const filter = store.getState().filter;
  if (filter === 'active') {
    return todoService.getTodos().filter((todo) => !todo.completed);
  } else if (filter === 'completed') {
    return todoService.getTodos().filter((todo) => todo.completed);
  }
  return todoService.getTodos();
}

// Display todo's list on the screen
function renderTodos() {

  // Display clear completed button if there are completed todos
  const clearCompletedButton = document.getElementById('clear-completed-button');
  if (clearCompletedButton) {
    const completedCount = todoService.getCompletedCount();
    if (completedCount > 0) {
      clearCompletedButton.style.display = 'inline-block';
    } else {
      clearCompletedButton.style.display = 'none';
    }
  }

  // Display footer if there are todos
  const all_todos = todoService.getTodos();
  const footer = document.querySelector("footer") as HTMLElement;
  if (all_todos.length > 0) {
    footer.style.display = "block";
  } else {
    footer.style.display = "none";
  }

  const counted_todos = todoService.getTodos().filter((todo) => !todo.completed);
  const todoCounter = document.getElementById('todo-counter');
  if (todoCounter) todoCounter.innerText = `${counted_todos.length} items left`;

  const todos = getFilteredTodos();

  const todoList = document.getElementById('todo-list');
  if (todoList) {
    todoList.innerHTML = '';
    todos.forEach((todo) => {
      const todoItem = document.createElement('li');
      todoItem.className = 'todo-item';

      const checkbox = document.createElement('input');
      checkbox.type = 'checkbox';
      checkbox.checked = todo.completed;

      events.on('change', checkbox, () => {
        todoService.toggleCompleted(todo.id);
        updateToggleAllCheckbox();
      });
      todoItem.appendChild(checkbox);

      const textInput = document.createElement('input');
      textInput.type = 'text';
      textInput.value = todo.text;
      textInput.readOnly = true;

      events.on('dblclick', textInput, () => {
        textInput.readOnly = false;
      });

      events.on('blur', textInput, () => {
        todoService.editTodo(todo.id, textInput.value);
        textInput.readOnly = true;
      });

      todoItem.appendChild(textInput);

      const deleteButton = document.createElement('button');
      deleteButton.textContent = 'Delete';

      events.on('click', deleteButton, () => {
        todoService.removeTodo(todo.id);
        updateToggleAllCheckbox();
      });

      todoItem.appendChild(deleteButton);

      todoList.appendChild(todoItem);
    });
  }
}

// Add event listener for toggle all checkbox
const toggleAllCheckbox = document.getElementById('toggle-all-checkbox') as HTMLInputElement;

let updatingToggleAll = false; // Flag to prevent change event from individual todos affecting toggleAllCheckbox

// Function to check/uncheck all todo checkboxes
function toggleAllTodos(checked: boolean) {
  const todos = todoService.getTodos();
  todos.forEach((todo) => (todo.completed = checked));
  updatingToggleAll = true; // Set the flag to true to prevent change event from affecting toggleAllCheckbox
  toggleAllCheckbox.checked = checked;
  updatingToggleAll = false; // Reset the flag to false after updating toggleAllCheckbox
  store.setState({ todos });
}

if (toggleAllCheckbox) {
  events.on('change', toggleAllCheckbox, () => {
    if (!updatingToggleAll) {
      const checked = toggleAllCheckbox.checked;
      toggleAllTodos(checked);
    }
  });
}

// Function to check the state of individual todo checkboxes and update the toggleAllCheckbox accordingly
function updateToggleAllCheckbox() {
  const todos = todoService.getTodos();
  const allChecked = todos.every((todo) => todo.completed);
  updatingToggleAll = true;
  toggleAllCheckbox.checked = allChecked;
  updatingToggleAll = false;
}

// Add todo to list (html "input" type=text, which add new todo after press enter)
const addTodoInput = document.getElementById('add-todo-input') as HTMLInputElement;
if (addTodoInput) {

  events.on('keypress', addTodoInput, (event: { key: string; }) => {
    if (event.key === 'Enter') {
      const text = addTodoInput.value;
      if (text.trim() !== '') {
        todoService.addTodo(text);
        addTodoInput.value = '';
      }
      updateToggleAllCheckbox();
    }
  });
}

// Add event listeners for filter buttons
const allButton = document.getElementById('all-button');
const activeButton = document.getElementById('active-button');
const completedButton = document.getElementById('completed-button');
if (allButton && activeButton && completedButton) {

  events.on('click', allButton, () => {
    router.navigate('/');
  });

  events.on('click', activeButton, () => {
    router.navigate('/active');
  });

  events.on('click', completedButton, () => {
    router.navigate('/completed');
  });

}

// add event listener for clear completed button
const clearCompletedButton = document.getElementById('clear-completed-button');
if (clearCompletedButton) {
  events.on('click', clearCompletedButton, () => {
    todoService.clearCompleted();
  });
}

// Subscribe to state changes and update the UI accordingly
store.subscribe(renderTodos);

// Initial render
renderTodos();

// Abstracting the DOM Routing System
router.addRoute('/', () => {
  store.setState({ filter: 'all' });
});

router.addRoute('/active', () => {
  store.setState({ filter: 'active' });
});

router.addRoute('/completed', () => {
  store.setState({ filter: 'completed' });
});

// Start the router
router.navigate(window.location.pathname);
