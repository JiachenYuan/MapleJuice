package global

const (
	LOGGER_PORT          = "55555"
	FD_PORT              = "55556"
	SDFS_PORT            = "55557"
	LEADER_ELECTION_PORT = "55558"

	QUORUM        = 6
	MAX_NUM_RETRY = 5
)

var (
	// This is supposed to be const, DO NOT MODIFY!!
	SERVER_ADDRS = []string{
		"fa23-cs425-1801.cs.illinois.edu", "fa23-cs425-1802.cs.illinois.edu",
		"fa23-cs425-1803.cs.illinois.edu", "fa23-cs425-1804.cs.illinois.edu",
		"fa23-cs425-1805.cs.illinois.edu", "fa23-cs425-1806.cs.illinois.edu",
		"fa23-cs425-1807.cs.illinois.edu", "fa23-cs425-1808.cs.illinois.edu",
		"fa23-cs425-1809.cs.illinois.edu", "fa23-cs425-1810.cs.illinois.edu"}
)

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// adapted from: https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
