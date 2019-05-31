package version

/*
	For each patch (a.k.a. class3 upgrade), define the corresponding height
	at which the software switches to the patched logic
*/
const (
	// 001: 0 - protocol version, 01 - patch number for protocol v0
	H001_UNDELEGATE_PATCH = 1159800
)
