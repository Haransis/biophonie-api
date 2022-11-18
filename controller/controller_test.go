package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cridenour/go-postgis"
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

var adminUser user.User = user.User{
	Id:       1,
	Name:     "admin",
	Password: "57aba9df-969f-4871-a095-e916d06ba38b",
	Admin:    true,
}
var standardUser user.User = user.User{
	Id:       2,
	Name:     "alice",
	Password: "57aca9df-969f-4861-a095-e916d06ba38b",
	Admin:    false,
}
var adminToken string
var standardToken string

var availableGeoPoint1 geopoint.DbGeoPoint = geopoint.DbGeoPoint{
	GeoPoint: &geopoint.GeoPoint{
		Id:         1,
		Title:      "Enabled",
		UserId:     1,
		Latitude:   1.02,
		Longitude:  1.0,
		CreatedOn:  time.Now(),
		Amplitudes: newAmplitudes(500),
		Picture:    "forest",
		Sound:      "sound1.wav",
		Available:  true,
	},
	Location: postgis.PointS{SRID: geopoint.WGS84, X: 1.02, Y: 1.0},
}
var availableGeoPoint2 geopoint.DbGeoPoint = geopoint.DbGeoPoint{
	GeoPoint: &geopoint.GeoPoint{
		Id:         2,
		Title:      "Enabled",
		UserId:     2,
		Latitude:   2.02,
		Longitude:  2.0,
		CreatedOn:  time.Now(),
		Amplitudes: newAmplitudes(500),
		Picture:    "forest",
		Sound:      "sound2.wav",
		Available:  true,
	},
	Location: postgis.PointS{SRID: geopoint.WGS84, X: 2.02, Y: 2.0},
}
var unavailableGeoPoint geopoint.DbGeoPoint = geopoint.DbGeoPoint{
	GeoPoint: &geopoint.GeoPoint{
		Id:         3,
		Title:      "Disabled",
		UserId:     1,
		Latitude:   3.02,
		Longitude:  3.0,
		CreatedOn:  time.Now(),
		Amplitudes: newAmplitudes(500),
		Picture:    "sea",
		Sound:      "sound4.wav",
		Available:  false,
	},
	Location: postgis.PointS{SRID: geopoint.WGS84, X: 3.02, Y: 3.0},
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	preparePublicDir()
	c = NewController()
	r = SetupRouter(c)
	c.clearDatabase()
	c.refreshGeoJson()
	c.createTokens()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestRefreshGeoJson(t *testing.T) {
	bytesGeoJson, err := ioutil.ReadFile(c.geoJsonPath)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	var geoJson geoJson
	if err := json.Unmarshal(bytesGeoJson, &geoJson); err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestPostUser(t *testing.T) {

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
		body, _ := json.Marshal(test.AddUser)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewReader(body))

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got user.User
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Errorf("response is not a user: %s", err)
			}
			if _, err := uuid.Parse(got.Password); err != nil {
				t.Errorf("password is not uuid: %s", err)
			}

			var created user.User
			if err := c.Db.Get(&created, "SELECT * FROM accounts WHERE name = $1", test.AddUser.Name); err != nil {
				t.Errorf("user was not stored: %s", err)
			}
			assert.Equal(t, test.AddUser.Name, got.Name)

			if err := bcrypt.CompareHashAndPassword([]byte(created.Password), []byte(got.Password)); err != nil {
				t.Errorf("password was not hashed properly: %s", err)
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
		{user.User{Name: "alice"}, "alice", http.StatusOK},
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
			assert.Equal(t, test.User.Name, got.Name)
			assert.Equal(t, got.Password, "") // password is not divulgated
		}
	}
}

