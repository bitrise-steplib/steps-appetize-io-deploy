package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-appetize-io-deploy/appetize"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-steputils/tools"
)

var debugMode bool

// Config ...
type Config struct {
	AppPath   string          `env:"app_path,required"`
	Token     stepconf.Secret `env:"appetize_token,required"`
	PublicKey string          `env:"public_key"`
	Verbose   bool            `env:"verbose,required"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	debugMode = cfg.Verbose
	log.SetEnableDebugLog(debugMode)

	stepconf.Print(cfg)
	fmt.Println()

	log.Infof("Checking provided file's platform")

	// Artifact Section
	artifact, err := appetize.NewArtifact(cfg.AppPath)
	if err != nil {
		failf("Upload failed!\nError: %s", err)
	}

	log.Printf("✅ Platform found: %s", artifact.Platform())
	fmt.Println()

	log.Infof("Checking provided file's extension")

	pth, err := artifact.EnsureExtension()
	if err != nil {
		failf("Upload failed!\nError: %s", err)
	}

	fmt.Println()
	log.Infof("Upload")

	// Network section
	client := appetize.NewClient(string(cfg.Token), pth, artifact, cfg.PublicKey)

	if cfg.PublicKey == "" {
		log.Warnf("🚨 No public key provided")
		log.Printf("Uploading new app to Appetize.io")
	} else {
		log.Printf("✅ Public key provided: %s", cfg.PublicKey)
		log.Printf("Updating the provided app at Appetize.io")
	}

	response, err := client.DirectFileUpload()
	if err != nil {
		failf("Upload failed %s", err)
	}

	fmt.Println()
	log.Printf("🎉 Upload succeeded")

	logDebugPretty(&response)

	appURL := "https://" + path.Join("appetize.io/app", response.PublicKey)
	log.Printf("You can check your app at: %s", appURL)
	fmt.Println()

	log.Infof("Generating output")

	// Output section
	if err := tools.ExportEnvironmentWithEnvman("APPETIZE_APP_URL", appURL); err != nil {
		failf("Failed to generate output - %s", "APPETIZE_APP_URL")
	}

	log.Donef("APPETIZE_APP_URL: %s", appURL)
	fmt.Println()
	log.Donef("Done")
}

// -------------------------------------
// -- Private methods

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}

func logDebugPretty(v interface{}) {
	if !debugMode {
		return
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("Response: %+v\n", string(b))
}
