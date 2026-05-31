export type Choice = 'rock' | 'paper' | 'scissors';
export type Result = 'win' | 'lose' | 'draw';

export interface GameState {
  playerScore: number;
  cpuScore: number;
  drawScore: number;
  history: Result[];
  streak: number;
  streakType: Result | null;
}

export const CHOICES: Choice[] = ['rock', 'paper', 'scissors'];

export const EMOJI_MAP: Record<Choice, string> = {
  rock: '🪨',
  paper: '📄',
  scissors: '✂️',
};

export const RESULT_EMOJI: Record<Result, string> = {
  win: '🎉',
  lose: '😤',
  draw: '🤝',
};

export const RESULT_TEXT: Record<Result, string> = {
  win: 'You win!',
  lose: 'You lose!',
  draw: "It's a draw!",
};

const WINS_AGAINST: Record<Choice, Choice> = {
  rock: 'scissors',
  paper: 'rock',
  scissors: 'paper',
};

export function getRandomChoice(): Choice {
  return CHOICES[Math.floor(Math.random() * CHOICES.length)];
}

export function determineResult(player: Choice, cpu: Choice): Result {
  if (player === cpu) return 'draw';
  return WINS_AGAINST[player] === cpu ? 'win' : 'lose';
}

export function getResultDetail(player: Choice, cpu: Choice, result: Result): string {
  if (result === 'draw') return `Both chose ${EMOJI_MAP[player]}`;
  const winner = result === 'win' ? player : cpu;
  const loser = result === 'win' ? cpu : player;
  return `${EMOJI_MAP[winner]} ${winner} beats ${loser} ${EMOJI_MAP[loser]}`;
}

export function createInitialState(): GameState {
  return {
    playerScore: 0,
    cpuScore: 0,
    drawScore: 0,
    history: [],
    streak: 0,
    streakType: null,
  };
}

export function updateState(state: GameState, result: Result): GameState {
  const next = { ...state };

  if (result === 'win') next.playerScore++;
  else if (result === 'lose') next.cpuScore++;
  else next.drawScore++;

  next.history = [...state.history, result];

  if (state.streakType === result) {
    next.streak = state.streak + 1;
  } else {
    next.streak = 1;
    next.streakType = result;
  }

  return next;
}
