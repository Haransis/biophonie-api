package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

	exitVal := m.Run()
	//c.clearDatabase()

	os.Exit(exitVal)
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestGetUser(t *testing.T) {
	c.clearDatabase()

	tests := []struct {
		User        *user.User
		RequestName string
		StatusCode  int
	}{
		{newUser(0, "alice", "2024-05-26T11:17:35.079344Z"), "eve", http.StatusNotFound},
		{newUser(1, "charles", "2024-05-26T11:17:35.079344Z"), "Charles", http.StatusNotFound},
		{newUser(2, "bob", "2022-05-26T11:17:35.079344Z"), "bob", http.StatusOK},
		{newUser(3, "bobdu42", "2022-05-26T11:17:35.079344Z"), "bobdu42", http.StatusOK},
	}

	for _, test := range tests {
		_, err := c.Db.NamedExec("INSERT INTO accounts (id, name, created_on) VALUES (:id, :name,:created_on) ON CONFLICT DO NOTHING", test.User)
		if err != nil {
			log.Fatal(err)
		}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/%s", test.RequestName), nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, *test.User, got)
		}
	}

}

func TestCreateUser(t *testing.T) {
	c.clearDatabase()

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

func newUser(id int, name string, date string) *user.User {
	return &user.User{
		Id:        &id,
		Name:      &name,
		CreatedOn: &date,
	}
}
