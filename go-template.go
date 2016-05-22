package main

import (
	"flag"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/app"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/routes"
	"github.com/kardianos/osext"
)

func main() {
	env := flag.String("env", "development", "environment to use for database")
	port := flag.String("port", "3000", "html port")
	developmentMode := flag.Bool("development-mode", false, "launch program in development mode")
	flag.Parse()

	rootDir, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

	settings := &app.ApplicationSettings{
		Environment:     *env,
		Port:            *port,
		RootDirectory:   rootDir,
		DevelopmentMode: *developmentMode,
	}

	app.Initialize(settings)

	serverChannel := routes.SetupApplication(settings)
	<-serverChannel // Wait for sort to finish; discard sent value.
}
