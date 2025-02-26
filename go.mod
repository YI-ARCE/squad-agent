module squad

go 1.21

replace (
	yiarce/core => ./core
	yiarce/core/date => ./core/date
	yiarce/core/file => ./core/file
	yiarce/core/frame => ./core/frame
	yiarce/core/log => ./core/log
	yiarce/core/timing => ./core/timing
	yiarce/core/yorm => ./core/yorm
	yiarce/core/yorm/mysql => ./core/yorm/mysql
)

require (
	github.com/fsnotify/fsnotify v1.8.0
	github.com/gorilla/websocket v1.5.3
	yiarce/core/date v0.0.0-00010101000000-000000000000
	yiarce/core/file v0.0.0-00010101000000-000000000000
	yiarce/core/frame v0.0.0-00010101000000-000000000000
	yiarce/core/log v0.0.0-00010101000000-000000000000
	yiarce/core/timing v0.0.0-00010101000000-000000000000
	yiarce/core/yorm v0.0.0-00010101000000-000000000000
	yiarce/core/yorm/mysql v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.29.0 // indirect
	yiarce/core v0.0.0-00010101000000-000000000000 // indirect
)
