package utils

type Locks struct {
	ForUpdate           bool `json:"for_update"`
	SessionAdvisoryLock bool `json:"session_advisory_lock"`
}
