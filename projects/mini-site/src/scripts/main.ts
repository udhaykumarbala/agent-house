import { CountdownTimer, type TimerState } from './timer'

const DURATION = 60 // seconds
const RING_CIRCUMFERENCE = 2 * Math.PI * 148 // matches SVG r=148

const COLOR_PRIMARY = '#E8C547'
const COLOR_DANGER = '#E85454'

// DOM elements
const minutesEl = document.getElementById('minutes')!
const secondsEl = document.getElementById('seconds')!
const statusEl = document.getElementById('status')!
const btnStart = document.getElementById('btn-start') as HTMLButtonElement
const btnReset = document.getElementById('btn-reset') as HTMLButtonElement
const progressRing = document.getElementById('progress-ring') as unknown as SVGCircleElement
const timerDisplay = document.getElementById('timer-display')!

function formatTime(totalSeconds: number): { min: string; sec: string } {
  const min = Math.floor(totalSeconds / 60)
  const sec = totalSeconds % 60
  return {
    min: String(min).padStart(2, '0'),
    sec: String(sec).padStart(2, '0'),
  }
}

function updateDisplay(remaining: number, total: number): void {
  const { min, sec } = formatTime(remaining)
  minutesEl.textContent = min
  secondsEl.textContent = sec

  // Update progress ring
  const progress = remaining / total
  const offset = RING_CIRCUMFERENCE * (1 - progress)
  progressRing.style.strokeDashoffset = String(offset)

  // Color shift when < 10s remaining
  if (remaining <= 10 && remaining > 0) {
    timerDisplay.classList.add('glow-danger')
    timerDisplay.classList.remove('glow')
    progressRing.style.stroke = COLOR_DANGER
    minutesEl.style.color = COLOR_DANGER
    secondsEl.style.color = COLOR_DANGER
  } else {
    timerDisplay.classList.remove('glow-danger')
    timerDisplay.classList.add('glow')
    progressRing.style.stroke = COLOR_PRIMARY
    minutesEl.style.color = ''
    secondsEl.style.color = ''
  }
}

function updateUI(state: TimerState): void {
  const labels: Record<TimerState, string> = {
    idle: 'Ready',
    running: 'Counting down',
    paused: 'Paused',
    finished: "Time\u2019s up",
  }

  statusEl.textContent = labels[state]
  statusEl.classList.add('animate-fade-in')
  requestAnimationFrame(() => {
    statusEl.classList.remove('animate-fade-in')
    void statusEl.offsetWidth // force reflow
    statusEl.classList.add('animate-fade-in')
  })

  // Button labels
  switch (state) {
    case 'idle':
      btnStart.textContent = 'Start'
      btnStart.classList.remove('hidden')
      btnReset.classList.add('hidden')
      break
    case 'running':
      btnStart.textContent = 'Pause'
      btnStart.classList.remove('hidden')
      btnReset.classList.remove('hidden')
      break
    case 'paused':
      btnStart.textContent = 'Resume'
      btnStart.classList.remove('hidden')
      btnReset.classList.remove('hidden')
      break
    case 'finished':
      btnStart.classList.add('hidden')
      btnReset.classList.remove('hidden')
      break
  }
}

function onFinish(): void {
  // Pulse animation on finish
  timerDisplay.classList.add('animate-pulse-slow')
  setTimeout(() => timerDisplay.classList.remove('animate-pulse-slow'), 4000)
}

// Initialize
const timer = new CountdownTimer(DURATION, {
  onTick: updateDisplay,
  onStateChange: updateUI,
  onFinish,
})

// Set initial display
updateDisplay(DURATION, DURATION)

// Event listeners
btnStart.addEventListener('click', () => {
  const state = timer.getState()
  if (state === 'idle' || state === 'finished') {
    if (state === 'finished') timer.reset()
    timer.start()
  } else if (state === 'running') {
    timer.pause()
  } else if (state === 'paused') {
    timer.resume()
  }
})

btnReset.addEventListener('click', () => {
  timer.reset()
})

// Keyboard shortcut: Space to start/pause, R to reset
document.addEventListener('keydown', (e) => {
  if (e.code === 'Space' && e.target === document.body) {
    e.preventDefault()
    btnStart.click()
  }
  if (e.code === 'KeyR' && e.target === document.body) {
    btnReset.click()
  }
})
