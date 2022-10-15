package fastcgi

// recType is a record type, as defined by
// http://www.fastcgi.com/devkit/doc/fcgi-spec.html#S8
type recType uint8

const (
	typeBeginRequest recType = iota + 1
	typeAbortRequest
	typeEndRequest
	typeParams
	typeStdin
	typeStdout
	typeStderr
	typeData
	typeGetValues
	typeGetValuesResult
	typeUnknownType
)

// Role for fastcgi application in spec
type Role uint16

const (
	RoleResponder Role = iota + 1
	RoleAuthorizer
	RoleFilter
)
