package v1alpha1

// Shell is a path to the shell used to execute Instructions.
type Shell string

// Get returns current shell.
func (sh Shell) Get() string {
	if sh == "" {
		return "/bin/sh"
	}

	return string(sh)
}
