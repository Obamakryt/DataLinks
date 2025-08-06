package slogger

import "errors"

var ScanError = errors.New("scan error")
var RowsError = errors.New("rows error")
var ServerError = errors.New("server error")
var DBError = errors.New("service_unavailable")

var AuthError = errors.New("invalid_credentials")

const SomeError = "opps... something went wrong try later"

type CodeOperation string

const (
	E001  CodeOperation = "Registration"
	E002  CodeOperation = "Login"
	E101  CodeOperation = "Add New Link In users Table"
	E1011 CodeOperation = "Add or Find Link in urls Table"
	E102  CodeOperation = "Find All User Links"
	E103  CodeOperation = "Delete Link in User Collection"
	E1031 CodeOperation = "Find Link for Delete"
	E104  CodeOperation = "Change Link in User Collection on New Link"
	E1041 CodeOperation = "Find Current Link Which Want Swag"
	E1042 CodeOperation = "Find Or Insert Swap Link"
)
