package supabase

import "net/http"

type GoTrueUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type GoTrueSession struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int        `json:"expires_in"`
	TokenType    string     `json:"token_type"`
	User         GoTrueUser `json:"user"`
}

// SignUp creates a new user via the GoTrue admin endpoint with email confirmation skipped.
// User metadata is stored on auth.users; the application profile row is created separately.
func (c *Client) SignUp(email, password string, metadata map[string]any) (*GoTrueUser, error) {
	body := map[string]any{
		"email":    email,
		"password": password,
		"data":     metadata,
	}
	path := "/auth/v1/signup"
	useAnon := true
	if c.HasServiceKey() {
		delete(body, "data")
		body["email_confirm"] = true
		body["user_metadata"] = metadata
		path = "/auth/v1/admin/users"
		useAnon = false
	}
	resp, err := c.do(http.MethodPost, path, body, useAnon, "")
	if err != nil {
		return nil, err
	}
	var u GoTrueUser
	if err := decode(resp, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) PasswordSignIn(email, password string) (*GoTrueSession, error) {
	body := map[string]any{"email": email, "password": password}
	resp, err := c.do(http.MethodPost, "/auth/v1/token?grant_type=password", body, true, "")
	if err != nil {
		return nil, err
	}
	var s GoTrueSession
	if err := decode(resp, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) RefreshSession(refreshToken string) (*GoTrueSession, error) {
	body := map[string]any{"refresh_token": refreshToken}
	resp, err := c.do(http.MethodPost, "/auth/v1/token?grant_type=refresh_token", body, true, "")
	if err != nil {
		return nil, err
	}
	var s GoTrueSession
	if err := decode(resp, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// SignOut invalidates the user's refresh token server-side.
func (c *Client) SignOut(userToken string) error {
	resp, err := c.do(http.MethodPost, "/auth/v1/logout", nil, true, userToken)
	if err != nil {
		return err
	}
	return decode(resp, nil)
}