func TestAuthorizeUser(t *testing.T) {

	tests := []struct {
		AuthUser   user.AuthUser
		StatusCode int
	}{
		{user.AuthUser{Name: standardUser.Name, Password: "random"}, http.StatusBadRequest},
		{user.AuthUser{Name: standardUser.Name, Password: ""}, http.StatusBadRequest},
		{user.AuthUser{Name: "charles", Password: "9b768967-d491-4baa-a812-24ea8a9c274d"}, http.StatusNotFound},
		{user.AuthUser{Name: standardUser.Name, Password: adminUser.Password}, http.StatusUnauthorized},
		{user.AuthUser{Name: standardUser.Name, Password: standardUser.Password}, http.StatusOK},
		{user.AuthUser{Name: adminUser.Name, Password: adminUser.Password}, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(test.AuthUser)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/authorize", bytes.NewReader(body))

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if w.Code == http.StatusOK {
			var got user.AccessToken
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			_, err := jwt.Parse(got.Token, func(token *jwt.Token) (interface{}, error) {
				return c.verifyKey, nil
			})
			if err != nil {
				t.Errorf("could not parse returned token: %s", err)
			}
		}
	}
}

func TestPingAuthenticated(t *testing.T) {
	unvalidTokens := c.wrongToken()
	tests := []struct {
		JWT        string
		StatusCode int
	}{
		{standardToken, http.StatusOK},
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

func TestPostGeoPoint(t *testing.T) {
	tests := []struct {
		SoundPath   string
		PicturePath string
		GeoPoint    geopoint.AddGeoPoint
		JWT         string
		StatusCode  int
	}{
		{"../testassets/merle.wav", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.1, Date: time.Now(), Amplitudes: newAmplitudes(100)}, standardToken, http.StatusOK},
		{"../testassets/merle.wav", "", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100), PictureTemplate: "clearing"}, standardToken, http.StatusOK},
		{"../testassets/merle.wav", "", geopoint.AddGeoPoint{Title: "Mountain by day", Latitude: 1.0, Longitude: 1.3, Date: time.Now(), Amplitudes: newAmplitudes(100), PictureTemplate: "desert"}, standardToken, http.StatusOK},
		{"../testassets/merle.wav", "", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, standardToken, http.StatusBadRequest},
		{"../testassets/merle.wav", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Fo", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, standardToken, http.StatusBadRequest},
		{"../testassets/merle.wav", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night very late at night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, standardToken, http.StatusBadRequest},
		{"../testassets/merle.wav", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 100000001.0, Longitude: 1000000000.2, Date: time.Now(), Amplitudes: newAmplitudes(100)}, standardToken, http.StatusBadRequest},
		{"../testassets/merle.wav", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now().Add(200000 * time.Hour), Amplitudes: newAmplitudes(100)}, standardToken, http.StatusBadRequest},
		{"../testassets/merle.wav", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, standardToken, http.StatusBadRequest},
		{"../testassets/merle.wav", "../main.go", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, standardToken, http.StatusBadRequest},
		{"../main.go", "../testassets/russie.jpg", geopoint.AddGeoPoint{Title: "Forest by night", Latitude: 1.0, Longitude: 1.2, Date: time.Now(), Amplitudes: newAmplitudes(1)}, standardToken, http.StatusBadRequest},
	}

	for _, test := range tests {
		userBytes, err := json.Marshal(test.GeoPoint)
		if err != nil {
			t.Error(err)
		}
		values := map[string]io.Reader{
			"sound":    mustOpen(test.SoundPath),
			"geopoint": strings.NewReader(string(userBytes)),
		}
		if test.PicturePath != "" {
			values["picture"] = mustOpen(test.PicturePath)
		}

		w := httptest.NewRecorder()
		req, err := buildFormData(values, "/api/v1/restricted/geopoint")
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))

		r.ServeHTTP(w, req)
		fmt.Println(w.Body.String())
		assert.Equal(t, test.StatusCode, w.Code)
	}
}

