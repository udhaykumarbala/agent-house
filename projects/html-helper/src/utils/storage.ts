import localforage from 'localforage';

/**
 * Storage manager using localForage for enhanced localStorage
 */
export interface UserProgress {
  completedLessons: number[];
  currentLesson: number;
  xp: number;
  streak: number;
  lastVisit: string;
  achievements: string[];
  preferences: {
    theme: 'dark' | 'light';
    soundEnabled: boolean;
    animationsEnabled: boolean;
  };
}

const DEFAULT_PROGRESS: UserProgress = {
  completedLessons: [],
  currentLesson: 1,
  xp: 0,
  streak: 0,
  lastVisit: new Date().toISOString(),
  achievements: [],
  preferences: {
    theme: 'dark',
    soundEnabled: false,
    animationsEnabled: true,
  },
};

class StorageManager {
  private store: typeof localforage;

  constructor() {
    this.store = localforage.createInstance({
      name: 'HTMLHelper',
      storeName: 'user_progress',
    });
  }

  /**
   * Get user progress
   */
  async getProgress(): Promise<UserProgress> {
    try {
      const progress = await this.store.getItem<UserProgress>('progress');
      return progress || DEFAULT_PROGRESS;
    } catch {
      return DEFAULT_PROGRESS;
    }
  }

  /**
   * Save user progress
   */
  async saveProgress(progress: UserProgress): Promise<void> {
    try {
      await this.store.setItem('progress', progress);
    } catch (error) {
      console.error('Failed to save progress:', error);
    }
  }

  /**
   * Update lesson completion
   */
  async completeLesson(lessonId: number, xpGained: number = 25): Promise<void> {
    const progress = await this.getProgress();

    if (!progress.completedLessons.includes(lessonId)) {
      progress.completedLessons.push(lessonId);
      progress.xp += xpGained;
    }

    progress.currentLesson = lessonId + 1;
    progress.lastVisit = new Date().toISOString();

    await this.saveProgress(progress);
  }

  /**
   * Unlock achievement
   */
  async unlockAchievement(achievementId: string): Promise<boolean> {
    const progress = await this.getProgress();

    if (!progress.achievements.includes(achievementId)) {
      progress.achievements.push(achievementId);
      await this.saveProgress(progress);
      return true;
    }

    return false;
  }

  /**
   * Update streak
   */
  async updateStreak(): Promise<number> {
    const progress = await this.getProgress();
    const lastVisit = new Date(progress.lastVisit);
    const now = new Date();
    const daysDiff = Math.floor((now.getTime() - lastVisit.getTime()) / (1000 * 60 * 60 * 24));

    if (daysDiff === 1) {
      // Consecutive day
      progress.streak += 1;
    } else if (daysDiff > 1) {
      // Streak broken
      progress.streak = 1;
    }
    // Same day visit doesn't change streak

    progress.lastVisit = now.toISOString();
    await this.saveProgress(progress);

    return progress.streak;
  }

  /**
   * Update preferences
   */
  async updatePreferences(preferences: Partial<UserProgress['preferences']>): Promise<void> {
    const progress = await this.getProgress();
    progress.preferences = { ...progress.preferences, ...preferences };
    await this.saveProgress(progress);
  }

  /**
   * Reset all progress
   */
  async resetProgress(): Promise<void> {
    await this.store.clear();
  }

  /**
   * Export progress (for backup)
   */
  async exportProgress(): Promise<string> {
    const progress = await this.getProgress();
    return JSON.stringify(progress, null, 2);
  }

  /**
   * Import progress (from backup)
   */
  async importProgress(data: string): Promise<void> {
    try {
      const progress = JSON.parse(data) as UserProgress;
      await this.saveProgress(progress);
    } catch (error) {
      throw new Error('Invalid progress data');
    }
  }
}

export const storage = new StorageManager();
