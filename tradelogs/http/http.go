package http

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/gin-gonic/gin"
)

const limitedTimeRange = 24 * time.Hour

// Server serve trade logs through http endpoint
type Server struct {
	storage     storage.Interface
	host        string
	sugar       *zap.SugaredLogger
	coreSetting core.Interface
}

type tradeLogsQuery struct {
	From uint64 `form:"from"`
	To   uint64 `form:"to"`
}

func (sv *Server) getTradeLogs(c *gin.Context) {
	var query tradeLogsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	fromTime := time.Unix(0, int64(query.From)*int64(time.Millisecond))
	toTime := time.Unix(0, int64(query.To)*int64(time.Millisecond))

	if toTime.After(fromTime.Add(limitedTimeRange)) {
		err := fmt.Errorf("time range is too broad, must be smaller or equal to %d milliseconds", limitedTimeRange/time.Millisecond)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if toTime.Equal(time.Unix(0, 0)) {
		toTime = time.Now()
		fromTime = toTime.Add(-time.Hour)
	}

	tradeLogs, err := sv.storage.LoadTradeLogs(fromTime, toTime)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "fromTime", fromTime, "toTime", toTime)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		tradeLogs,
	)
}

func (sv *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/trade-logs", sv.getTradeLogs)
	r.GET("/asset-volume", sv.getAssetVolume)
	return r
}

// Start running http server to serve trade logs data
func (sv *Server) Start() error {
	r := sv.setupRouter()
	return r.Run(sv.host)
}

// NewServer returns an instance of HttpApi to serve trade logs
func NewServer(storage storage.Interface, host string, sugar *zap.SugaredLogger, sett core.Interface) *Server {
	return &Server{storage: storage, host: host, sugar: sugar, coreSetting: sett}
}
