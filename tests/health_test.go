package tests

import (
	"net/http"
	"testing"

	"github.com/juanmaAV/go-echo-blueprint/internal/health"
	"github.com/juanmaAV/go-echo-blueprint/tests/helpers"
	echotest "github.com/juanmaAV/go-utils/testutil/echo"
	httptest "github.com/juanmaAV/go-utils/testutil/http"
)

func Test_HealthCheck(t *testing.T) {
	srv := helpers.NewTestServer()
	handler := health.NewHandler(health.NewService())

	cases := []echotest.Case{
		{
			Name: "returns OK",
			Request: echotest.Request{
				Method: http.MethodGet,
				Url:    "/health",
			},
			Response: echotest.ExpectedResponse{Status: http.StatusOK},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, rec := echotest.PrepareContext(srv.Echo, tc)
			if err := handler.Check(ctx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			httptest.AssertStatus(t, rec, tc.Response.Status)
			httptest.AssertJSONField(t, rec, "status", "OK")
		})
	}
}
