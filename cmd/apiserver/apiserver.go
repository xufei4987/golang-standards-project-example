package main

import (
	"golang-standards-project-example/internal/apiserver"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	apiserver.NewApp("user-apiserver").Run()
}