func TestEnableGeoPoint(t *testing.T) {
	defer c.Db.MustExec("UPDATE geopoints SET available = FALSE WHERE id = $1", unavailableGeoPoint.Id)

	tests := []struct {
		GeoId      int
		JWT        string
		StatusCode int
	}{
		{-1, adminToken, http.StatusBadRequest},
		{9999, adminToken, http.StatusNotFound},
		{unavailableGeoPoint.Id, standardToken, http.StatusUnauthorized},
		{unavailableGeoPoint.Id, adminToken, http.StatusOK},
		{unavailableGeoPoint.Id, adminToken, http.StatusNotFound},
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

	bytesGeoJson, err := ioutil.ReadFile(c.geoJsonPath)
	if err != nil {
		t.Errorf("cannot read geojson file: %s", err)
	}

	var geoJson geoJson
	if err := json.Unmarshal(bytesGeoJson, &geoJson); err != nil {
		t.Errorf("geojson was not properly refreshed: %s", err)
	}

	assert.NotEqual(t, len(geoJson.Features), 0)
	assert.Equal(t, geoJson.Features[unavailableGeoPoint.Id-1].Properties.Name, unavailableGeoPoint.Title)
}

func TestGetGeoPoint(t *testing.T) {
	tests := []struct {
		GeoPoint   geopoint.GeoPoint
		StatusCode int
	}{
		{geopoint.GeoPoint{Id: 9999, Title: "Random"}, http.StatusNotFound},
		{*unavailableGeoPoint.GeoPoint, http.StatusForbidden},
		{*availableGeoPoint1.GeoPoint, http.StatusOK},
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

func TestGetAssets(t *testing.T) {
	tests := []struct {
		GeoPoint   geopoint.GeoPoint
		StatusCode int
	}{
		{geopoint.GeoPoint{Id: 9999, Title: "Random"}, http.StatusNotFound},
		{*unavailableGeoPoint.GeoPoint, http.StatusForbidden},
		{*availableGeoPoint1.GeoPoint, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/geopoint/%d/assets", test.GeoPoint.Id), nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got geopoint.Assets
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.GeoPoint.Picture, got.Picture)
			assert.Equal(t, test.GeoPoint.Sound, got.Sound)
		}
	}
}

func TestGetClosestGeoPoint(t *testing.T) {
	tests := []struct {
		Latitude   float64
		Longitude  float64
		Not        []string
		StatusCode int
		IdResult   int
	}{
		{1, -10000000000, []string{}, http.StatusBadRequest, 0},
		{-10000, 1.0, []string{}, http.StatusBadRequest, 0},
		{1.0, 1.1, []string{}, http.StatusOK, availableGeoPoint1.Id},
		{1.0, 1.1, []string{fmt.Sprint(availableGeoPoint1.Id)}, http.StatusOK, availableGeoPoint2.Id},
		{1.0, 1.1, []string{fmt.Sprint(availableGeoPoint1.Id), fmt.Sprint(availableGeoPoint2.Id)}, http.StatusNotFound, 0},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()

		query := fmt.Sprintf("/api/v1/geopoint/closest/to/%f/%f", test.Latitude, test.Longitude)
		values := url.Values{"not[]": test.Not}
		query += "?" + values.Encode()
		req, _ := http.NewRequest(http.MethodGet, query, nil)

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got struct{ Id int }
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.IdResult, got.Id)
		}
	}
}

func TestGetRestrictedGeoPoint(t *testing.T) {
	tests := []struct {
		GeoPoint   geopoint.GeoPoint
		JWT        string
		StatusCode int
	}{
		{geopoint.GeoPoint{Id: 9999, Title: "Forest by night"}, standardToken, http.StatusUnauthorized},
		{*unavailableGeoPoint.GeoPoint, standardToken, http.StatusUnauthorized},
		{geopoint.GeoPoint{Id: 9999, Title: "Forest by night"}, adminToken, http.StatusNotFound},
		{*unavailableGeoPoint.GeoPoint, adminToken, http.StatusOK},
		{*availableGeoPoint1.GeoPoint, adminToken, http.StatusOK},
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

func TestDeleteGeoPoint(t *testing.T) {
	tests := []struct {
		IdToDelete int
		JWT        string
		StatusCode int
	}{
		{availableGeoPoint1.Id, "", http.StatusUnauthorized},
		{availableGeoPoint1.Id, standardToken, http.StatusUnauthorized},
		{12000, adminToken, http.StatusNotFound},
		{availableGeoPoint1.Id, adminToken, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/restricted/geopoint/%d", test.IdToDelete), nil)

		if test.JWT != "" {
			// Write authorization token in header
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))
		}

		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var geoPoint geopoint.DbGeoPoint
			err := c.Db.Get(&geoPoint, "SELECT * FROM geopoints WHERE id = $1", test.IdToDelete)
			assert.Equal(t, err, sql.ErrNoRows)
		}
	}
}

