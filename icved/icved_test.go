package main

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const crypted = "敬业民主民主法治敬业民主民主爱国敬业民主文明富强"
const plain = "abc"

type Case struct {
	t *testing.T
}

func (c *Case) Encrypt(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(c.t, crypted, r.Body.String())
	assert.Equal(c.t, http.StatusOK, r.Code)
}

func (c *Case) Decrypt(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	assert.Equal(c.t, plain, r.Body.String())
	assert.Equal(c.t, http.StatusOK, r.Code)
}

func TestHandle(t *testing.T) {
	r := gofight.New()
	e := Engine()
	var c Case
	c.t = t
	r.GET("/e/"+plain).SetDebug(true).Run(e, c.Encrypt)
	r.GET("/d/"+crypted).SetDebug(true).Run(e, c.Decrypt)
}
