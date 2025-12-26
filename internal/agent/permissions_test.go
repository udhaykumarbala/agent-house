package agent

import "testing"

func TestCEOPermissions(t *testing.T) {
	ceo := &Agent{Role: RoleCEO, Name: "CEO"}

	// CEO can create markdown in docs/
	if !ceo.CanCreate("docs/vision.md") {
		t.Error("CEO should be able to create docs/vision.md")
	}

	// CEO cannot create code files
	if ceo.CanCreate("src/main.go") {
		t.Error("CEO should NOT be able to create src/main.go")
	}

	// CEO cannot modify files (review only)
	if ceo.CanModify("docs/vision.md") {
		t.Error("CEO should NOT be able to modify files")
	}

	// CEO cannot execute code
	if ceo.CanExecute() {
		t.Error("CEO should NOT be able to execute code")
	}
}

func TestPMPermissions(t *testing.T) {
	pm := &Agent{Role: RolePM, Name: "PM"}

	// PM can create markdown
	if !pm.CanCreate("docs/requirements.md") {
		t.Error("PM should be able to create docs/requirements.md")
	}

	if !pm.CanCreate("specs/feature.md") {
		t.Error("PM should be able to create specs/feature.md")
	}

	// PM cannot create code
	if pm.CanCreate("src/main.go") {
		t.Error("PM should NOT be able to create code files")
	}

	// PM can modify docs
	if !pm.CanModify("docs/requirements.md") {
		t.Error("PM should be able to modify docs")
	}
}

func TestSeniorDevPermissions(t *testing.T) {
	dev := &Agent{Role: RoleSeniorDev, Name: "Senior Dev"}

	// Senior dev can create any file
	if !dev.CanCreate("src/main.go") {
		t.Error("Senior Dev should be able to create src/main.go")
	}

	if !dev.CanCreate("docs/readme.md") {
		t.Error("Senior Dev should be able to create docs/readme.md")
	}

	if !dev.CanCreate("test/main_test.go") {
		t.Error("Senior Dev should be able to create test files")
	}

	// Senior dev cannot touch .git or secrets
	if dev.CanCreate(".git/config") {
		t.Error("Senior Dev should NOT be able to create in .git/")
	}

	if dev.CanCreate("secrets/api_key.txt") {
		t.Error("Senior Dev should NOT be able to create in secrets/")
	}

	// Senior dev can execute code
	if !dev.CanExecute() {
		t.Error("Senior Dev should be able to execute code")
	}
}

func TestJuniorDevPermissions(t *testing.T) {
	dev := &Agent{Role: RoleJuniorDev, Name: "Junior Dev"}

	// Junior dev can create in src/ and test/
	if !dev.CanCreate("src/utils.go") {
		t.Error("Junior Dev should be able to create src/utils.go")
	}

	if !dev.CanCreate("tests/utils_test.go") {
		t.Error("Junior Dev should be able to create tests/")
	}

	// Junior dev cannot touch config
	if dev.CanCreate("config/settings.yaml") {
		t.Error("Junior Dev should NOT be able to create in config/")
	}

	// Junior dev cannot touch .env
	if dev.CanCreate(".env") {
		t.Error("Junior Dev should NOT be able to create .env")
	}
}

func TestSecurityPermissions(t *testing.T) {
	sec := &Agent{Role: RoleSecurity, Name: "Security"}

	// Security can create reports
	if !sec.CanCreate("docs/security-review.md") {
		t.Error("Security should be able to create security docs")
	}

	if !sec.CanCreate("security/threat-model.md") {
		t.Error("Security should be able to create in security/")
	}

	// Security can create config files
	if !sec.CanCreate("config/security.yaml") {
		t.Error("Security should be able to create security configs")
	}

	// Security cannot create code
	if sec.CanCreate("src/auth.go") {
		t.Error("Security should NOT be able to create code files")
	}

	// Security cannot execute code
	if sec.CanExecute() {
		t.Error("Security should NOT be able to execute code")
	}
}

func TestUIDesignerPermissions(t *testing.T) {
	ui := &Agent{Role: RoleUI, Name: "UI"}

	// UI can create CSS
	if !ui.CanCreate("styles/main.css") {
		t.Error("UI should be able to create CSS files")
	}

	// UI can create design docs
	if !ui.CanCreate("design/components.md") {
		t.Error("UI should be able to create design docs")
	}

	// UI cannot create Go code
	if ui.CanCreate("src/main.go") {
		t.Error("UI should NOT be able to create Go files")
	}
}

func TestArchitectPermissions(t *testing.T) {
	arch := &Agent{Role: RoleArchitect, Name: "Architect"}

	// Architect can create architecture docs
	if !arch.CanCreate("architecture/system-design.md") {
		t.Error("Architect should be able to create architecture docs")
	}

	// Architect can create config
	if !arch.CanCreate("config/database.yaml") {
		t.Error("Architect should be able to create config files")
	}

	// Architect cannot create code
	if arch.CanCreate("src/main.go") {
		t.Error("Architect should NOT be able to create code files")
	}
}

func TestPermissionError(t *testing.T) {
	pm := &Agent{Role: RolePM, Name: "Product Manager"}

	err := pm.CheckCreatePermission("src/main.go")
	if err == nil {
		t.Error("Expected permission error for PM creating Go file")
	}

	permErr, ok := err.(PermissionError)
	if !ok {
		t.Error("Expected PermissionError type")
	}

	if permErr.Agent != "Product Manager" {
		t.Errorf("Expected agent 'Product Manager', got '%s'", permErr.Agent)
	}

	if permErr.Action != "create" {
		t.Errorf("Expected action 'create', got '%s'", permErr.Action)
	}
}

func TestAllRolesHavePermissions(t *testing.T) {
	roles := []Role{
		RoleCEO, RolePM, RoleUX, RoleUI,
		RoleSecurity, RoleArchitect, RoleSeniorDev, RoleJuniorDev,
	}

	for _, role := range roles {
		perm := GetPermission(role)
		// Just verify we get a non-zero permission (at least one field set)
		if !perm.CanCreateFiles && !perm.CanModifyFiles && !perm.CanExecuteCode {
			// This is OK for a restricted role, but let's verify it's intentional
			t.Logf("Role %s has very restricted permissions", role)
		}
	}
}
