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

  off(eventName: string, listenTo: EventTarget, callback: Function, info = '') {
    if (info) {
      console.log(`Removing ${eventName} event listener on ${listenTo} for ${info}`);
    } else {
      console.log(`Removing ${eventName} event listener on ${listenTo}`);
    }

    const callbacksMap = this.eventMap.get(listenTo);
    if (!callbacksMap) {
      console.warn(`No callbacksMap for ${eventName} on ${listenTo}`);
      return;
    }

    const callbacks = callbacksMap.get(eventName);
    if (!callbacks) {
      console.warn(`No callbacks for ${eventName} on ${listenTo}`);
      return;
    }

    const index = callbacks.indexOf(callback);
    if (index !== -1) {
      callbacks.splice(index, 1);
    }

    // If there are no more callbacks for the given eventName and listenTo, remove the event listener
    if (callbacks.length === 0) {
      console.log(`Removing event listener for ${eventName} on ${listenTo}`);
      listenTo.removeEventListener(eventName, (event) => {
        const callbacks = callbacksMap?.get(eventName);
        if (callbacks) {
          callbacks.forEach((callback) => callback(event));
        }
      });
      callbacksMap.delete(eventName);
    }

    // If there are no more event listeners for the given listenTo, remove it from the eventMap
    if (callbacksMap.size === 0) {
      console.log(`Removing ${listenTo} from eventMap`);
      this.eventMap.delete(listenTo);
    }
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
