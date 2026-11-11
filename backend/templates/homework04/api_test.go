package homework04_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"homework04/internal/database"
	"homework04/internal/handler"
	"homework04/internal/router"

	"github.com/gin-gonic/gin"
)

func testRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := database.OpenSQLiteMemory()
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := &handler.Handler{DB: db, Logger: logger}
	return router.New(h, logger)
}

func doJSON(t *testing.T, r http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestBlogAPIFlow(t *testing.T) {
	r := testRouter(t)

	w := doJSON(t, r, http.MethodPost, "/api/register", "", map[string]string{
		"username": "alice",
		"password": "secret1",
		"email":    "alice@example.com",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("register: %d %s", w.Code, w.Body.String())
	}

	w = doJSON(t, r, http.MethodPost, "/api/login", "", map[string]string{
		"username": "alice",
		"password": "secret1",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("login: %d %s", w.Code, w.Body.String())
	}
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil || loginResp.Token == "" {
		t.Fatalf("login token: %v body=%s", err, w.Body.String())
	}
	token := loginResp.Token

	w = doJSON(t, r, http.MethodPost, "/api/posts", token, map[string]string{
		"title":   "Hello",
		"content": "world",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create post: %d %s", w.Code, w.Body.String())
	}
	var createResp struct {
		Post struct {
			ID uint `json:"ID"`
		} `json:"post"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &createResp)
	postID := createResp.Post.ID
	if postID == 0 {
		t.Fatalf("expected post id, body=%s", w.Body.String())
	}

	w = doJSON(t, r, http.MethodGet, "/api/posts", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list posts: %d", w.Code)
	}

	w = doJSON(t, r, http.MethodGet, "/api/posts/1", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get post: %d", w.Code)
	}

	w = doJSON(t, r, http.MethodPost, "/api/posts/1/comments", token, map[string]string{
		"content": "nice post",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment: %d %s", w.Code, w.Body.String())
	}

	w = doJSON(t, r, http.MethodGet, "/api/posts/1/comments", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list comments: %d", w.Code)
	}

	w = doJSON(t, r, http.MethodPut, "/api/posts/1", token, map[string]string{
		"title":   "Updated",
		"content": "body",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("update post: %d %s", w.Code, w.Body.String())
	}

	w = doJSON(t, r, http.MethodPost, "/api/posts", "", map[string]string{
		"title":   "x",
		"content": "y",
	})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("create without auth: expected 401, got %d", w.Code)
	}

	w = doJSON(t, r, http.MethodDelete, "/api/posts/"+strconv.FormatUint(uint64(postID), 10), token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("delete post: %d %s", w.Code, w.Body.String())
	}
}

func TestUpdateForbiddenForOtherUser(t *testing.T) {
	r := testRouter(t)

	doJSON(t, r, http.MethodPost, "/api/register", "", map[string]string{
		"username": "author",
		"password": "secret12",
		"email":    "a@example.com",
	})
	doJSON(t, r, http.MethodPost, "/api/register", "", map[string]string{
		"username": "other",
		"password": "secret12",
		"email":    "b@example.com",
	})

	w := doJSON(t, r, http.MethodPost, "/api/login", "", map[string]string{
		"username": "author",
		"password": "secret12",
	})
	var loginResp struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &loginResp)
	authorToken := loginResp.Token

	doJSON(t, r, http.MethodPost, "/api/posts", authorToken, map[string]string{
		"title":   "mine",
		"content": "content",
	})

	w = doJSON(t, r, http.MethodPost, "/api/login", "", map[string]string{
		"username": "other",
		"password": "secret12",
	})
	_ = json.Unmarshal(w.Body.Bytes(), &loginResp)
	otherToken := loginResp.Token

	w = doJSON(t, r, http.MethodPut, "/api/posts/1", otherToken, map[string]string{
		"title":   "hack",
		"content": "hack",
	})
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
