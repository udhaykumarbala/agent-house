package skill

import (
	"regexp"
	"strings"
)

// MatchTrigger evaluates a trigger expression against the current context
// Supports:
//   - Simple conditions: "phase:research", "role:pm", "template:go-api"
//   - AND logic: "phase:development AND role:senior_dev"
//   - OR logic: "phase:research OR phase:planning"
//   - Mixed: "role:pm AND (phase:research OR phase:planning)"
func MatchTrigger(trigger string, ctx Context) bool {
	if trigger == "" {
		return true // Empty trigger matches everything
	}

	trigger = strings.TrimSpace(trigger)

	// Handle parentheses by evaluating inner expressions first
	for strings.Contains(trigger, "(") {
		trigger = evaluateParentheses(trigger, ctx)
	}

	return evaluateExpression(trigger, ctx)
}

// evaluateParentheses finds and evaluates the innermost parenthetical expression
func evaluateParentheses(expr string, ctx Context) string {
	// Find innermost parentheses
	start := -1
	for i, ch := range expr {
		if ch == '(' {
			start = i
		} else if ch == ')' && start >= 0 {
			// Evaluate the content inside parentheses
			inner := expr[start+1 : i]
			result := evaluateExpression(inner, ctx)
			// Replace with TRUE or FALSE
			replacement := "FALSE"
			if result {
				replacement = "TRUE"
			}
			return expr[:start] + replacement + expr[i+1:]
		}
	}
	return expr
}

// evaluateExpression handles AND/OR logic without parentheses
func evaluateExpression(expr string, ctx Context) bool {
	expr = strings.TrimSpace(expr)

	// Handle OR first (lower precedence)
	if strings.Contains(expr, " OR ") {
		parts := strings.Split(expr, " OR ")
		for _, part := range parts {
			if evaluateExpression(strings.TrimSpace(part), ctx) {
				return true
			}
		}
		return false
	}

	// Handle AND (higher precedence)
	if strings.Contains(expr, " AND ") {
		parts := strings.Split(expr, " AND ")
		for _, part := range parts {
			if !evaluateExpression(strings.TrimSpace(part), ctx) {
				return false
			}
		}
		return true
	}

	// Handle TRUE/FALSE literals (from parenthesis evaluation)
	if expr == "TRUE" {
		return true
	}
	if expr == "FALSE" {
		return false
	}

	// Handle NOT prefix
	if strings.HasPrefix(expr, "NOT ") {
		return !evaluateCondition(strings.TrimPrefix(expr, "NOT "), ctx)
	}

	// Single condition
	return evaluateCondition(expr, ctx)
}

// evaluateCondition evaluates a single key:value condition
func evaluateCondition(condition string, ctx Context) bool {
	condition = strings.TrimSpace(condition)

	parts := strings.SplitN(condition, ":", 2)
	if len(parts) != 2 {
		return false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "phase":
		return matchValue(ctx.Phase, value)
	case "role":
		return matchValue(ctx.Role, value)
	case "template":
		return matchValue(ctx.Template, value)
	default:
		return false
	}
}

// matchValue compares context value against pattern
// Supports exact match and wildcard (*)
func matchValue(contextValue, pattern string) bool {
	if pattern == "*" {
		return true
	}

	// Support simple wildcard patterns like "static-*"
	if strings.Contains(pattern, "*") {
		// Convert glob to regex
		regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
		regexPattern = strings.ReplaceAll(regexPattern, "\\*", ".*")
		matched, err := regexp.MatchString(regexPattern, contextValue)
		if err != nil {
			return false
		}
		return matched
	}

	return strings.EqualFold(contextValue, pattern)
}

// ValidateTrigger checks if a trigger expression is syntactically valid
func ValidateTrigger(trigger string) error {
	if trigger == "" {
		return nil
	}

	// Check balanced parentheses
	depth := 0
	for _, ch := range trigger {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		}
		if depth < 0 {
			return &TriggerError{Message: "unbalanced parentheses"}
		}
	}
	if depth != 0 {
		return &TriggerError{Message: "unbalanced parentheses"}
	}

	// Check for valid condition patterns
	conditionPattern := regexp.MustCompile(`(phase|role|template):[\w\-\*]+`)

	// Remove AND/OR/NOT and parentheses to find conditions
	cleaned := trigger
	cleaned = strings.ReplaceAll(cleaned, " AND ", " ")
	cleaned = strings.ReplaceAll(cleaned, " OR ", " ")
	cleaned = strings.ReplaceAll(cleaned, "NOT ", "")
	cleaned = strings.ReplaceAll(cleaned, "(", " ")
	cleaned = strings.ReplaceAll(cleaned, ")", " ")

	parts := strings.Fields(cleaned)
	for _, part := range parts {
		if part == "" {
			continue
		}
		if !conditionPattern.MatchString(part) {
			return &TriggerError{Message: "invalid condition: " + part}
		}
	}

	return nil
}

// TriggerError represents an error in trigger syntax
type TriggerError struct {
	Message string
}

func (e *TriggerError) Error() string {
	return "trigger error: " + e.Message
}

// ParseTriggerConditions extracts all conditions from a trigger expression
func ParseTriggerConditions(trigger string) []Condition {
	var conditions []Condition

	if trigger == "" {
		return conditions
	}

	// Find all key:value patterns
	conditionPattern := regexp.MustCompile(`(phase|role|template):([\w\-\*]+)`)
	matches := conditionPattern.FindAllStringSubmatch(trigger, -1)

	for _, match := range matches {
		if len(match) == 3 {
			conditions = append(conditions, Condition{
				Key:   match[1],
				Value: match[2],
			})
		}
	}

	return conditions
}

// Condition represents a single trigger condition
type Condition struct {
	Key   string // "phase", "role", or "template"
	Value string // The expected value
}

// ContextMatchesAny checks if context matches any of the given triggers
func ContextMatchesAny(ctx Context, triggers []string) bool {
	for _, trigger := range triggers {
		if MatchTrigger(trigger, ctx) {
			return true
		}
	}
	return false
}

// ContextMatchesAll checks if context matches all of the given triggers
func ContextMatchesAll(ctx Context, triggers []string) bool {
	for _, trigger := range triggers {
		if !MatchTrigger(trigger, ctx) {
			return false
		}
	}
	return true
}
