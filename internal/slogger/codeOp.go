package slogger

import "errors"

var ScanError = errors.New("scan error")
var RowsError = errors.New("rows error")
var ServerError = errors.New("server error")

const SomeError = "opps... something went wrong"

type CodeOperation string

const (
	E101  CodeOperation = "Add New Link In users Table"
	E1011 CodeOperation = "Add or Find Link in urls Table"
	E102  CodeOperation = "Find All User Links"
	E103  CodeOperation = "Delete Link in User Collection"
	E1031 CodeOperation = "Find Link for Delete"
	E104  CodeOperation = "Swap Processing"
	E1041 CodeOperation = "Change Link in User Collection on New Link"
	E1042 CodeOperation = "Find Current Link Which Want Swag"
	E1043 CodeOperation = "Find Or Insert Swap Link"
)
