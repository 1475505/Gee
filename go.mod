module gee.com/gee

go 1.18

require gee v0.0.0
replace (
	gee => ./gee
	geerpc => ./gee-rpc
	codec => ./gee-rpc/codec
)