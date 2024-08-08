module flightmonitor

go 1.21.11

require (
	github.com/claranet/go-zabbix-api v1.0.2
	github.com/go-sql-driver/mysql v1.8.1
	github.com/slack-go/slack v0.13.1
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/claranet/go-zabbix-api => github.com/elastic-infra/go-zabbix-api v1.3.1
