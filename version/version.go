package version

var (
	Version = "0.7.5"

	// git hash should be filled by:
	// 	go build -ldflags="-X github.com/caivega/cayley/version.GitHash=xxxx"

	GitHash   = "dev snapshot"
	BuildDate string
)
