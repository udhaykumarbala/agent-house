import {
  type Choice,
  type GameState,
  type Result,
  RESULT_EMOJI,
  RESULT_TEXT,
  EMOJI_MAP,
  getRandomChoice,
  determineResult,
  getResultDetail,
  createInitialState,
  updateState,
} from './utils';

// DOM elements
const playerScoreEl = document.getElementById('player-score') as HTMLElement;
const cpuScoreEl = document.getElementById('cpu-score') as HTMLElement;
const drawScoreEl = document.getElementById('draw-score') as HTMLElement;
const resultArea = document.getElementById('result-area') as HTMLElement;
const resultEmoji = document.getElementById('result-emoji') as HTMLElement;
const resultText = document.getElementById('result-text') as HTMLElement;
const resultDetail = document.getElementById('result-detail') as HTMLElement;
const historySection = document.getElementById('history-section') as HTMLElement;
const historyList = document.getElementById('history-list') as HTMLElement;
const resetBtn = document.getElementById('reset-btn') as HTMLElement;
const streakBadge = document.getElementById('streak-badge') as HTMLElement;
const streakText = document.getElementById('streak-text') as HTMLElement;
const choiceButtons = document.querySelectorAll<HTMLButtonElement>('[data-choice]');

let state: GameState = createInitialState();
let isAnimating = false;

function animateScore(el: HTMLElement, value: number): void {
  el.textContent = String(value);
  el.classList.remove('score-bump');
  void el.offsetWidth; // force reflow
  el.classList.add('score-bump');
}

function updateScoreboard(): void {
  animateScore(playerScoreEl, state.playerScore);
  animateScore(cpuScoreEl, state.cpuScore);
  animateScore(drawScoreEl, state.drawScore);
}

function showResult(playerChoice: Choice, cpuChoice: Choice, result: Result): void {
  // Clear previous result classes
  resultArea.classList.remove('result-win', 'result-lose', 'result-draw');
  resultArea.classList.add(`result-${result}`);

  // Animate result emoji
  resultEmoji.style.opacity = '0';
  resultEmoji.style.transform = 'scale(0.5)';
  requestAnimationFrame(() => {
    resultEmoji.textContent = `${EMOJI_MAP[playerChoice]} vs ${EMOJI_MAP[cpuChoice]}`;
    resultEmoji.style.transition = 'all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)';
    resultEmoji.style.opacity = '1';
    resultEmoji.style.transform = 'scale(1)';
  });

  // Animate result text
  resultText.style.opacity = '0';
  setTimeout(() => {
    resultText.textContent = `${RESULT_EMOJI[result]} ${RESULT_TEXT[result]}`;
    resultText.className = 'font-display text-lg font-semibold animate-fade-up';

    if (result === 'win') resultText.classList.add('text-accent-win');
    else if (result === 'lose') resultText.classList.add('text-accent-lose');
    else resultText.classList.add('text-accent-draw');

    resultText.style.opacity = '1';
  }, 150);

  // Show detail
  resultDetail.textContent = getResultDetail(playerChoice, cpuChoice, result);
  resultDetail.classList.remove('hidden');
  resultDetail.className = 'text-text-muted text-sm mt-1 animate-fade-up';
}

function highlightButton(choice: Choice, result: Result): void {
  choiceButtons.forEach((btn) => {
    const btnChoice = btn.dataset.choice as Choice;
    btn.classList.remove('winner-highlight', 'loser-highlight');

    if (btnChoice === choice) {
      btn.classList.add(result === 'win' ? 'winner-highlight' : result === 'lose' ? 'loser-highlight' : '');
    }
  });
}

function addHistoryPip(result: Result): void {
  historySection.classList.remove('hidden');

  const pip = document.createElement('div');
  pip.className = `history-pip ${result} animate-bounce-in`;
  pip.title = result.charAt(0).toUpperCase() + result.slice(1);

  // Keep only the last 20 pips visible
  if (historyList.children.length >= 20) {
    historyList.removeChild(historyList.firstChild!);
  }

  historyList.appendChild(pip);
}

function updateStreak(): void {
  if (state.streak >= 3 && state.streakType !== 'draw') {
    const label = state.streakType === 'win' ? '🔥 Win' : '💀 Lose';
    streakText.textContent = `${label} streak: ${state.streak}`;
    streakBadge.classList.remove('hidden');

    // Re-trigger animation
    streakBadge.classList.remove('animate-bounce-in');
    void streakBadge.offsetWidth;
    streakBadge.classList.add('animate-bounce-in');
  } else {
    streakBadge.classList.add('hidden');
  }
}

function setButtonsDisabled(disabled: boolean): void {
  choiceButtons.forEach((btn) => {
    if (disabled) {
      btn.classList.add('disabled');
    } else {
      btn.classList.remove('disabled');
    }
  });
}

function shakeResultArea(): void {
  resultArea.classList.remove('animate-shake');
  void resultArea.offsetWidth;
  resultArea.classList.add('animate-shake');
}

async function playRound(playerChoice: Choice): Promise<void> {
  if (isAnimating) return;
  isAnimating = true;
  setButtonsDisabled(true);

  // Brief suspense delay
  resultEmoji.textContent = '🎲';
  resultEmoji.style.opacity = '0.5';
  resultText.textContent = 'Choosing...';
  resultText.className = 'font-display text-lg font-semibold text-text-secondary';
  resultDetail.classList.add('hidden');
  resultArea.classList.remove('result-win', 'result-lose', 'result-draw');

  shakeResultArea();

  await new Promise((resolve) => setTimeout(resolve, 500));

  const cpuChoice = getRandomChoice();
  const result = determineResult(playerChoice, cpuChoice);

  state = updateState(state, result);

  showResult(playerChoice, cpuChoice, result);
  updateScoreboard();
  highlightButton(playerChoice, result);
  addHistoryPip(result);
  updateStreak();

  setTimeout(() => {
    isAnimating = false;
    setButtonsDisabled(false);
    choiceButtons.forEach((btn) => {
      btn.classList.remove('winner-highlight', 'loser-highlight');
    });
  }, 800);
}

function resetGame(): void {
  state = createInitialState();
  playerScoreEl.textContent = '0';
  cpuScoreEl.textContent = '0';
  drawScoreEl.textContent = '0';
  resultEmoji.textContent = '⚔️';
  resultEmoji.style.opacity = '0.5';
  resultEmoji.style.transform = 'scale(1)';
  resultText.textContent = 'Choose your weapon';
  resultText.className = 'font-display text-lg font-semibold text-text-secondary';
  resultDetail.classList.add('hidden');
  resultArea.classList.remove('result-win', 'result-lose', 'result-draw');
  historyList.innerHTML = '';
  historySection.classList.add('hidden');
  streakBadge.classList.add('hidden');
}

// Event listeners
choiceButtons.forEach((btn) => {
  btn.addEventListener('click', () => {
    const choice = btn.dataset.choice as Choice;
    playRound(choice);
  });
});

resetBtn.addEventListener('click', resetGame);

// Keyboard support
document.addEventListener('keydown', (e: KeyboardEvent) => {
  const keyMap: Record<string, Choice> = {
    '1': 'rock',
    '2': 'paper',
    '3': 'scissors',
    r: 'rock',
    p: 'paper',
    s: 'scissors',
  };

  const choice = keyMap[e.key.toLowerCase()];
  if (choice) playRound(choice);
});
