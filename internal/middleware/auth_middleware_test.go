package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func generateTestToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))
	return signed
}

func setupMiddlewareRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", AuthMiddleware(testSecret), func(c *gin.Context) {
		hospitalID, _ := c.Get("hospital_id")
		c.JSON(http.StatusOK, gin.H{"hospital_id": hospitalID})
	})
	return r
}

func TestMiddleware_NoToken(t *testing.T) {
	r := setupMiddlewareRouter()

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	r := setupMiddlewareRouter()

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	r := setupMiddlewareRouter()

	token := generateTestToken(jwt.MapClaims{
		"sub":         1,
		"hospital_id": float64(1),
		"exp":         time.Now().Add(-time.Hour).Unix(), // expired
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", w.Code)
	}
}

func TestMiddleware_ValidToken(t *testing.T) {
	r := setupMiddlewareRouter()

	token := generateTestToken(jwt.MapClaims{
		"sub":         1,
		"hospital_id": float64(1),
		"exp":         time.Now().Add(time.Hour).Unix(),
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMiddleware_WrongSecret(t *testing.T) {
	r := setupMiddlewareRouter()

	// Sign with wrong secret
	wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         1,
		"hospital_id": float64(1),
		"exp":         time.Now().Add(time.Hour).Unix(),
	})
	signed, _ := wrongToken.SignedString([]byte("wrong-secret"))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong secret, got %d", w.Code)
	}
}
