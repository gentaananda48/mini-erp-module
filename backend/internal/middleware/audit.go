package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/internal/model"
)

const maxAuditBodySize = 64 * 1024

type AuditRepository interface {
    Create(ctx context.Context, audit model.AuditLog) error
}

func Audit(repo AuditRepository) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestBody := captureRequestBody(r)

            recorder := &responseRecorder{
                ResponseWriter: w,
            }

            next.ServeHTTP(recorder, r)

            audit := model.AuditLog{
                Method:         r.Method,
                Endpoint:       r.URL.Path,
                RequestBody:    requestBody,
                ResponseStatus: recorder.status,
                ResponseBody:   recorder.body(),
            }

            ctx, cancel := context.WithTimeout(
                context.Background(),
                2*time.Second,
            )
            defer cancel()

            if err := repo.Create(ctx, audit); err != nil {
                log.Printf("failed to save audit log: %v", err)
            }
        })
    }
}

func captureRequestBody(r *http.Request) string {
    if r.Body == nil || !isJSON(r.Header.Get("Content-Type")) {
        return ""
    }

    body, err := io.ReadAll(io.LimitReader(
        r.Body,
        maxAuditBodySize+1,
    ))
    if err != nil {
        return "[request body could not be read]"
    }

    r.Body = io.NopCloser(io.MultiReader(
        bytes.NewReader(body),
        r.Body,
    ))

    if len(body) > maxAuditBodySize {
        return "[request body omitted: exceeds audit limit]"
    }

    return sanitizeJSON(body)
}

type responseRecorder struct {
    http.ResponseWriter
    status     int
    bodyBuffer bytes.Buffer
    tooLarge   bool
}

func (r *responseRecorder) WriteHeader(status int) {
    if r.status != 0 {
        return
    }

    r.status = status
    r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(body []byte) (int, error) {
    if r.status == 0 {
        r.WriteHeader(http.StatusOK)
    }

    remaining := maxAuditBodySize - r.bodyBuffer.Len()

    if remaining > 0 {
        if len(body) <= remaining {
            _, _ = r.bodyBuffer.Write(body)
        } else {
            _, _ = r.bodyBuffer.Write(body[:remaining])
            r.tooLarge = true
        }
    } else if len(body) > 0 {
        r.tooLarge = true
    }

    return r.ResponseWriter.Write(body)
}

func (r *responseRecorder) Unwrap() http.ResponseWriter {
    return r.ResponseWriter
}

func (r *responseRecorder) body() string {
    if r.tooLarge {
        return "[response body omitted: exceeds audit limit]"
    }

    if !isJSON(r.Header().Get("Content-Type")) {
        return ""
    }

    return sanitizeJSON(r.bodyBuffer.Bytes())
}

func isJSON(contentType string) bool {
    contentType = strings.ToLower(contentType)

    return strings.Contains(contentType, "application/json") ||
        strings.Contains(contentType, "+json")
}

func sanitizeJSON(body []byte) string {
    if len(body) == 0 {
        return ""
    }

    var data any

    if err := json.Unmarshal(body, &data); err != nil {
        return "[invalid JSON body omitted]"
    }

    sanitized, err := json.Marshal(data)
    if err != nil {
        return "[JSON sanitization failed]"
    }

    return string(sanitized)
}