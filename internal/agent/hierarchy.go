package agent

// DelegationHierarchy defines which agents can delegate to which.
// This enforces a logical flow: managers → designers → developers
// to prevent chaotic delegation patterns.
var DelegationHierarchy = map[Role][]Role{
	RoleCEO:       {RolePM, RoleArchitect, RoleSecurity},
	RolePM:        {RoleUX, RoleUI, RoleArchitect, RoleSeniorDev},
	RoleArchitect: {RoleSeniorDev, RoleSecurity},
	RoleUX:        {RoleUI, RoleSeniorDev},
	RoleUI:        {RoleSeniorDev},
	RoleSecurity:  {RoleSeniorDev},
	RoleSeniorDev: {RoleJuniorDev},
	RoleJuniorDev: {}, // Leaf node - cannot delegate
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
