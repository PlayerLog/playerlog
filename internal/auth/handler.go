package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/PlayerLog/playerlog/cmd/web/components"
	"github.com/PlayerLog/playerlog/cmd/web/views/auth_views"
	"github.com/PlayerLog/playerlog/internal/database"
	"github.com/PlayerLog/playerlog/internal/types"
	"github.com/PlayerLog/playerlog/pkg/etag"
	"github.com/PlayerLog/playerlog/pkg/jsonutil"
	"github.com/PlayerLog/playerlog/pkg/render"
	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
)

type Handler struct {
	service *Service
	store   *sessions.CookieStore
}

func NewHandler(db database.Service, sessionKey []byte) *Handler {
	return &Handler{
		service: NewService(db),
		store:   sessions.NewCookieStore(sessionKey),
	}
}

func (h Handler) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", h.RegisterUserHandler)
	mux.HandleFunc("POST /login", h.LoginUserHandler)
	mux.HandleFunc("GET /register-form", h.RegisterFormTypeHandler)
	mux.HandleFunc("GET /users", h.UsersHandler)
	mux.HandleFunc("GET /avatar", h.UserAvatarHandler)
	mux.HandleFunc("POST /user", h.UserHandler)

	return mux
}

func (h *Handler) LoginPage() *templ.ComponentHandler {
	return render.WrapHandler(auth_views.LoginPage())
}

func (h *Handler) RegisterPage() *templ.ComponentHandler {
	return render.WrapHandler(auth_views.RegisterPage())
}

func (h *Handler) RegisterFormTypeHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	values := r.Form

	var formType types.RegisterFormType

	err = form.NewDecoder().Decode(&formType, values)

	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	switch formType.Type {
	case "team":
		render.Render(w, auth_views.RegisterTeamForm(types.RegisterTeamValues{}, map[string]string{}))
	case "organization":
		render.Render(w, auth_views.OrganizationForm(types.RegisterTeamValues{}, map[string]string{}))
	default:
		render.Render(w, auth_views.RegisterTeamForm(types.RegisterTeamValues{}, map[string]string{}))
	}
	return
}

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(time.Second * 1)

	if r.Header.Get("Hx-Request") != "true" {
		http.Error(w, "Only HTMX Request", http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	values := r.Form

	var registerValues types.RegisterTeamValues

	err = form.NewDecoder().Decode(&registerValues, values)

	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	errs := registerValues.Validate()

	if len(errs) > 0 {
		render.Render(w, auth_views.RegisterTeamForm(registerValues, errs))
		return
	}

	err = h.service.CreateAccount(registerValues)

	if err != nil {
		errs["global"] = err.Error()
		// http.Error(w, err.Error(), http.StatusBadRequest)
		render.Render(w, auth_views.RegisterTeamForm(registerValues, errs))
		return
	}

	w.Header().Set("Hx-Redirect", "/dashboard")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	values := r.Form

	var loginValues types.LoginUserValues

	err = form.NewDecoder().Decode(&loginValues, values)

	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	errs := loginValues.Validate()

	if len(errs) > 0 {
		render.Render(w, auth_views.LoginForm(loginValues, errs))
		return
	}

	user, err := h.service.ValidateUser(loginValues.Email, loginValues.Password)

	if err != nil {
		errs["global"] = err.Error()
		render.Render(w, auth_views.LoginForm(loginValues, errs))
		return
	}

	session, _ := h.store.New(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Values["avatar_url"] = "https://i.pravatar.cc/150?img=55"

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,
		Secure:   true, // Enable in production
		SameSite: http.SameSiteStrictMode,
	}

	if err := session.Save(r, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Hx-Redirect", "/dashboard")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()

	if err != nil {
		return
	}

	jsonutil.WriteJson(w, http.StatusOK, users)
	return
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.ValidateUser("gurraarre5@gmail.com", "Gurrag090122")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonutil.WriteJson(w, http.StatusOK, user)
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.store.Get(r, "auth-session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// // Optional: Get fresh user data if needed
		// if userID, ok := session.Values["user_id"].(int64); ok {
		//     user, err := h.authService.GetUserByID(r.Context(), userID)
		//     if err != nil {
		//         session.Options.MaxAge = -1
		//         session.Save(r, w)
		//         http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		//         return
		//     }
		//     // Store user in request context if needed
		//     ctx := context.WithValue(r.Context(), "user", user)
		//     r = r.WithContext(ctx)
		// }

		next.ServeHTTP(w, r)
	})
}

// func (h *Handler) OpenAuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Default to unauthenticated
// 		ctx := context.WithValue(r.Context(), "authenticated", false)

// 		// Only attempt to get session if store exists
// 		if h.store != nil {
// 			if session, err := h.store.Get(r, "auth-session"); err == nil {
// 				// Type assert with single operation
// 				if auth, ok := session.Values["authenticated"].(bool); ok {
// 					ctx = context.WithValue(r.Context(), "authenticated", auth)
// 				}
// 			}
// 		}

// 		// Use the context directly without additional WithContext call
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func (h *Handler) OpenAuthMiddleware(next http.Handler) http.Handler {
	// Create a sync.Pool to reuse session objects
	sessionPool := &sync.Pool{
		New: func() interface{} {
			return &sessions.Session{}
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fast path for HTMX requests - check if we can reuse the auth state
		if hx := r.Header.Get("HX-Request"); hx != "" {
			// Check for existing auth state in header
			if authHeader := r.Header.Get("X-Auth-State"); authHeader != "" {
				isAuth := authHeader == "true"
				ctx := context.WithValue(r.Context(), "authenticated", isAuth)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Get session from pool
		session := sessionPool.Get().(*sessions.Session)
		defer sessionPool.Put(session)

		var isAuthenticated bool

		// Only attempt to get session if store exists
		if h.store != nil {
			if sess, err := h.store.Get(r, "auth-session"); err == nil {
				if auth, ok := sess.Values["authenticated"].(bool); ok {
					isAuthenticated = auth
				}
			}
		}

		// Set auth state in response header for future HTMX requests
		w.Header().Set("X-Auth-State", strconv.FormatBool(isAuthenticated))

		ctx := context.WithValue(r.Context(), "authenticated", isAuthenticated)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) UserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Standing on business")
	session, err := h.store.Get(r, "auth-session")

	if err != nil {
		render.Render(w, components.Profile(""))
		return
	}

	url, ok := session.Values["avatar_url"].(string)
	if !ok {
		render.Render(w, components.Profile(""))
		return
	}
	// Set Cache-Control headers to cache the avatar for 1 hour
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Set ETag to avoid unnecessary re-fetching
	etag := etag.CalculateETagForAvatar(url) // Replace with your method to calculate the ETag based on the avatar URL
	w.Header().Set("ETag", etag)

	// Handle conditional requests based on ETag
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	fmt.Println("Standing on business")

	render.Render(w, components.Profile(url))
	return
}
