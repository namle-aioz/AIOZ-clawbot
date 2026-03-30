package session

import (
	"crypto/sha256"
	"net/http"
	"time"

	"backend/utils/response"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	sessionCookieName = "signup_session_id"
	EmailKey          = "email"
)

type Store struct {
	store *sessions.CookieStore
	name  string
	ttl   time.Duration
}

func NewStore(ttl time.Duration, secret string) *Store {
	key := deriveKey(secret)

	cookieStore := sessions.NewCookieStore(key)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	}

	return &Store{
		store: cookieStore,
		name:  sessionCookieName,
		ttl:   ttl,
	}
}

func (s *Store) Set(ctx echo.Context, key, value string) error {
	if key == "" {
		return response.NewHttpError(nil, response.ErrInvalidRequestBody, http.StatusBadRequest)
	}

	session, err := s.loadSession(ctx)
	if err != nil {
		return err
	}

	SetSessionValue(session, key, value)
	return s.saveSession(ctx, session, int(s.ttl.Seconds()))
}

func (s *Store) Get(ctx echo.Context, key string) (string, error) {
	if key == "" {
		return "", response.NewHttpError(nil, response.ErrInvalidRequestBody, http.StatusBadRequest)
	}

	session, err := s.loadSession(ctx)
	if err != nil {
		return "", err
	}

	value, ok := GetSessionValue[string](session, key)
	if !ok || value == "" {
		return "", response.NewHttpError(nil, response.ErrOTPSessionNotFound, http.StatusBadRequest)
	}

	return value, nil
}

func (s *Store) Clear(ctx echo.Context) {
	session, err := s.loadSession(ctx)
	if err != nil {
		session = sessions.NewSession(s.store, s.name)
	}

	session.Values = map[interface{}]interface{}{}
	if err := s.saveSession(ctx, session, -1); err != nil {
		_ = err
	}
}

func (s *Store) loadSession(ctx echo.Context) (*sessions.Session, error) {
	session, err := s.store.Get(ctx.Request(), s.name)
	if err != nil {
		return nil, response.NewHttpError(err, response.ErrOTPSessionNotFound, http.StatusBadRequest)
	}

	return session, nil
}

func (s *Store) saveSession(ctx echo.Context, session *sessions.Session, maxAge int) error {
	session.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}

	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return response.NewInternalError(err)
	}

	return nil
}

func SetSessionValue[T any](session *sessions.Session, key string, value T) {
	if session == nil {
		return
	}

	if session.Values == nil {
		session.Values = make(map[interface{}]interface{})
	}

	session.Values[key] = value
}

func GetSessionValue[T any](session *sessions.Session, key string) (T, bool) {
	var zero T
	if session == nil {
		return zero, false
	}

	rawValue, ok := session.Values[key]
	if !ok {
		return zero, false
	}

	value, ok := rawValue.(T)
	if !ok {
		return zero, false
	}

	return value, true
}

func deriveKey(secret string) []byte {
	sum := sha256.Sum256([]byte(secret))
	key := make([]byte, len(sum))
	copy(key, sum[:])

	return key
}
