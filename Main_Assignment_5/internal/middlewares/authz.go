package middlewares

import (
	"aaa_bekend/internal/lmodel"
	"net/http"
)

type RoleGuard struct {
	next         http.Handler
	allowedRoles []string
}

func (rg *RoleGuard) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, err := lmodel.UserFromContext(r.Context())
	// TODO
	if err != nil {
		htttp.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAllowed := false
	for _, role := range rg.allowedRoles {
		if user.Role == role {
			isAllowed = true
			break
		}
	}

	if isAllowed {
		rg.next.ServerHTTP(w, r)
		return
	}
	
	http.Error(w, "Forbidden", http.StatusForbidden)

}

func RequireRole(allowedRules ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &RoleGuard{
			next:         next,
			allowedRoles: allowedRules,
		}
	}
}
