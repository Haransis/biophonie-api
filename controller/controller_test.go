package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/haran/biophonie-api/controller/geopoint"
	"github.com/haran/biophonie-api/controller/user"
	"golang.org/x/crypto/bcrypt"
)

var c *Controller
var r *gin.Engine

const (
	geoIdEnabled  = 1
	geoIdDisabled = 2
)

var adminUser user.User = user.User{
	Id:        1,
	Name:      "admin",
	Password:  "57aba9df-969f-4871-a095-e916d06ba38b",
	CreatedOn: time.Now().String(),
}

var adminPassword, _ = bcrypt.GenerateFromPassword([]byte(adminUser.Password), bcrypt.DefaultCost)

var validUsers []user.User
var validTokens []string

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	preparePublicDir()
	c = NewController()
	r = SetupRouter(c)
	c.clearDatabase()
	c.refreshGeoJson()
	validUsers = make([]user.User, 0)
	validUsers = append(validUsers, adminUser)
	validTokens = make([]string, 0)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestRefreshGeoJson(t *testing.T) {
	bytesGeoJson, err := ioutil.ReadFile(c.geoJsonPath)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	var geoJson geopoint.GeoJson
	if err := json.Unmarshal(bytesGeoJson, &geoJson); err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(geoJson.Features) != 0 {
		t.Error("test-related: geoJson was not cleared properly")
	}
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
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
		{user.AddUser{Name: "admin"}, http.StatusConflict},
		{user.AddUser{Name: "bob"}, http.StatusOK},
		{user.AddUser{Name: "bobdu42"}, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(test.AddUser)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewReader(body))

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)
		if w.Code == http.StatusOK { // checks that token is not hashed
			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, http.StatusOK, w.Code)
			if _, err := uuid.Parse(got.Password); err != nil {
				t.Error(err)
			}
			validUsers = append(validUsers, got)
		}

		if test.StatusCode == http.StatusOK { // checks that user is stored properly
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/%s", test.AddUser.Name), nil)
			r.ServeHTTP(w, req)

			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, test.AddUser.Name, got.Name)
			if _, err := uuid.Parse(got.Password); err == nil {
				t.Fail()
			}
		}
	}
}

func TestGetUser(t *testing.T) {

	tests := []struct {
		User        user.User
		RequestName string
		StatusCode  int
	}{
		{user.User{Name: "alice"}, "eve", http.StatusNotFound},
		{user.User{Name: "charles"}, "Charles", http.StatusNotFound},
		{user.User{Name: "bob"}, "bob", http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/user/%s", test.RequestName), nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK { // Check if user was created properly
			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.User.Name, got.Name)
			assert.Equal(t, got.Password, "")
		}
	}
}

func TestCreateToken(t *testing.T) {

	tests := []struct {
		AuthUser   user.AuthUser
		StatusCode int
	}{
		{user.AuthUser{Name: validUsers[0].Name, Password: "random"}, http.StatusBadRequest},
		{user.AuthUser{Name: validUsers[0].Name, Password: ""}, http.StatusBadRequest},
		{user.AuthUser{Name: "charles", Password: "9b768967-d491-4baa-a812-24ea8a9c274d"}, http.StatusNotFound},
		{user.AuthUser{Name: validUsers[0].Name, Password: validUsers[1].Password}, http.StatusUnauthorized},
		{user.AuthUser{Name: validUsers[0].Name, Password: validUsers[0].Password}, http.StatusOK},
		{user.AuthUser{Name: validUsers[1].Name, Password: validUsers[1].Password}, http.StatusOK},
		{user.AuthUser{Name: validUsers[2].Name, Password: validUsers[2].Password}, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(test.AuthUser)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/token", bytes.NewReader(body))

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if w.Code == http.StatusOK {
			token := w.Body.Bytes()
			validTokens = append(validTokens, string(token))
		}
	}
}

func TestPingAuthenticated(t *testing.T) {
	unvalidTokens := c.wrongToken()
	tests := []struct {
		JWT        string
		StatusCode int
	}{
		{validTokens[0], http.StatusOK},
		{validTokens[0], http.StatusOK},
		{unvalidTokens[0], http.StatusNotFound},
		{unvalidTokens[1], http.StatusUnauthorized},
		{unvalidTokens[2], http.StatusUnauthorized},
	}
	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/restricted/ping", nil)

		if test.JWT != "" {
			// Write authorization token in header
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))
		}

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)
	}

}

func TestCreateGeoPoint(t *testing.T) {
	tests := []struct {
		SoundPath   string
		PicturePath string
		GeoPoint    geopoint.AddGeoPoint
		JWT         string
		StatusCode  int
	}{
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, validTokens[0], http.StatusOK},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, validTokens[0], http.StatusOK},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Fo", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, validTokens[0], http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night very late at night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, validTokens[0], http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 100000001.0, Longitude: 1000000000.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, validTokens[0], http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now().Add(200000 * time.Hour), Amplitudes: newAmplitudes(100)}, validTokens[0], http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, validTokens[0], http.StatusBadRequest},
		{"../testgeopoint/merle.wav", "../main.go", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, validTokens[0], http.StatusBadRequest},
		{"../main.go", "../testgeopoint/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, validTokens[0], http.StatusBadRequest},
	}

	for _, test := range tests {
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
		req, err := buildFormData(values, "/api/v1/restricted/geopoint")
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)
	}
}

