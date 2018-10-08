package tools

import (
	"encoding/json"
	"net/http"
	"time"
)

type Client interface {
	GetUser(id string) (*User, error)
}

type client struct {
	http.Client
	baseURL string
}

type User struct {
	name  string
	email string
}

func New(url string) Client {
	return &client{
		http.Client{
			Timeout: time.Duration(1) * time.Second,
		},
		url,
	}
}

func (c *client) GetUser(id string) (*User, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/user/"+id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var user *User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
