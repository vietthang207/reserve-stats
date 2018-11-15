package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type tokenHeatmapQuery struct {
	httputil.TimeRangeQuery
	Asset string `form:"asset" binding:"required"`
}

func (sv *Server) getTokenHeatMap(c *gin.Context) {
	var (
		query tokenHeatmapQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	from, to, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	asset, err := core.LookupToken(sv.coreSetting, query.Asset)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	heatmap, err := sv.storage.GetTokenHeatmap(asset, from, to)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, heatmap)
}
