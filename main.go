package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/ziputil"
	"github.com/bitrise-steplib/steps-appetize-io-deploy/appetize"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-steputils/tools"
)

// Config ...
type Config struct {
	AppPath   string          `env:"app_path,required"`
	Platform  string          `env:"platform,opt[ios,android]"`
	Token     stepconf.Secret `env:"appetize_token,required"`
	PublicKey string          `env:"public_key"`
	Verbose   bool            `env:"verbose,required"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	log.SetEnableDebugLog(cfg.Verbose)

	stepconf.Print(cfg)
	fmt.Println()

	log.Infof("Checking provided file extension")
	var pth string
	var err error
	if pth, err = ensureExtension(cfg.AppPath, cfg.Platform); err != nil {
		failf("Upload failed!\nError: %s", err)
	}

	fmt.Println()
	log.Infof("Upload")
	// curl https://tok_vjbxr9m95cwe7r3fpjpqh3y94w@api.appetize.io/v1/apps -F "file=@/Users/birmachera/Desktop/ip/XcodeArchiveTest.ipa" -F "platform=ios"
	client := appetize.NewClient(string(cfg.Token), pth, cfg.Platform, cfg.PublicKey)

	if cfg.PublicKey == "" {
		log.Warnf("No public key provided  🚨")
		log.Printf("Uploading new app to Appetize.io")
	} else {
		log.Printf("Public key provided: %s  ✅", cfg.PublicKey)
		log.Printf("Updating the provided app at Appetize.io")
	}

	var response appetize.Response
	if response, err = client.DirectFileUpload(); err != nil {
		failf("Upload failed %s", err)
	}

	fmt.Println()
	log.Printf("Upload succeeded  🎉")

	logDebugPretty(&response)

	appURL := path.Join("https://appetize.io/app", response.PublicKey)
	log.Printf("You can check your app at: %s", appURL)
	fmt.Println()

	log.Infof("Generating output")

	if err := tools.ExportEnvironmentWithEnvman("APPETIZE_APP_URL", appURL); err != nil {
		failf("Failed to generate output - %s", "APPETIZE_APP_URL")
	}

	log.Donef("APPETIZE_APP_URL: %s", appURL)
	fmt.Println()
	log.Donef("Done")
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	log.Warnf("For more details you can enable the debug logs by turning on the verbose step input.")
	os.Exit(1)
}

func logDebugPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("Response: %+v\n", string(b))
}

// ensureExtension checks the extension for the given platform.
// For iOS platform the accepted extensions: .app, .zip, .tar.gz.
// For Android platform the accepted extension: .apk.
// If the platorm is iOS and the given file's extension is .app it creates a new .zip and returns the .zip's path.
func ensureExtension(pth, platform string) (string, error) {
	extension := path.Ext(pth)

	if platform == "ios" {
		if extension == ".zip" || extension == ".gz" {
			log.Printf("Provided file is %s  ✅", extension)
			return pth, nil
		} else if extension == ".app" {
			log.Warnf("Provided file is %s  🚨", extension)
			log.Printf("Need to compress it...")

			zipPth := strings.Replace(pth, ".app", ".zip", 1)
			return zipPth, ensureZIP(pth, path.Base(zipPth))
		}
	} else {
		if extension == ".apk" {
			log.Printf("Provided file is %s  ✅", extension)
			return pth, nil
		}
	}

	return "", fmt.Errorf("bad file extension. For iOS, upload a .zip or .tar.gz file containing your compressed .app bundle. For Android, upload the .apk containing your app. Provided file: %s and provided platform: %s You can read more about it here: https://appetize.io/upload", path.Base(pth), platform)
}

func ensureZIP(sourcePath string, destination string) error {
	fmt.Println()
	log.Infof("Creating %s from %s", destination, sourcePath)

	info, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return ziputil.ZipDir(sourcePath, destination, false)
	}

	return ziputil.ZipFile(sourcePath, destination)
}
