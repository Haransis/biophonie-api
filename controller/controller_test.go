package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/haran/biophonie-api/controller/user"
)

var c *Controller
var r *gin.Engine

func (c *Controller) clearDatabase() {
	tx := c.Db.MustBegin()
	tx.MustExec("TRUNCATE TABLE accounts")
	tx.MustExec("TRUNCATE TABLE geopoints")
	tx.Commit()
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	c = NewController()
	r = SetupRouter(c)

	c.clearDatabase()
	exitVal := m.Run()
	c.clearDatabase()

	os.Exit(exitVal)
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestCreateUser(t *testing.T) {

	tests := []struct {
		AddUser    user.AddUser
		StatusCode int
	}{
		{user.AddUser{Name: "ev"}, http.StatusBadRequest},
		{user.AddUser{Name: "alalalalalalalalalalalalalalalal"}, http.StatusBadRequest},
		{user.AddUser{Name: "bob"}, http.StatusOK},
		{user.AddUser{Name: "bob"}, http.StatusConflict},
		{user.AddUser{Name: "bobdu42"}, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		userBytes, _ := json.Marshal(test.AddUser)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewReader(userBytes))

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/%s", test.AddUser.Name), nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestGetUser(t *testing.T) {

	tests := []struct {
		User        *user.User
		RequestName string
		StatusCode  int
	}{
		{newUser("alice"), "eve", http.StatusNotFound},
		{newUser("charles"), "Charles", http.StatusNotFound},
		{newUser("bob"), "bob", http.StatusOK},
		{newUser("bobdu42"), "bobdu42", http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/%s", test.RequestName), nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, *test.User.Name, got.Name)
		}
	}
}

func newUser(name string) *user.User {
	return &user.User{
		Name: &name,
	}
}
