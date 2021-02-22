package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/SmallTianTian/fresh-go/pkg/logger"
)

// FirstMod 仅在首次创建项目的时候需要执行。
// 如果有 go.mod 文件，将直接返回.
func FirstMod(dir, oAp string, vendor bool) bool {
	if IsExist(filepath.Join(dir, "go.mod")) {
		logger.Debug("Has go.mod, DON'T init again.")
		return true
	}

	logger.Debugf("Begin exec `go mod init %s` in %s.", oAp, dir)
	if err := Exec(dir, "go", "mod", "init", oAp); err != nil {
		logger.Errorf("Go mod init failed.\n Please exec `go mod init %s`.\n Error: %v", oAp, err)
		return false
	}

	if vendor {
		logger.Debug("User vendor")
		if err := os.Mkdir(filepath.Join(dir, "vendor"), os.ModePerm); err != nil {
			return false
		}
	}
	return GoModRebuild(dir)
}

// GoModRebuild 重新整理 Go mod
// 如果有 vendor 目录，会同时重新构建 vendor 目录.
func GoModRebuild(dir string) bool {
	logger.Debugf("Begin exec `go mod tidy` in %s.", dir)
	if err := Exec(dir, "go", "mod", "tidy"); err != nil {
		logger.Errorf("Go mod tidy failed.\n Please exec `go mod tidy` in `%s`.\n Error: %v", dir, err)
		return false
	}
	// 如果不包含 vendor 目录，将直接返回
	if !CheckUseVendor(dir) {
		return true
	}
	logger.Debugf("Begin exec `go mod vendor` in %s.", dir)
	if err := Exec(dir, "go", "mod", "vendor"); err != nil {
		logger.Errorf("Go mod tidy failed.\n Please exec `go mod tidy` in `%s`.\n Error: %v", dir, err)
		return false
	}
	return true
}

// GoFmtCode 格式化 golang 代码.
func GoFmtCode(path string) bool {
	if err := Exec(path, "gofmt", "-s", "-w", "."); err != nil {
		logger.Errorf("Go fmt code failed.\n Please exec `gofmt -s -w .` in `%s`.\n Error: %v", path, err)
		return false
	}
	return true
}

// GetOrganizationAndProjectName 从地址中获取组织名和项目名。
// 只查看 go.mod 文件，不包含任何 go 文件也可正常获取.
func GetOrganizationAndProjectName(path string) (org, pro string) {
	lines := ReadTxtFileEachLine(filepath.Join(path, "go.mod"))
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "module") {
			oAp := strings.TrimSpace(line)[len("module "):]
			org = filepath.Dir(oAp)
			pro = filepath.Base(oAp)
			return
		}
	}
	return
}

// CheckGoProject 检查是否是 Golang 项目.
// 必须使用 go.mod，否则也认为不是 go 项目，不予支持.
func CheckGoProject(path string) bool {
	return IsExist(filepath.Join(path, "go.mod"))
}

// CheckUseVendor 检查是否使用 vendor 目录.
// vendor 一定是目录，存在 vendor 文件也不会被认可.
func CheckUseVendor(path string) bool {
	info, err := os.Stat(filepath.Join(path, "vendor"))
	if err != nil && !os.IsExist(err) {
		return false
	}
	return info.IsDir()
}
