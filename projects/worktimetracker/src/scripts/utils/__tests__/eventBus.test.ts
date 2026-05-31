import { describe, it, expect, beforeEach, vi } from 'vitest';
import { eventBus } from '../eventBus';

describe('eventBus', () => {
  beforeEach(() => {
    eventBus._reset();
  });

  it('should call handler when event is emitted', () => {
    const handler = vi.fn();
    eventBus.on('test', handler);
    eventBus.emit('test', 'arg1', 'arg2');
    expect(handler).toHaveBeenCalledWith('arg1', 'arg2');
  });

  it('should support multiple handlers for the same event', () => {
    const handler1 = vi.fn();
    const handler2 = vi.fn();
    eventBus.on('test', handler1);
    eventBus.on('test', handler2);
    eventBus.emit('test');
    expect(handler1).toHaveBeenCalledTimes(1);
    expect(handler2).toHaveBeenCalledTimes(1);
  });

  it('should remove handler with off()', () => {
    const handler = vi.fn();
    eventBus.on('test', handler);
    eventBus.off('test', handler);
    eventBus.emit('test');
    expect(handler).not.toHaveBeenCalled();
  });

  it('should not throw when emitting event with no listeners', () => {
    expect(() => eventBus.emit('nonexistent')).not.toThrow();
  });

  it('should not throw when removing non-existent handler', () => {
    expect(() => eventBus.off('test', vi.fn())).not.toThrow();
  });

  it('should clear all listeners with _reset()', () => {
    const handler = vi.fn();
    eventBus.on('test', handler);
    eventBus._reset();
    eventBus.emit('test');
    expect(handler).not.toHaveBeenCalled();
  });
});
