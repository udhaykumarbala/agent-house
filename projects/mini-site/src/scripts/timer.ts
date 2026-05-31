export type TimerState = 'idle' | 'running' | 'paused' | 'finished'

export interface TimerCallbacks {
  onTick: (remaining: number, total: number) => void
  onStateChange: (state: TimerState) => void
  onFinish: () => void
}

export class CountdownTimer {
  private total: number
  private remaining: number
  private intervalId: number | null = null
  private state: TimerState = 'idle'
  private callbacks: TimerCallbacks

  constructor(seconds: number, callbacks: TimerCallbacks) {
    this.total = seconds
    this.remaining = seconds
    this.callbacks = callbacks
  }

  start(): void {
    if (this.state === 'running') return

    this.setState('running')
    this.intervalId = window.setInterval(() => {
      this.remaining--
      this.callbacks.onTick(this.remaining, this.total)

      if (this.remaining <= 0) {
        this.stop()
        this.setState('finished')
        this.callbacks.onFinish()
      }
    }, 1000)
  }

  pause(): void {
    if (this.state !== 'running') return
    this.stop()
    this.setState('paused')
  }

  resume(): void {
    if (this.state !== 'paused') return
    this.start()
  }

  reset(): void {
    this.stop()
    this.remaining = this.total
    this.setState('idle')
    this.callbacks.onTick(this.remaining, this.total)
  }

  getState(): TimerState {
    return this.state
  }

  getRemaining(): number {
    return this.remaining
  }

  private stop(): void {
    if (this.intervalId !== null) {
      clearInterval(this.intervalId)
      this.intervalId = null
    }
  }

  private setState(state: TimerState): void {
    this.state = state
    this.callbacks.onStateChange(state)
  }
}
