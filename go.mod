module 7-days-golang

go 1.19

require (
	geeRpc v0.0.0
	github.com/mattn/go-sqlite3 v1.14.15
	geeorm v0.0.0
)

replace geeRpc => ./geeRpc
replace (
	geeorm => ./geeORM
)
