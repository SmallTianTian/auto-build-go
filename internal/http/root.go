package http

import (
	"fmt"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/templates"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	ast_util "github.com/SmallTianTian/fresh-go/utils/ast"
	config_util "github.com/SmallTianTian/fresh-go/utils/config"
)

var fileAndTmpl = map[string]string{
	"ui/http/root.go": templates.ReadTemplateFile("http/gin/ui_http_root.go.tmpl"),
	"server/http.go":  templates.ReadTemplateFile("http/server_http.go.tmpl"),
}

func NewHTTP() {
	pro := config.DefaultConfig.Project.Name
	org := config.DefaultConfig.Project.Org
	dir := config.DefaultConfig.Project.Path
	module := filepath.Join(org, pro)
	httpPort := config.DefaultConfig.HTTP.Port

	if httpPort <= 0 {
		logger.Debug("Not set real http port, will skip create http server.")
		return
	}
	logger.Debugf("Project name: %s\nOrganization: %s\nPath: %s\nHTTP port: %d", pro, org, dir, httpPort)

	var kRv = map[string]interface{}{"module": module}
	utils.WriteByTemplate(dir, fileAndTmpl, kRv)
	logger.Debug("Writing to the http file is complete.")

	addConfig(dir, httpPort)
	logger.Debug("Add config is complete.")

	addCmdRun(dir, module)
	logger.Debug("Add cmd run is complete.")
}

func addConfig(path string, httpPort int) {
	pg := filepath.Join(path, "config/config.go")
	fga := utils.File2GoAST(pg)
	ast_util.AddField2AstFile(fga, "Http", "int", []string{"Config", "Port"})
	ast_util.AddField2AstFile(fga, "HttpPrefix", "string", []string{"Config", "Application"})
	utils.WriteAstFile(pg, "", fga)

	config_util.WriteConfig(path, "http", httpPort, []string{"port"})
	config_util.WriteConfig(path, "HttpPrefix", "", []string{"application"})
}

func addCmdRun(path, module string) {
	path = filepath.Join(path, "cmd/base.go")
	fga := utils.File2GoAST(path)
	ast_util.AppendFuncCall2AstFile(fga, "server.RunHttp", []string{}, []string{"start"})
	ast_util.AppendFuncCall2AstFile(fga, "server.StopHttp", []string{}, []string{"stop"})
	ast_util.SetImport2AstFile(fga, fmt.Sprintf("%s/server", module))
	utils.WriteAstFile(path, "", fga)
}
