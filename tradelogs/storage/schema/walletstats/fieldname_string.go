// Code generated by "stringer -type=FieldName -linecomment"; DO NOT EDIT.

package walletstats

import "strconv"

const _FieldName_name = "timeunique_addresseseth_volumeusd_volumetotal_tradeusd_per_tradeeth_per_tradetotal_burn_feenew_unique_addresseskycedwallet_addr"

var _FieldName_index = [...]uint8{0, 4, 20, 30, 40, 51, 64, 77, 91, 111, 116, 127}

func (i FieldName) String() string {
	if i < 0 || i >= FieldName(len(_FieldName_index)-1) {
		return "FieldName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldName_name[_FieldName_index[i]:_FieldName_index[i+1]]
}