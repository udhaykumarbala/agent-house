package agent

import "testing"

func TestCanDelegateTo(t *testing.T) {
	tests := []struct {
		from     Role
		to       Role
		expected bool
	}{
		// CEO can delegate to anyone
		{RoleCEO, RolePM, true},
		{RoleCEO, RoleArchitect, true},
		{RoleCEO, RoleSecurity, true},
		{RoleCEO, RoleSeniorDev, true}, // CEO can delegate to anyone
		{RoleCEO, RoleUX, true},        // CEO can delegate to anyone
		{RoleCEO, RoleUI, true},
		{RoleCEO, RoleJuniorDev, true},

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
	// CEO tries to delegate to multiple agents - all should be valid now
	delegations := []Role{RolePM, RoleArchitect, RoleSeniorDev, RoleUX}
	valid := FilterValidDelegations(RoleCEO, delegations)

	if len(valid) != 4 {
		t.Errorf("Expected 4 valid delegations for CEO, got %d", len(valid))
	}

	// PM tries to delegate - cannot delegate to CEO
	pmDelegations := []Role{RoleCEO, RoleUX, RoleSeniorDev}
	validPM := FilterValidDelegations(RolePM, pmDelegations)

	if len(validPM) != 2 {
		t.Errorf("Expected 2 valid delegations for PM, got %d", len(validPM))
	}
}

func TestCanReviewTo(t *testing.T) {
	tests := []struct {
		from     Role
		to       Role
		expected bool
	}{
		// Junior Dev can review to Senior Dev
		{RoleJuniorDev, RoleSeniorDev, true},
		{RoleJuniorDev, RoleCEO, false},

		// Senior Dev can review to Architect or Security
		{RoleSeniorDev, RoleArchitect, true},
		{RoleSeniorDev, RoleSecurity, true},
		{RoleSeniorDev, RoleCEO, false},

		// PM can review to CEO
		{RolePM, RoleCEO, true},
		{RolePM, RoleSeniorDev, false},

		// CEO cannot review to anyone
		{RoleCEO, RolePM, false},
		{RoleCEO, RoleArchitect, false},
	}

	for _, tt := range tests {
		result := CanReviewTo(tt.from, tt.to)
		if result != tt.expected {
			t.Errorf("CanReviewTo(%s, %s) = %v, expected %v", tt.from, tt.to, result, tt.expected)
		}
	}
}

func TestFilterValidReviews(t *testing.T) {
	// Senior Dev tries to review to multiple agents
	reviews := []Role{RoleArchitect, RoleSecurity, RoleCEO, RolePM}
	valid := FilterValidReviews(RoleSeniorDev, reviews)

	if len(valid) != 2 {
		t.Errorf("Expected 2 valid reviews for Senior Dev, got %d", len(valid))
	}
}

func TestGetAllowedDelegations(t *testing.T) {
	allowed := GetAllowedDelegations(RoleCEO)
	if len(allowed) != 7 { // CEO can delegate to all 7 agents
		t.Errorf("Expected 7 allowed delegations for CEO, got %d", len(allowed))
	}

	juniorAllowed := GetAllowedDelegations(RoleJuniorDev)
	if len(juniorAllowed) != 0 {
		t.Errorf("Expected 0 allowed delegations for Junior Dev, got %d", len(juniorAllowed))
	}
}
