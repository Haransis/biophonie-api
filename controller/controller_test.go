package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/haran/biophonie-api/controller/geopoint"
	"github.com/haran/biophonie-api/controller/user"
)

var c *Controller
var r *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	c = NewController()
	r = SetupRouter(c)

	c.clearDatabase()
	c.preparePublicDir()
	exitVal := m.Run()

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

func TestCreateGeoPoint(t *testing.T) {
	tests := []struct {
		SoundPath   string
		PicturePath string
		GeoPoint    geopoint.AddGeoPoint
		StatusCode  int
	}{
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, http.StatusOK},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Fo", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night very late at night", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 9999, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, http.StatusNotFound},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 1, Latitude: 100000001.0, Longitude: 1000000000.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now().Add(200000 * time.Hour), Amplitudes: newAmplitudes(100)}, http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../main.go", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, http.StatusBadRequest},
		{"../main.go", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", UserId: 1, Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, http.StatusBadRequest},
	}

	for id, test := range tests {
		userBytes, _ := json.Marshal(test.GeoPoint)
		jsonFile, err := os.CreateTemp("", "*.json")
		if err != nil {
			t.Error(err)
		}
		jsonFile.Write(userBytes)
		jsonFile.Seek(0, 0)
		values := map[string]io.Reader{
			"sound":    mustOpen(test.SoundPath),
			"picture":  mustOpen(test.PicturePath),
			"geopoint": jsonFile,
		}

		w := httptest.NewRecorder()
		req, err := buildFormData(values, "/api/v1/geopoint")
		if err != nil {
			t.Error(err)
		}

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/geopoint/%d", id+1), nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestGetGeoPoint(t *testing.T) {
	tests := []struct {
		GeoPoint   geopoint.GeoPoint
		StatusCode int
	}{
		{geopoint.GeoPoint{Id: 9999, Title: "Forest by night"}, http.StatusNotFound},
		{geopoint.GeoPoint{Id: 1, Title: "Forest by night"}, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/geopoint/%d", test.GeoPoint.Id), nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got geopoint.GeoPoint
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.GeoPoint.Title, got.Title)
		}
	}
}

func (c *Controller) clearDatabase() {
	tx := c.Db.MustBegin()
	tx.MustExec("TRUNCATE TABLE accounts RESTART IDENTITY")
	tx.MustExec("TRUNCATE TABLE geopoints RESTART IDENTITY")
	tx.Commit()
}

func (c *Controller) preparePublicDir() {
	os.RemoveAll("/tmp/public")
	os.MkdirAll("/tmp/public/picture", os.ModePerm)
	os.MkdirAll("/tmp/public/sound", os.ModePerm)
}

func newUser(name string) *user.User {
	return &user.User{
		Name: &name,
	}
}

func newAmplitudes(length int) []int64 {
	ampl := make([]int64, length)
	for i := 0; i < length; i++ {
		ampl[i] = rand.Int63n(100) - rand.Int63n(100)
	}
	return ampl
}

func buildFormData(values map[string]io.Reader, url string) (*http.Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var err error
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return nil, err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	w.Close()
	req, _ := http.NewRequest(http.MethodPost, url, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
