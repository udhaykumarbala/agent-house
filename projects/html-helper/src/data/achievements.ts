export interface Achievement {
  id: string;
  title: string;
  description: string;
  icon: string;
  condition: (progress: any) => boolean;
}

export const achievements: Achievement[] = [
  {
    id: 'first_tag',
    title: 'Tag Master',
    description: 'Created your first HTML tag',
    icon: '🏆',
    condition: (progress) => progress.completedLessons.length >= 1,
  },
  {
    id: 'tag_closer',
    title: 'Tag Closer',
    description: 'Closed 25 tags correctly',
    icon: '⭐',
    condition: (progress) => progress.completedLessons.length >= 5,
  },
  {
    id: 'link_legend',
    title: 'Link Legend',
    description: 'Created your first working link',
    icon: '💪',
    condition: (progress) => progress.completedLessons.includes(5),
  },
  {
    id: 'speed_demon',
    title: 'Speed Typer',
    description: 'Completed a lesson in under 60 seconds',
    icon: '🚀',
    condition: () => false, // Track separately with timing
  },
  {
    id: 'halfway_hero',
    title: 'Halfway Hero',
    description: 'Completed 7 lessons',
    icon: '🎨',
    condition: (progress) => progress.completedLessons.length >= 7,
  },
  {
    id: 'streak_keeper',
    title: 'Streak Keeper',
    description: 'Practiced 3 days in a row',
    icon: '🔥',
    condition: (progress) => progress.streak >= 3,
  },
  {
    id: 'html_wizard',
    title: 'HTML Wizard',
    description: 'Completed all 15 lessons!',
    icon: '💎',
    condition: (progress) => progress.completedLessons.length >= 15,
  },
  {
    id: 'perfectionist',
    title: 'Perfectionist',
    description: 'Completed 5 lessons without hints',
    icon: '✨',
    condition: () => false, // Track separately
  },
];

export function checkAchievements(progress: any): string[] {
  return achievements
    .filter(achievement =>
      !progress.achievements.includes(achievement.id) &&
      achievement.condition(progress)
    )
    .map(achievement => achievement.id);
}

export function getAchievementById(id: string): Achievement | undefined {
  return achievements.find(a => a.id === id);
}
