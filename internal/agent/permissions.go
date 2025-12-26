package agent

import (
	"path/filepath"
	"strings"
)

// Permission defines what an agent can do
type Permission struct {
	CanCreateFiles  bool
	CanModifyFiles  bool
	CanExecuteCode  bool
	AllowedExtensions []string // Empty means all allowed
	AllowedPaths      []string // Path prefixes allowed (e.g., "docs/", "src/")
	DeniedPaths       []string // Path prefixes denied
}

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role]Permission{
	RoleCEO: {
		CanCreateFiles:    true,
		CanModifyFiles:    false, // Review only
		CanExecuteCode:    false,
		AllowedExtensions: []string{".md"},
		AllowedPaths:      []string{"docs/"},
	},
	RolePM: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    false,
		AllowedExtensions: []string{".md", ".txt"},
		AllowedPaths:      []string{"docs/", "specs/", "requirements/"},
	},
	RoleUX: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    false,
		AllowedExtensions: []string{".md", ".txt", ".svg"},
		AllowedPaths:      []string{"docs/", "design/", "wireframes/"},
	},
	RoleUI: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    false,
		AllowedExtensions: []string{".md", ".css", ".scss", ".svg", ".json"},
		AllowedPaths:      []string{"docs/", "design/", "styles/", "assets/"},
	},
	RoleSecurity: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    false,
		AllowedExtensions: []string{".md", ".txt", ".json", ".yaml", ".yml"},
		AllowedPaths:      []string{"docs/", "security/", "config/"},
	},
	RoleArchitect: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    false,
		AllowedExtensions: []string{".md", ".txt", ".json", ".yaml", ".yml", ".svg"},
		AllowedPaths:      []string{"docs/", "architecture/", "config/"},
	},
	RoleSeniorDev: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    true,
		AllowedExtensions: []string{}, // All extensions allowed
		AllowedPaths:      []string{}, // All paths allowed
		DeniedPaths:       []string{".git/", ".env", "secrets/"},
	},
	RoleJuniorDev: {
		CanCreateFiles:    true,
		CanModifyFiles:    true,
		CanExecuteCode:    true,
		AllowedExtensions: []string{}, // All extensions allowed
		AllowedPaths:      []string{"src/", "test/", "tests/", "lib/"},
		DeniedPaths:       []string{".git/", ".env", "secrets/", "config/"},
	},
}

// GetPermission returns the permission for a role
func GetPermission(role Role) Permission {
	if perm, ok := RolePermissions[role]; ok {
		return perm
	}
	// Default: no permissions
	return Permission{}
}

// CanCreate checks if the agent can create a file at the given path
func (a *Agent) CanCreate(filePath string) bool {
	perm := GetPermission(a.Role)

	if !perm.CanCreateFiles {
		return false
	}

	return checkPathPermission(filePath, perm)
}

// CanModify checks if the agent can modify a file at the given path
func (a *Agent) CanModify(filePath string) bool {
	perm := GetPermission(a.Role)

	if !perm.CanModifyFiles {
		return false
	}

	return checkPathPermission(filePath, perm)
}

// CanExecute checks if the agent can execute code
func (a *Agent) CanExecute() bool {
	return GetPermission(a.Role).CanExecuteCode
}

// checkPathPermission verifies if a path is allowed
func checkPathPermission(filePath string, perm Permission) bool {
	// Normalize path
	filePath = filepath.Clean(filePath)
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check denied paths first
	for _, denied := range perm.DeniedPaths {
		if strings.HasPrefix(filePath, denied) || strings.Contains(filePath, denied) {
			return false
		}
	}

	// Check extension if restrictions exist
	if len(perm.AllowedExtensions) > 0 {
		extAllowed := false
		for _, allowedExt := range perm.AllowedExtensions {
			if ext == allowedExt {
				extAllowed = true
				break
			}
		}
		if !extAllowed {
			return false
		}
	}

	// Check path prefix if restrictions exist
	if len(perm.AllowedPaths) > 0 {
		pathAllowed := false
		for _, allowedPath := range perm.AllowedPaths {
			if strings.HasPrefix(filePath, allowedPath) {
				pathAllowed = true
				break
			}
		}
		if !pathAllowed {
			return false
		}
	}

	return true
}

// PermissionError represents a permission denied error
type PermissionError struct {
	Agent    string
	Action   string
	Path     string
	Reason   string
}

func (e PermissionError) Error() string {
	return e.Agent + " cannot " + e.Action + " " + e.Path + ": " + e.Reason
}

// CheckCreatePermission returns an error if creation is not allowed
func (a *Agent) CheckCreatePermission(filePath string) error {
	if !a.CanCreate(filePath) {
		perm := GetPermission(a.Role)
		reason := "not allowed for this role"

		ext := filepath.Ext(filePath)
		if len(perm.AllowedExtensions) > 0 {
			allowed := false
			for _, e := range perm.AllowedExtensions {
				if e == ext {
					allowed = true
					break
				}
			}
			if !allowed {
				reason = "extension " + ext + " not allowed (allowed: " + strings.Join(perm.AllowedExtensions, ", ") + ")"
			}
		}

		if len(perm.AllowedPaths) > 0 {
			reason = "path not in allowed directories: " + strings.Join(perm.AllowedPaths, ", ")
		}

		return PermissionError{
			Agent:  a.Name,
			Action: "create",
			Path:   filePath,
			Reason: reason,
		}
	}
	return nil
}
