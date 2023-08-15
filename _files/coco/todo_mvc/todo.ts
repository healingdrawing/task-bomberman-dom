import { store, events } from '../framework/framework';

export interface Todo {
  id: number;
  text: string;
  completed: boolean;
}

export class TodoService {
  private todos: Todo[];

  constructor() {
    this.todos = [];
  }

  addTodo(text: string) {
    const todo: Todo = {
      id: Date.now(),
      text,
      completed: false,
    };
    this.todos.push(todo);
    store.setState({ todos: this.todos });
  }

  removeTodo(id: number) {
    this.todos = this.todos.filter((todo) => todo.id !== id);
    store.setState({ todos: this.todos });
  }

  editTodo(id: number, text: string) {
    this.todos = this.todos.map((todo) =>
      todo.id === id ? { ...todo, text } : todo
    );
    store.setState({ todos: this.todos });
  }

  toggleCompleted(id: number) {
    this.todos = this.todos.map((todo) =>
      todo.id === id ? { ...todo, completed: !todo.completed } : todo
    );
    store.setState({ todos: this.todos });
  }

  clearCompleted() {
    this.todos = this.todos.filter((todo) => !todo.completed);
    store.setState({ todos: this.todos });
  }

  getCompletedCount(): number {
    return this.todos.filter((todo) => todo.completed).length;
  }

  getTodos(): Todo[] {
    return this.todos;
  }
}

export const todoService = new TodoService();
