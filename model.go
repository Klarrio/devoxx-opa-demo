package main

type policyRequest struct {
	Resource file
	Subject  user
}

type page struct {
	Files []displayFile
	User  map[string]any
}

type displayFile struct {
	File  map[string]any
	Authz bool
}

type file struct {
	Owner          string
	Name           string
	Location       string
	Classification string
	Environment    string
	EmployeeID     string
}

type user struct {
	// Name            string `yaml:"Name"`
	Email           string `yaml:"Email"`
	WorkingLocation string `yaml:"WorkingLocation"`
	// EmployeeID      string `yaml:"EmployeeID"`
	Function string `yaml:"Function"`
}
