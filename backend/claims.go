package main

type Claims struct {
	Subject           string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
}

func (c Claims) User() string {
	if c.PreferredUsername != "" {
		return c.PreferredUsername
	}
	if c.Email != "" {
		return c.Email
	}
	return c.Subject
}
