// Code generated by "stringer -type=FieldName -linecomment"; DO NOT EDIT.

package walletfee

import "strconv"

const _FieldName_name = "timetx_hashreserve_addrwallet_addrlog_indextrade_log_indexamountcountry"

var _FieldName_index = [...]uint8{0, 4, 11, 23, 34, 43, 58, 64, 71}

func (i FieldName) String() string {
	if i < 0 || i >= FieldName(len(_FieldName_index)-1) {
		return "FieldName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldName_name[_FieldName_index[i]:_FieldName_index[i+1]]
}