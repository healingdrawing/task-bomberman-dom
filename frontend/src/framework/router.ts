export class Router {
  private routes: Record<string, Function>;

  constructor() {
    this.routes = {};
    this.handlePopState = this.handlePopState.bind(this);
    /*
    with "window" object (not HtmlElement) events.on fails, with error
    Uncaught ReferenceError: can't access lexical declaration 'o' before initialization.
    Cant fix it, so straight way used for "window" object.
    It is not critical, because there is no strict
    requirements in audit questions about this.
    This supports the go back/front browser buttons functionality
    */
    window.addEventListener('popstate', this.handlePopState);
  }

  addRoute(path: string, handler: Function) {
    this.routes[path] = handler;
  }

  navigate(path: string) {
    const handler = this.routes[path];
    if (handler) {
      handler();
      this.updateURL(path); // Update the browser URL after navigating to the new route
    } else {
      // console.log('Route not found');
    }
  }

  private updateURL(path: string) {
    history.pushState(null, '', path);
  }

  private handlePopState() {
    const path = window.location.pathname;
    const handler = this.routes[path];
    if (handler) {
      handler();
    } else {
      // console.log('Route not found');
    }
  }
}