func TestMakeAdmin(t *testing.T) {
	defer c.Db.Exec("UPDATE accounts SET admin = TRUE WHERE id = $1", standardUser.Id)

	tests := []struct {
		Id         int
		JWT        string
		StatusCode int
	}{
		{standardUser.Id, standardToken, http.StatusUnauthorized},
		{99999, adminToken, http.StatusNotFound},
		{standardUser.Id, adminToken, http.StatusOK},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/restricted/user/%d", test.Id), nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.JWT))
		r.ServeHTTP(w, req)
		assert.Equal(t, test.StatusCode, w.Code)

		if test.StatusCode == http.StatusOK {
			var got user.User
			if err := c.Db.Get(&got, "SELECT * FROM accounts WHERE id = $1", test.Id); err != nil {
				t.Errorf("admined user not in database: %s", err)
			}
			assert.Equal(t, true, got.Admin)
		}
	}
}

func (c *Controller) clearDatabase() {
	var hashAdminPwd, _ = bcrypt.GenerateFromPassword([]byte(adminUser.Password), bcrypt.DefaultCost)
	var hashAlicePwd, _ = bcrypt.GenerateFromPassword([]byte(standardUser.Password), bcrypt.DefaultCost)
	tx := c.Db.MustBegin()
	tx.MustExec("TRUNCATE TABLE accounts RESTART IDENTITY")
	tx.MustExec("TRUNCATE TABLE geopoints RESTART IDENTITY")
	tx.MustExec("INSERT INTO accounts (name, created_on, password, admin) VALUES ($1,now(),$2,$3) ON CONFLICT DO NOTHING", adminUser.Name, hashAdminPwd, adminUser.Admin)
	tx.MustExec("INSERT INTO accounts (name, created_on, password, admin) VALUES ($1,now(),$2,$3) ON CONFLICT DO NOTHING", standardUser.Name, hashAlicePwd, adminUser.Admin)
	tx.NamedExec("INSERT INTO geopoints (title, user_id, location, amplitudes, picture, sound, created_on, available) VALUES (:title,:user_id,GeomFromEWKB(:location),:amplitudes,:picture,:sound,:created_on,:available)", availableGeoPoint1)
	tx.NamedExec("INSERT INTO geopoints (title, user_id, location, amplitudes, picture, sound, created_on, available) VALUES (:title,:user_id,GeomFromEWKB(:location),:amplitudes,:picture,:sound,:created_on,:available)", availableGeoPoint2)
	tx.NamedExec("INSERT INTO geopoints (title, user_id, location, amplitudes, picture, sound, created_on, available) VALUES (:title,:user_id,GeomFromEWKB(:location),:amplitudes,:picture,:sound,:created_on,:available)", unavailableGeoPoint)
	tx.Commit()
}

func (c *Controller) createTokens() {
	adminToken, _ = c.createToken(adminUser.Name, adminUser.Admin)
	standardToken, _ = c.createToken(standardUser.Name, standardUser.Admin)
}

func preparePublicDir() {
	os.RemoveAll("/tmp/public")
	os.MkdirAll("/tmp/public/assets/picture", os.ModePerm)
	os.MkdirAll("/tmp/public/assets/sound", os.ModePerm)
	os.Create("/tmp/public/assets/geojson.json")
}

func newAmplitudes(length int) []float64 {
	ampl := make([]float64, length)
	for i := 0; i < length; i++ {
		ampl[i] = rand.Float64() - rand.Float64()
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

	t.Claims = &user.CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
		},
		UserInfo: user.UserInfo{Name: "random", Admin: true},
	}

	token, err := t.SignedString(c.signKey)
	if err != nil {
		panic(err)
	}
	tokens = append(tokens, token)

	t.Claims = &user.CustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now()),
		},
		UserInfo: user.UserInfo{Name: adminUser.Password, Admin: false},
	}

	token, err = t.SignedString(c.signKey)
	if err != nil {
		panic(err)
	}

	tokens = append(tokens, token)
	tokens = append(tokens, "")
	return tokens
}

type geoJson struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

type feature struct {
	Type        string     `json:"type"`
	Coordinates []float64  `json:"coordinates"`
	Properties  properties `json:"properties"`
}

type properties struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
