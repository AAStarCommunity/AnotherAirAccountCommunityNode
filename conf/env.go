package conf

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	Name     string // 	环境名称
	Debugger bool   // 是否调试模式
}

func (env *Env) IsDevelopment() bool {
	return strings.EqualFold("dev", env.Name)
}

func (env *Env) IsProduction() bool {
	return strings.EqualFold("prod", env.Name)
}

func (env *Env) GetEnvName() *string {
	return &env.Name
}

func getConfFilePath() *string {
	if os.Getenv("app_conf_path") != "" {
		path := os.Getenv("app_conf_path")
		return &path
	}
	path := fmt.Sprintf("conf/appsettings.%s.yaml", strings.ToLower(Environment.Name))
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		path = "conf/appsettings.yaml"
	}
	return &path
}

var Environment *Env

func init() {
	envName := "prod"
	if len(os.Getenv("Env")) > 0 {
		envName = os.Getenv("Env")
	}
	Environment = &Env{
		Name:     envName,
		Debugger: strings.EqualFold("dev", envName),
	}
}
