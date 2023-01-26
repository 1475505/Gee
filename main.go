package main

import (
	"gee.com/gee/gee"
)

func main() {
	r := gee.New()
	r.Static("/assets", "/home/ll/go/Gee")
	// 或相对路径 r.Static("/assets", "./static")
	// eg.http://localhost:9999/assets/go.mod
	r.Run(":9999")
}
