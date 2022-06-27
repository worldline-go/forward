package handler

import "strings"

// FilterMethods using in handler http.
// default Allow has all methods.
type FilterMethods struct {
	Allow map[string]struct{} `json:"allow"`
	Deny  map[string]struct{} `json:"deny"`
}

// NewFilterMethods creates a new FilterMethods.
func NewFilterMethods() *FilterMethods {
	return &FilterMethods{
		Allow: map[string]struct{}{"*": {}},
		Deny:  make(map[string]struct{}),
	}
}

// String returns a string representation of the FilterMethods.
func (f *FilterMethods) String() string {
	return "allow: " + printMethods(f.Allow) + "; deny: " + printMethods(f.Deny)
}

// IsAllow returns true if the method is allowed.
func (f *FilterMethods) IsAllow(method string) bool {
	_, ok := f.Allow[strings.ToUpper(method)]
	return ok
}

// IsDeny returns true if the method is denied.
func (f *FilterMethods) IsDeny(method string) bool {
	_, ok := f.Deny[strings.ToUpper(method)]
	return ok
}

// Match returns true if the method is allowed.
func (f *FilterMethods) Match(method string) bool {
	if f.IsDeny("*") || f.IsDeny(method) {
		return false
	}

	if f.IsAllow("*") || f.IsAllow(method) {
		return true
	}

	return false
}

// Parse parses a string into a FilterMethods.
func (f *FilterMethods) Parse(methods []string) {
	delete(f.Allow, "*")

	allowUsed := false

	for _, m := range methods {
		if len(m) > 0 && m[0] == '-' {
			f.Deny[strings.ToUpper(m[1:])] = struct{}{}
		} else {
			f.Allow[strings.ToUpper(m)] = struct{}{}
			allowUsed = true
		}
	}

	if !allowUsed {
		f.Allow["*"] = struct{}{}
	}
}

func printMethods(methods map[string]struct{}) string {
	var keys []string
	for k := range methods {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}
