package orchestrator

import (
	"fmt"
	"regexp"
	"strings"

	"pty-claude-test/internal/agent"
	"pty-claude-test/internal/message"
)

// executeQAPhase runs QA review for a completed development phase
func (o *Orchestrator) executeQAPhase(phase *DevelopmentPhase, taskStr, projectID string, result *Result) error {
	if o.config.Verbose {
		fmt.Printf("\n%s QA REVIEW - Phase %d: %s (Iteration %d)\n", GetPhaseEmoji(PhaseQA), phase.Index, phase.Name, phase.Iteration)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	// Build QA context with completion criteria
	qaContext := buildQAContext(taskStr, phase)

	// CEO acts as QA reviewer
	reviewer := agent.RoleCEO

	// Process QA review
	if err := o.processAgentInPhase(reviewer, qaContext, projectID, result); err != nil {
		return fmt.Errorf("QA review failed: %w", err)
	}

	// Parse QA decision from the last message
	lastMsg := result.Messages[len(result.Messages)-1]
	qaStatus, feedback := parseQADecision(lastMsg.Content)

	// Record QA review
	review := CreateQAReview(
		result.TaskID,
		phase.Index,
		phase.Name,
		reviewer,
		qaStatus,
		feedback,
		extractFailedCriteria(lastMsg.Content),
		phase.Iteration,
	)

	if err := RecordQAReview(o.config.ProjectDir, review); err != nil {
		if o.config.Verbose {
			fmt.Printf("⚠️  Failed to record QA review: %v\n", err)
		}
	}

	// Update phase QA status in development plan
	plan, err := LoadDevelopmentPlan(o.config.ProjectDir)
	if err != nil {
		return fmt.Errorf("failed to load development plan: %w", err)
	}

	if err := plan.UpdatePhaseQAStatus(o.config.ProjectDir, phase.Index, qaStatus, feedback); err != nil {
		return fmt.Errorf("failed to update QA status: %w", err)
	}

	// Log QA decision
	if o.config.Verbose {
		if qaStatus == QAStatusApproved {
			fmt.Printf("✅ QA APPROVED: Phase %d passed review\n", phase.Index)
		} else if qaStatus == QAStatusRejected {
			fmt.Printf("❌ QA REJECTED: Phase %d needs revision (Iteration %d/3)\n", phase.Index, phase.Iteration)
			fmt.Printf("   Feedback: %s\n", feedback)
		}
	}

	// Update phase reference to reflect changes
	*phase = *plan.GetPhase(phase.Index)

	return nil
}

// buildQAContext creates the context for QA review
func buildQAContext(originalTask string, phase *DevelopmentPhase) string {
	var sb strings.Builder

	sb.WriteString("## PHASE: QA REVIEW\n\n")
	sb.WriteString("You are reviewing a completed development phase. Your job is to verify the implementation meets all completion criteria.\n\n")

	sb.WriteString("### Original Task\n")
	sb.WriteString(originalTask)
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("### Development Phase: %s\n", phase.Name))
	sb.WriteString(fmt.Sprintf("**Description:** %s\n", phase.Description))
	sb.WriteString(fmt.Sprintf("**Iteration:** %d/3\n\n", phase.Iteration))

	if phase.Iteration > 1 && phase.QAFeedback != "" {
		sb.WriteString("### Previous QA Feedback\n")
		sb.WriteString(phase.QAFeedback)
		sb.WriteString("\n\n")
	}

	sb.WriteString("### Subtasks & Completion Criteria\n\n")
	for i, subtask := range phase.SubTasks {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, subtask.Title))
		if subtask.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", subtask.Description))
		}
		sb.WriteString("   \n   **Completion Criteria:**\n")
		for _, criterion := range subtask.CompletionCriteria {
			sb.WriteString(fmt.Sprintf("   - %s\n", criterion))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("### Your QA Review Process\n\n")
	sb.WriteString("1. **Read the specifications:**\n")
	sb.WriteString("   - `.plans/specs/ui-spec.md` - Check colors, fonts, spacing\n")
	sb.WriteString("   - `.plans/specs/ux-spec.md` - Check layout, wireframes\n")
	sb.WriteString("   - `.plans/specs/product-spec.md` - Check features, behavior\n\n")

	sb.WriteString("2. **Review the implementation files:**\n")
	sb.WriteString("   - Check if each completion criterion is met\n")
	sb.WriteString("   - Verify specs were followed\n")
	sb.WriteString("   - Test functionality if applicable\n\n")

	sb.WriteString("3. **Make your decision:**\n\n")

	sb.WriteString("**If all criteria are met:**\n")
	sb.WriteString("```\n")
	sb.WriteString(fmt.Sprintf("QA_APPROVED: %s\n\n", phase.Name))
	sb.WriteString("✅ All completion criteria met:\n")
	sb.WriteString("- [Criterion 1]: Verified working\n")
	sb.WriteString("- [Criterion 2]: Matches specification\n\n")
	sb.WriteString("Phase is complete and meets quality standards.\n")
	sb.WriteString("```\n\n")

	sb.WriteString("**If issues are found:**\n")
	sb.WriteString("```\n")
	sb.WriteString(fmt.Sprintf("QA_REJECTED: %s\n\n", phase.Name))
	sb.WriteString(fmt.Sprintf("❌ Issues found (Iteration %d/3):\n\n", phase.Iteration))
	sb.WriteString("FAILED CRITERIA:\n")
	sb.WriteString("- [Criterion]: [Specific issue with file references]\n\n")
	sb.WriteString("REQUIRED FIXES:\n")
	sb.WriteString("1. [Specific fix needed]\n")
	sb.WriteString("2. [Specific fix needed]\n")
	sb.WriteString("```\n\n")

	sb.WriteString("### Quality Standards\n\n")
	sb.WriteString("- **Be specific**: Reference file names and line numbers\n")
	sb.WriteString("- **Check specs**: Verify colors, fonts, layout match specifications\n")
	sb.WriteString("- **Test thoroughly**: Check responsive design, interactive elements\n")
	sb.WriteString("- **No generic colors**: Reject AI defaults like #3B82F6 or #10B981\n")

	return sb.String()
}

