package compile

import (
	"flag"
	"fmt"
	"os"
)

var (
	Built     string // 编译时间
	GitCommit string // Git commit hash
	GoVersion string // Golang版本
	OsArch    string // 操作系统
)

func init() {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "print build info")
	flag.Parse()
	if showVersion {
		fmt.Print(getBuildInfoStr())
		os.Exit(1)
	}
}

// getBuildInfoStr return build info string.
func getBuildInfoStr() string {
	return fmt.Sprintf(`Go Version : %s
Git Commit : %s
Built      : %s 
OS/Arch    : %s
`, GoVersion, GitCommit, Built, OsArch)
}
