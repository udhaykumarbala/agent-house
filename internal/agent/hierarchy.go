package agent

// DelegationHierarchy defines which agents can delegate to which.
// CEO can delegate to anyone. Others follow a logical flow.
var DelegationHierarchy = map[Role][]Role{
	RoleCEO:       {RolePM, RoleArchitect, RoleSecurity, RoleUX, RoleUI, RoleSeniorDev, RoleJuniorDev, RoleHR, RoleProjectMgr, RoleProcurement, RoleSiteEngineer, RoleHSE, RoleQAInspector}, // CEO can delegate to anyone
	RolePM:        {RoleUX, RoleUI, RoleArchitect, RoleSeniorDev},
	RoleArchitect: {RoleSeniorDev, RoleSecurity},
	RoleUX:        {RoleUI, RoleSeniorDev},
	RoleUI:        {RoleSeniorDev},
	RoleSecurity:  {RoleSeniorDev},
	RoleSeniorDev: {RoleJuniorDev},
	RoleJuniorDev: {}, // Leaf node - cannot delegate

	// EPC hierarchy
	RoleProjectMgr:   {RoleSiteEngineer, RoleProcurement, RoleHSE, RoleQAInspector, RoleHR},
	RoleHR:           {},
	RoleProcurement:  {},
	RoleSiteEngineer: {RoleQAInspector},
	RoleHSE:          {},
	RoleQAInspector:  {},
}

// ReviewHierarchy defines which agents can escalate to which for decisions/review.
// This is the inverse of delegation - agents can ask superiors for approval.
var ReviewHierarchy = map[Role][]Role{
	RoleJuniorDev: {RoleSeniorDev},
	RoleSeniorDev: {RoleArchitect, RoleSecurity},
	RoleUI:        {RoleUX},
	RoleUX:        {RolePM},
	RoleArchitect: {RoleCEO},
	RolePM:        {RoleCEO},
	RoleSecurity:  {RoleArchitect, RoleCEO},
	RoleCEO:       {}, // CEO is top - no one to escalate to

	// EPC review hierarchy
	RoleHR:           {RoleProjectMgr, RoleCEO},
	RoleProjectMgr:   {RoleCEO},
	RoleProcurement:  {RoleProjectMgr, RoleCEO},
	RoleSiteEngineer: {RoleProjectMgr},
	RoleHSE:          {RoleProjectMgr, RoleCEO},
	RoleQAInspector:  {RoleSiteEngineer, RoleProjectMgr},
}

// CanDelegateTo checks if the 'from' role can delegate to the 'to' role
func CanDelegateTo(from, to Role) bool {
	allowed, exists := DelegationHierarchy[from]
	if !exists {
		return false
	}

	for _, role := range allowed {
		if role == to {
			return true
		}
	}
	return false
}

// FilterValidDelegations returns only the valid delegations from the given role
func FilterValidDelegations(from Role, delegations []Role) []Role {
	var valid []Role
	for _, to := range delegations {
		if CanDelegateTo(from, to) {
			valid = append(valid, to)
		}
	}
	return valid
}

// GetAllowedDelegations returns the list of agents a role can delegate to
func GetAllowedDelegations(from Role) []Role {
	return DelegationHierarchy[from]
}

// CanReviewTo checks if the 'from' role can escalate to the 'to' role for review
func CanReviewTo(from, to Role) bool {
	allowed, exists := ReviewHierarchy[from]
	if !exists {
		return false
	}

	for _, role := range allowed {
		if role == to {
			return true
		}
	}
	return false
}

// FilterValidReviews returns only the valid review escalations from the given role
func FilterValidReviews(from Role, reviews []Role) []Role {
	var valid []Role
	for _, to := range reviews {
		if CanReviewTo(from, to) {
			valid = append(valid, to)
		}
	}
	return valid
}

// GetAllowedReviews returns the list of agents a role can escalate to for review
func GetAllowedReviews(from Role) []Role {
	return ReviewHierarchy[from]
}