// parseQADecision extracts QA status and feedback from response
func parseQADecision(content string) (QAStatus, string) {
	content = strings.TrimSpace(content)

	// Check for QA_APPROVED marker
	approvedRegex := regexp.MustCompile(`(?i)QA_APPROVED:\s*(.+)`)
	if matches := approvedRegex.FindStringSubmatch(content); len(matches) > 1 {
		// Extract everything after QA_APPROVED as feedback
		feedback := strings.TrimSpace(strings.Split(content, "QA_APPROVED:")[1])
		return QAStatusApproved, feedback
	}

	// Check for QA_REJECTED marker
	rejectedRegex := regexp.MustCompile(`(?i)QA_REJECTED:\s*(.+)`)
	if matches := rejectedRegex.FindStringSubmatch(content); len(matches) > 1 {
		// Extract everything after QA_REJECTED as feedback
		feedback := strings.TrimSpace(strings.Split(content, "QA_REJECTED:")[1])
		return QAStatusRejected, feedback
	}

	// Default to pending if no clear decision
	return QAStatusPending, content
}

// extractFailedCriteria extracts failed criteria from QA feedback
func extractFailedCriteria(content string) []string {
	var criteria []string

	// Look for FAILED CRITERIA section
	lines := strings.Split(content, "\n")
	inFailedSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(strings.ToUpper(line), "FAILED CRITERIA") {
			inFailedSection = true
			continue
		}

		if inFailedSection {
			// Stop at next section header
			if strings.Contains(strings.ToUpper(line), "REQUIRED FIXES") ||
				strings.Contains(strings.ToUpper(line), "FILES WITH ISSUES") {
				break
			}

			// Extract criteria from list items
			if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
				criterion := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"))
				if criterion != "" {
					criteria = append(criteria, criterion)
				}
			}
		}
	}

	return criteria
}

// buildIterationContext creates context for development iteration with QA feedback
func buildIterationContext(taskStr string, phase *DevelopmentPhase) string {
	var sb strings.Builder

	sb.WriteString("## PHASE: DEVELOPMENT ITERATION\n\n")
	sb.WriteString(fmt.Sprintf("### Phase %d: %s (Iteration %d/3)\n\n", phase.Index, phase.Name, phase.Iteration))

	sb.WriteString("The previous implementation was reviewed by QA and needs revision.\n\n")

	sb.WriteString("### QA Feedback\n\n")
	sb.WriteString(phase.QAFeedback)
	sb.WriteString("\n\n")

	sb.WriteString("### Original Requirements\n")
	sb.WriteString(taskStr)
	sb.WriteString("\n\n")

	sb.WriteString("### Your Task\n\n")
	sb.WriteString("Review the QA feedback above and fix the identified issues.\n\n")

	sb.WriteString("**IMPORTANT:**\n")
	sb.WriteString("- Address EVERY issue mentioned in the QA feedback\n")
	sb.WriteString("- Reference the specification files (ui-spec.md, ux-spec.md, product-spec.md)\n")
	sb.WriteString("- Do not introduce new issues\n")
	sb.WriteString("- Test your changes before submitting\n\n")

	sb.WriteString("### Subtasks & Completion Criteria\n\n")
	for i, subtask := range phase.SubTasks {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, subtask.Title))
		sb.WriteString("   **Completion Criteria:**\n")
		for _, criterion := range subtask.CompletionCriteria {
			sb.WriteString(fmt.Sprintf("   - %s\n", criterion))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("When done, signal: COMPLETE: Iteration finished, ready for QA review\n")

	return sb.String()
}

// notifyQADecision sends a notification about QA decision
func (o *Orchestrator) notifyQADecision(phase *DevelopmentPhase, qaStatus QAStatus, projectID, taskID string, result *Result) {
	var msgType message.MessageType
	var content string

	if qaStatus == QAStatusApproved {
		msgType = message.TypeSystem
		content = fmt.Sprintf("✅ QA Approved: Phase %d - %s", phase.Index, phase.Name)
	} else if qaStatus == QAStatusRejected {
		msgType = message.TypeSystem
		content = fmt.Sprintf("❌ QA Rejected: Phase %d - %s (Iteration %d/3)", phase.Index, phase.Name, phase.Iteration)
	} else {
		return
	}

	msg := message.NewMessage(msgType, "orchestrator", "all", content)
	msg.Metadata.ProjectID = projectID
	msg.Metadata.TaskID = taskID

	o.store.Add(msg)
	o.notify(msg)
	result.Messages = append(result.Messages, msg)
}
