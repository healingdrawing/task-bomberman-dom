export class Events {
  private eventMap: Map<EventTarget, Map<string, Function[]>>;

  constructor() {
    this.eventMap = new Map();
  }

  on(eventName: string, listenTo: EventTarget, callback: Function) {
    let callbacksMap = this.eventMap.get(listenTo);
    if (!callbacksMap) {
      callbacksMap = new Map();
      this.eventMap.set(listenTo, callbacksMap);
    }

    let callbacks = callbacksMap.get(eventName);
    if (!callbacks) {
      callbacks = [];
      callbacksMap.set(eventName, callbacks);
      listenTo.addEventListener(eventName, (event) => {
        const callbacks = callbacksMap?.get(eventName);
        if (callbacks) {
          callbacks.forEach((callback) => callback(event));
        }
      });
    }

    callbacks.push(callback);
  }

  emit(eventName: string, listenTo: EventTarget, data?: any) {
    const callbacksMap = this.eventMap.get(listenTo);
    if (!callbacksMap) return;

    const callbacks = callbacksMap.get(eventName);
    if (callbacks) {
      callbacks.forEach((callback) => callback(data));
    }
  }
}
