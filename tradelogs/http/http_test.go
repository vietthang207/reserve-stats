package http

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type mockStorage struct {
}

func (s *mockStorage) SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error {
	return nil
}

func (s *mockStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	return nil, nil
}

func (s *mockStorage) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	return nil, nil
}

func newTestServer() (*Server, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return &Server{
		storage:     &mockStorage{},
		sugar:       sugar,
		coreSetting: &mockCore{}}, nil
}

func TestTradeLogsRoute(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test valid request",
			Endpoint: "/trade-logs",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result []common.TradeLog
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
			},
		},
		{
			Msg:      "Test invalid time range",
			Endpoint: fmt.Sprintf("/trade-logs?from=0&to=%d", time.Hour/time.Millisecond*25),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				assert.Contains(t, result.Error, "time range is too broad")
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}

func TestBurnFeeRoute(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test valid request",
			Endpoint: "/burn-fee?freq=h&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result map[ethereum.Address]map[string]float64
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
			},
		},
		{
			Msg:      "Test missing reserve address",
			Endpoint: "/burn-fee?freq=h",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result map[ethereum.Address]map[string]float64
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
			},
		},
		{
			Msg:      "Test invalid reserve address",
			Endpoint: "/burn-fee?freq=h&reserve=invalidAddress",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				assert.Contains(t, result.Error, "Field validation for 'ReserveAddrs[0]' failed on the 'isAddress' tag")
			},
		},
		{
			Msg:      "Test invalid frequency",
			Endpoint: "/burn-fee?freq=invalid&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				assert.Contains(t, result.Error, "your query frequency is not supported")
			},
		},
		{
			Msg:      "Test time range too broad",
			Endpoint: fmt.Sprintf("/burn-fee?from=0&to=%d&freq=h&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F", hourlyBurnFeeMaxDuration/time.Millisecond+1),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				expectedErrMsg := fmt.Sprintf("your query time range exceeds the duration limit %s", hourlyBurnFeeMaxDuration)
				assert.Equal(t, expectedErrMsg, result.Error)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}