func TestEnableGeoPoint(t *testing.T) {
	tests := []struct {
		GeoId      int
		JWT        string
		StatusCode int
	}{
		{-1, validTokens[0], http.StatusBadRequest},
		{9999, validTokens[0], http.StatusNotFound},
		{geoIdEnabled, validTokens[1], http.StatusUnauthorized},
		{geoIdEnabled, validTokens[0], http.StatusOK},
		{geoIdEnabled, validTokens[0], http.StatusNotFound},
	}
	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/restricted/geopoint/%d/enable", test.GeoId), nil)

		if test.JWT != "" {
			// Write authorization token in header
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))
		}

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)
	}
}

func TestAppendGeoJson(t *testing.T) {
	bytesGeoJson, err := ioutil.ReadFile(c.geoJsonPath)
	if err != nil {
		t.Error(err)
	}

	var geoJson geopoint.GeoJson
	if err := json.Unmarshal(bytesGeoJson, &geoJson); err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, len(geoJson.Features), 0)
	assert.Equal(t, geoJson.Features[0].Properties.Name, "Forest by night")
}

func TestGetGeoPoint(t *testing.T) {
	tests := []struct {
		GeoPoint   geopoint.GeoPoint
		StatusCode int
	}{
		{geopoint.GeoPoint{Id: 9999, Title: "Forest by night"}, http.StatusNotFound},
		{geopoint.GeoPoint{Id: geoIdDisabled, Title: "Forest by night"}, http.StatusForbidden},
		{geopoint.GeoPoint{Id: geoIdEnabled, Title: "Forest by night"}, http.StatusOK},
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

func TestGetRestrictedGeoPoint(t *testing.T) {
	tests := []struct {
		GeoPoint   geopoint.GeoPoint
		JWT        string
		StatusCode int
	}{
		{geopoint.GeoPoint{Id: 9999, Title: "Forest by night"}, validTokens[1], http.StatusUnauthorized},
		{geopoint.GeoPoint{Id: geoIdDisabled, Title: "Forest by night"}, validTokens[1], http.StatusUnauthorized},
		{geopoint.GeoPoint{Id: geoIdEnabled, Title: "Forest by night"}, validTokens[1], http.StatusUnauthorized},
		{geopoint.GeoPoint{Id: 9999, Title: "Forest by night"}, validTokens[0], http.StatusNotFound},
		{geopoint.GeoPoint{Id: geoIdDisabled, Title: "Forest by night"}, validTokens[0], http.StatusOK},
		{geopoint.GeoPoint{Id: geoIdEnabled, Title: "Forest by night"}, validTokens[0], http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/restricted/geopoint/%d", test.GeoPoint.Id), nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))
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

func TestMakeAdmin(t *testing.T) {
	tests := []struct {
		Id         int
		JWT        string
		StatusCode int
	}{
		{validUsers[1].Id, validTokens[1], http.StatusUnauthorized},
		{99999, validTokens[0], http.StatusNotFound},
		{validUsers[1].Id, validTokens[0], http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/restricted/user/%d", test.Id), nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))
		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/"+validUsers[1].Name, nil)
			r.ServeHTTP(w, req)
			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, true, got.Admin)
		}
	}
}

func (c *Controller) clearDatabase() {
	tx := c.Db.MustBegin()
	tx.MustExec("TRUNCATE TABLE accounts RESTART IDENTITY")
	tx.MustExec("TRUNCATE TABLE geopoints RESTART IDENTITY")
	tx.MustExec("INSERT INTO accounts (name, created_on, password, admin) VALUES ($1,now(),$2,'t') ON CONFLICT DO NOTHING", "admin", adminPassword)
	tx.Commit()
}

func preparePublicDir() {
	os.RemoveAll("/tmp/public")
	os.MkdirAll("/tmp/public/picture", os.ModePerm)
	os.MkdirAll("/tmp/public/sound", os.ModePerm)
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

func (c *Controller) wrongToken() []string {
	tokens := make([]string, 0)

	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &CustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
		},
		UserInfo{"random", true},
	}

	token, err := t.SignedString(c.signKey)
	if err != nil {
		panic(err)
	}
	tokens = append(tokens, token)

	t.Claims = &CustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now()),
		},
		UserInfo{validUsers[0].Name, false},
	}

	token, err = t.SignedString(c.signKey)
	if err != nil {
		panic(err)
	}

	tokens = append(tokens, token)
	tokens = append(tokens, "")
	return tokens
}
