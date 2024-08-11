package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func SHA256Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the signature from the header
			signature := c.Request().Header.Get("X-Hub-Signature-256")

			// Remove the prefix (if any)
			prefix := "sha256="
			signature = strings.TrimPrefix(signature, prefix)

			// Calculate the expected signature
			mac := hmac.New(sha256.New, []byte("your-secret"))

			// Read the body into a byte slice
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed_to_read_body")
			}

			mac.Write(bodyBytes)
			expectedSignature := hex.EncodeToString(mac.Sum(nil))

			// Compare the signatures
			if subtle.ConstantTimeCompare([]byte(signature), []byte(expectedSignature)) != 1 {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid_signature")
			}

			return next(c)
		}
	}
}
