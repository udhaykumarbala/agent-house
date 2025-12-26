package agent

import "testing"

func TestCanDelegateTo(t *testing.T) {
	tests := []struct {
		from     Role
		to       Role
		expected bool
	}{
		// CEO can delegate to PM, Architect, Security
		{RoleCEO, RolePM, true},
		{RoleCEO, RoleArchitect, true},
		{RoleCEO, RoleSecurity, true},
		{RoleCEO, RoleSeniorDev, false}, // Not allowed
		{RoleCEO, RoleUX, false},        // Not allowed

		// PM delegations
		{RolePM, RoleUX, true},
		{RolePM, RoleUI, true},
		{RolePM, RoleSeniorDev, true},
		{RolePM, RoleCEO, false}, // Cannot delegate back up

		// Senior Dev can only delegate to Junior Dev
		{RoleSeniorDev, RoleJuniorDev, true},
		{RoleSeniorDev, RolePM, false},
		{RoleSeniorDev, RoleUI, false},

		// Junior Dev cannot delegate to anyone
		{RoleJuniorDev, RoleSeniorDev, false},
		{RoleJuniorDev, RolePM, false},
	}

	for _, tt := range tests {
		result := CanDelegateTo(tt.from, tt.to)
		if result != tt.expected {
			t.Errorf("CanDelegateTo(%s, %s) = %v, expected %v", tt.from, tt.to, result, tt.expected)
		}
	}
}

func TestFilterValidDelegations(t *testing.T) {
	// CEO tries to delegate to multiple agents
	delegations := []Role{RolePM, RoleArchitect, RoleSeniorDev, RoleUX}
	valid := FilterValidDelegations(RoleCEO, delegations)

	if len(valid) != 2 {
		t.Errorf("Expected 2 valid delegations, got %d", len(valid))
	}

	// Check that PM and Architect are in valid list
	hasPM, hasArchitect := false, false
	for _, r := range valid {
		if r == RolePM {
			hasPM = true
		}
		if r == RoleArchitect {
			hasArchitect = true
		}
	}

	if !hasPM || !hasArchitect {
		t.Errorf("Expected PM and Architect in valid delegations, got %v", valid)
	}
}

func TestGetAllowedDelegations(t *testing.T) {
	allowed := GetAllowedDelegations(RoleCEO)
	if len(allowed) != 3 {
		t.Errorf("Expected 3 allowed delegations for CEO, got %d", len(allowed))
	}

	juniorAllowed := GetAllowedDelegations(RoleJuniorDev)
	if len(juniorAllowed) != 0 {
		t.Errorf("Expected 0 allowed delegations for Junior Dev, got %d", len(juniorAllowed))
	}
}
