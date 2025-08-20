package constants

var EntryDuplicationDelay = 5

type AdminLevel string

const (
	SuperAdmin AdminLevel = "super_admin"
	OrgAdmin   AdminLevel = "org_admin"
)
