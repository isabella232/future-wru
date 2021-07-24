package wru

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/ymotongpoo/datemaki"
)

var ErrInvalidSessionToken = errors.New("invalid session token")

type SessionStatus int

const (
	BeforeLogin SessionStatus = iota
	ActiveSession
	IdleTimeoutSession
	AbsoluteTimeoutSession // This is not used
)

type SingleSessionData struct {
	ID             string    `docstore:"id" json:"-"`
	UserID         string    `docstore:"user_id" json:"-"`
	LoginAt        time.Time `docstore:"login_at" json:"login_at"`
	LastAccessAt   time.Time `docstore:"last_access_at" json:"last_access_at"`
	CurrentSession bool      `docstore:"-" json:"current"`

	LoginInfo map[string]string `docstore:"loginInfo" json:"login_info"`
}

func (s SingleSessionData) LoginAtFormat() string {
	return s.LoginAt.Format("2006/Jan/02 15:04")
}

func (s SingleSessionData) LoginAtForHuman() string {
	return datemaki.FormatDuration(time.Now().Sub(s.LoginAt))
}

func (s SingleSessionData) LastAccessAtFormat() string {
	return s.LastAccessAt.Format("2006/Jan/02 15:04")
}

func (s SingleSessionData) LastAccessAtForHuman() string {
	return datemaki.FormatDuration(time.Now().Sub(s.LastAccessAt))
}

func (s SingleSessionData) OS() string {
	return s.LoginInfo["os"]
}

func (s SingleSessionData) Browser() string {
	return s.LoginInfo["browser"]
}

type AllUserSessions []SingleSessionData

func (as AllUserSessions) WriteAsJson(w io.Writer) error {
	type Sessions struct {
		Sessions []SingleSessionData `json:"sessions"`
	}

	e := json.NewEncoder(w)
	return e.Encode(&Sessions{
		Sessions: as,
	})
}

type UserSession struct {
	ID       string            `docstore:"id"`
	Sessions []string          `docstore:"singleSessions"`
	Data     map[string]string `docstore:"data"`

	// User Informations
	DisplayName  string   `docstore:"name"`
	Email        string   `docstore:"email"`
	Organization string   `docstore:"org"`
	Scopes       []string `docstore:"scopes"`
}

type UnixTime time.Time

func (u UnixTime) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, time.Time(u).UnixNano(), 10), nil
}

type Session struct {
	LoginAt      UnixTime          `json:"login_at"`
	ExpireAt     UnixTime          `json:"expire_at"`
	LastAccessAt UnixTime          `json:"last_access_at"`
	UserID       string            `json:"id"`
	DisplayName  string            `json:"name"`
	Email        string            `json:"email"`
	Organization string            `json:"org"`
	Scopes       []string          `json:"scopes"`
	Status       SessionStatus     `json:"-"`
	Data         map[string]string `json:"data"`
}

type Directive struct {
	Key   string
	Value string
}

var directiveRe = regexp.MustCompile(`\s*(\S+)\s*=\s*(.*)`)

func parseDirective(src string) (*Directive, error) {
	match := directiveRe.FindStringSubmatch(src)
	if len(match) == 0 {
		return nil, fmt.Errorf("parse directive error: %s", src)
	}
	return &Directive{
		Key:   match[1],
		Value: match[2],
	}, nil
}

type SessionStorage interface {
	// StartLogin is called before login session
	// info keeps information like redirect URL
	StartLogin(ctx context.Context, info map[string]string) (sessionID string, err error)
	// AddLoginInfo adds extra login information for IDP.
	AddLoginInfo(ctx context.Context, oldSessionID string, info map[string]string) (newSessionID string, err error)
	// StartSessionAndRedirect is called after authorization and it renews login session ID and return info that is stored in StartLogin
	StartSession(ctx context.Context, oldSessionID string, user *User, r *http.Request) (newSessionID string, info map[string]string, err error)
	Logout(ctx context.Context, sessionID string) error
	GetUserSessions(ctx context.Context, userID string) ([]SingleSessionData, error)
	FindBySessionToken(ctx context.Context, sessionID string) (*Session, error)
	UpdateSessionData(ctx context.Context, sessionID string, directives []string) (err error)
	RenewSession(ctx context.Context, oldSessionID string) (sessionID string, err error)
}

func NewSessionStorage(ctx context.Context, c *Config, out io.Writer) (SessionStorage, error) {
	// todo: switch redis
	if c.SessionStorage == "" {
		return NewMemorySessionStorage(ctx, c, "")
	} else {
		return NewServerlessSessionStorage(ctx, c, "")
	}
}

type RedisSessionStorage struct {
}

func (s RedisSessionStorage) UpdateSessionData(ctx context.Context, sessionID string, directives []string) (err error) {
	panic("implement me")
}

func (s RedisSessionStorage) RenewSession(ctx context.Context, oldSessionID string) (sessionID string, err error) {
	panic("implement me")
}

func (s RedisSessionStorage) StartLogin(ctx context.Context, info map[string]string) (sessionID string, err error) {
	panic("implement me")
}

func (s RedisSessionStorage) AddLoginInfo(ctx context.Context, oldSessionID string, info map[string]string) (newSessionID string, err error) {
	panic("implement me")
}

func (s RedisSessionStorage) StartSession(ctx context.Context, oldSessionID string, user *User, r *http.Request) (sessionID string, info map[string]string, err error) {
	panic("implement me")
}

func (s RedisSessionStorage) Logout(ctx context.Context, sessionID string) error {
	panic("implement me")
}

func (s RedisSessionStorage) GetUserSessions(ctx context.Context, userID string) ([]SingleSessionData, error) {
	panic("implement me")
}

func (s RedisSessionStorage) FindBySessionToken(ctx context.Context, token string) (*Session, error) {
	panic("implement me")
}

var _ SessionStorage = &RedisSessionStorage{}