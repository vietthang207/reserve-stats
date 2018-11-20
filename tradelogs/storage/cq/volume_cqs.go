package cq

import (
	"bytes"
	"text/template"

	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
)

const (
	// the trades from WETH-ETH doesn't count. Hence the select clause skips every trade ETH-WETH or WETH-ETH
	// These trades are excluded by its src_addr and dst_addr, which is 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE for ETH
	// and 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 for WETH
	rsvVolTemplate = `SELECT SUM({{.AmountType}}) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume  ` +
		`INTO {{.MeasurementName}} FROM ` +
		`(SELECT {{.AmountType}}, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE` +
		`((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR ` +
		`(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) ` +
		`AND {{.RsvAddressType}}!='') GROUP BY {{.AddressType}},{{.RsvAddressType}}`
	rsvVolHourMsmName = `rsv_volume_hour`
	rsvVolDayMsmName  = `rsv_volume_day`
)

// CreateAssetVolumeCqs return a set of cqs required for asset volume aggregation
func CreateAssetVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	assetVolDstHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour "+
			"FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE "+
			"((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR "+
			"(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY dst_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstHourCqs)
	assetVolSrcHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour "+
			"FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE "+
			"((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR "+
			"(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY src_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcHourCqs)
	assetVolDstDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_day",
		dbName,
		"1h",
		"2d",
		"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_day FROM "+
			"(SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE "+
			"((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR "+
			"(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY dst_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstDayCqs)

	assetVolSrcDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_day FROM "+
			"(SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE "+
			"((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR "+
			"(src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY src_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcDayCqs)
	return result, nil
}

//CreateUserVolumeCqs continueous query for aggregate user volume
func CreateUserVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	userVolumeDayCqs, err := libcq.NewContinuousQuery(
		"user_volume_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume "+
			"INTO user_volume_day FROM (SELECT eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY user_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, userVolumeDayCqs)
	userVolumeHourCqs, err := libcq.NewContinuousQuery(
		"user_volume_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(eth_amount) as eth_volume, SUM(usd_amount) as usd_volume "+
			"INTO user_volume_hour FROM (SELECT eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY user_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, userVolumeHourCqs)
	return result, nil
}

// RsvFieldsType declare the set of names requires to completed a reserveVolume Cqs
type RsvFieldsType struct {
	// AmountType: it can be dst_amount or src_amount
	AmountType string
	// RsvAddressType: it can be dst_rsv_amount or src_rsv_amount
	RsvAddressType string
	// Addresstype: it can be dst_addr or src_addr
	AddressType string
}

func renderRsvCqFromTemplate(tmpl *template.Template, mName string, types RsvFieldsType) (string, error) {
	var query bytes.Buffer
	err := tmpl.Execute(&query, struct {
		RsvFieldsType
		MeasurementName string
	}{
		RsvFieldsType:   types,
		MeasurementName: mName,
	})
	if err != nil {
		return "", err
	}
	return query.String(), nil
}

// CreateReserveVolumeCqs return a set of cqs required for asset volume aggregation
func CreateReserveVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result     []*libcq.ContinuousQuery
		cqsGroupBY = map[string]RsvFieldsType{
			"rsv_volume_src_src": {
				AmountType:     "src_amount",
				RsvAddressType: "src_rsv_addr",
				AddressType:    "src_addr"},
			"rsv_volume_src_dst": {
				AmountType:     "src_amount",
				RsvAddressType: "dst_rsv_addr",
				AddressType:    "src_addr"},
			"rsv_volume_dst_src": {
				AmountType:     "dst_amount",
				RsvAddressType: "src_rsv_addr",
				AddressType:    "dst_addr"},
			"rsv_volume_dst_dst": {
				AmountType:     "dst_amount",
				RsvAddressType: "dst_rsv_addr",
				AddressType:    "dst_addr"},
		}
	)

	tpml, err := template.New("cq.CreateReserveVolumeCqs").Parse(rsvVolTemplate)
	if err != nil {
		return nil, err
	}

	for name, types := range cqsGroupBY {
		query, err := renderRsvCqFromTemplate(tpml, rsvVolHourMsmName, types)
		if err != nil {
			return nil, err
		}
		hourCQ, err := libcq.NewContinuousQuery(
			name+"_hour",
			dbName,
			hourResampleInterval,
			hourResampleFor,
			query,
			"1h",
			nil,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, hourCQ)

		query, err = renderRsvCqFromTemplate(tpml, rsvVolDayMsmName, types)
		if err != nil {
			return nil, err
		}
		dayCQ, err := libcq.NewContinuousQuery(
			name+"_day",
			dbName,
			dayResampleInterval,
			dayResampleFor,
			query,
			"1d",
			nil,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, dayCQ)
	}

	return result, nil
}
