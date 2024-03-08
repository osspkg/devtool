package appgoppy

const tmplMainGO = `package main

import (
	app_{{.app_name}} "{{.app_module}}/app"

	"go.osspkg.com/goppy"
	{{if .mod_metrics}}"go.osspkg.com/goppy/metrics"
{{end}}{{if .mod_geoip}}"go.osspkg.com/goppy/geoip"
{{end}}{{if or .mod_oauth .mod_auth_jwt}}"go.osspkg.com/goppy/auth"
{{end}}{{if .mod_db_mysql}}"go.osspkg.com/goppy/ormmysql"
{{end}}{{if .mod_db_sqlite}}"go.osspkg.com/goppy/ormsqlite"
{{end}}{{if .mod_db_postgre}}"go.osspkg.com/goppy/ormpgsql"
{{end}}{{if or .mod_web_server .mod_web_client}}"go.osspkg.com/goppy/web"
{{end}}{{if or .mod_websocket_server .mod_websocket_client}}"go.osspkg.com/goppy/ws"
{{end}}{{if or .mod_unixsocket_server .mod_unixsocket_client}}"go.osspkg.com/goppy/unixsocket"
{{end}}{{if or .mod_dns_server .mod_dns_client}}"go.osspkg.com/goppy/xdns"
{{end}}
)

var Version = "v0.0.0-dev"

func main() {
	app := goppy.New()
	app.AppName("{{.app_name}}")
	app.AppVersion(Version)
	app.Plugins(
		{{if .mod_metrics}}metrics.WithMetrics(),{{end}}
		{{if .mod_geoip}}geoip.WithMaxMindGeoIP(),{{end}}
		{{if .mod_oauth}}auth.WithOAuth(),{{end}}
		{{if .mod_auth_jwt}}auth.WithJWT(),{{end}}
		{{if .mod_db_mysql}}ormmysql.WithMySQL(),{{end}}
		{{if .mod_db_sqlite}}ormsqlite.WithSQLite(),{{end}}
		{{if .mod_db_postgre}}ormpgsql.WithPostgreSQL(),{{end}}
		{{if .mod_web_server}}web.WithHTTP(),{{end}}
		{{if .mod_web_client}}web.WithHTTPClient(),{{end}}
		{{if .mod_websocket_server}}ws.WithWebsocketServer(),{{end}}
		{{if .mod_websocket_client}}ws.WithWebsocketClient(),{{end}}
		{{if .mod_dns_server}}xdns.WithDNSServer(),{{end}}
		{{if .mod_dns_client}}xdns.WithDNSClient(),{{end}}
		{{if .mod_unixsocket_server}}unixsocket.WithServer(),{{end}}
		{{if .mod_unixsocket_client}}unixsocket.WithClient(),{{end}}
	)
	app.Plugins(app_{{.app_name}}.Plugins...)
	app.Run()
}
`

const tmplAppGo = `package app

import (
	"go.osspkg.com/goppy/plugins"
)

var Plugins = plugins.Inject()

`

const tmplReadMe = `# {{.app_name}}

`

const tmplGitIgnore = `
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
*.dev.yaml

vendor/
build/
.idea/
.vscode/
.tools/

`

const tmplDockerFile = `version: '2.4'

networks:
  database:
    name: {{.app_name}}-dev-net

services:

  db:
    image: library/mysql:5.7.25
    restart: on-failure
    environment:
      MYSQL_ROOT_PASSWORD: 'root'
      MYSQL_USER: 'test'
      MYSQL_PASSWORD: 'test'
      MYSQL_DATABASE: 'test_database'
    healthcheck:
      test: [ "CMD", "mysql", "--user=root", "--password=root", "-e", "SHOW DATABASES;" ]
      interval: 15s
      timeout: 30s
      retries: 30
    ports:
      - "127.0.0.1:3306:3306"
    networks:
      - database

  adminer:
    image: adminer:latest
    restart: on-failure
    links:
      - db
    ports:
      - "127.0.0.1:8000:8080"
    networks:
      - database
`

const tmplGoMod = `module {{.app_module}}

go {{.go_version}}
`
