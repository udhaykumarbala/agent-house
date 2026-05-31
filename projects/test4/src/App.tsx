/**
 * DesignFlow — UX Designer Todo App
 * Root component. Theme is applied via dark class on <html>.
 */
function App() {
  return (
    <div className="min-h-dvh flex flex-col">
      {/* App will be built in subsequent phases */}
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <h1 className="font-heading text-h1 text-text-primary mb-2">DesignFlow</h1>
          <p className="text-body text-text-secondary">UX Designer Todo App</p>
          <p className="text-meta text-text-tertiary mt-4">Phase 1: Foundation — scaffolding complete</p>
        </div>
      </div>
    </div>
  )
}

export default App
