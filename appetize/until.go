package appetize

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/ziputil"
)

// Platform [IOS, Android, Unknown]
type Platform string

// Extension  [.zip, .app, .gz, .apk]
type Extension string

// const
const (
	IOS     Platform = "ios"
	Android Platform = "android"
	Unknown Platform = "unknown"

	ZIP Extension = ".zip"
	APP Extension = ".app"
	GZ  Extension = ".gz"
	APK Extension = ".apk"
)

// Artifact ...
type Artifact struct {
	filePath  string
	platform  Platform
	extension Extension
}

// -------------------------------------
// -- Artifact methods

// NewArtifact checks the extension for the given file.
// For iOS platform the accepted extensions: .app, .zip, .tar.gz.
// For Android platform the accepted extension: .apk.
// Returns a new Artifact with the platform of the provided file [iOS, Android] and it's extension.
// **If the extension is not .app, .zip, .tar.gz. or .apk it will return an error.**
func NewArtifact(filePath string) (Artifact, error) {
	extension := path.Ext(filePath)

	switch extension {
	case ".zip":
		art := Artifact{filePath: filePath, platform: IOS, extension: ZIP}
		return art, nil
	case ".gz":
		art := Artifact{filePath: filePath, platform: IOS, extension: GZ}
		return art, nil
	case ".app":
		art := Artifact{filePath: filePath, platform: IOS, extension: APP}
		return art, nil
	case ".apk":
		art := Artifact{filePath: filePath, platform: Android, extension: APK}
		return art, nil
	}

	return Artifact{}, fmt.Errorf("bad file extension. For iOS, upload a .zip or .tar.gz file containing your compressed .app bundle. For Android, upload the .apk containing your app. Provided file's extension: %s. You can read more about it here: https://appetize.io/upload", extension)
}

// EnsureExtension checks the extension for the given platform.
// For iOS platform the accepted extensions: .app, .zip, .tar.gz.
// For Android platform the accepted extension: .apk.
// Returns the path of the file if the extension is valid.
// **If the platorm is iOS and the given file's extension is .app it creates a new .zip and returns the .zip's path.**
func (a Artifact) EnsureExtension() (string, error) {
	switch a.Platform() {
	case IOS:
		if a.Ext() == ZIP || a.Ext() == GZ {
			log.Printf("✅ Provided file is %s", a.Ext())
			return a.Path(), nil
		} else if a.Ext() == APP {
			log.Warnf("🚨 Provided file is %s", a.Ext())
			log.Printf("Need to compress it...")

			zipPth := strings.Replace(a.Path(), ".app", ".zip", 1)
			return zipPth, ensureZIP(a.Path(), path.Base(zipPth))
		}
		log.Warnf("platform.Ext(): %s", a.Ext())
	case Android:
		if a.Ext() == APK {
			log.Printf("✅ Provided file is %s", a.Ext())
			return a.Path(), nil
		}
	}

	return "", fmt.Errorf("bad file extension. For iOS, upload a .zip or .tar.gz file containing your compressed .app bundle. For Android, upload the .apk containing your app. Provided file: %s and provided platform: %s You can read more about it here: https://appetize.io/upload", path.Base(a.Path()), a.Platform())
}

// Path returns the artifact's path
func (a Artifact) Path() string {
	return a.filePath
}

// Ext returns the artifact's extension
func (a Artifact) Ext() Extension {
	return a.extension
}

// Platform returns the artifact's platform
func (a Artifact) Platform() Platform {
	return a.platform
}

func (a Artifact) String() string {
	return a.platform.String()
}

func (p Platform) String() string {
	return string(p)
}

func (e Extension) String() string {
	return string(e)
}

// -------------------------------------
// -- Private methods

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
