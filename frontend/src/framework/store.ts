type State = Record<string, any>;

export class Store {
  private state: State;
  private subscribers: Function[];

  constructor() {
    this.state = {};
    this.subscribers = [];
  }

  getState() {
    return this.state;
  }

  setState(newState: State) {
    this.state = { ...this.state, ...newState };
    this.notifySubscribers();
  }

  subscribe(subscriber: Function) {
    this.subscribers.push(subscriber);
  }

  private notifySubscribers() {
    this.subscribers.forEach((subscriber) => subscriber(this.state));
  }
}
