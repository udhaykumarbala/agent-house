/**
 * EventBus — Pub/sub event system for cross-view communication.
 * Services emit events on data changes; views subscribe and re-render.
 */

type EventHandler = (...args: unknown[]) => void;

const listeners: Map<string, Set<EventHandler>> = new Map();

export const eventBus = {
  on(event: string, handler: EventHandler): void {
    if (!listeners.has(event)) {
      listeners.set(event, new Set());
    }
    listeners.get(event)!.add(handler);
  },

  off(event: string, handler: EventHandler): void {
    const handlers = listeners.get(event);
    if (handlers) {
      handlers.delete(handler);
      if (handlers.size === 0) {
        listeners.delete(event);
      }
    }
  },

  emit(event: string, ...args: unknown[]): void {
    const handlers = listeners.get(event);
    if (handlers) {
      handlers.forEach((handler) => handler(...args));
    }
  },

  /** Remove all listeners — for testing */
  _reset(): void {
    listeners.clear();
  },
};
