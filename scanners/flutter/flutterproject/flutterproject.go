package flutterproject

type SDKVersionsReader interface {
	ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error)
}

type Project struct {
	rootDir    string
	fileOpener FileOpener
}

func New(rootDir string, fileOpener FileOpener) Project {
	return Project{
		rootDir:    rootDir,
		fileOpener: fileOpener,
	}
}

type FlutterAndDartSDKVersions struct {
	FlutterSDKVersions []VersionConstraint
	DartSDKVersions    []VersionConstraint
}

func (p Project) FlutterAndDartSDKVersions() (FlutterAndDartSDKVersions, error) {
	versionReaders := []SDKVersionsReader{
		NewFVMVersionReader(p.fileOpener),
		NewASDFVersionReader(p.fileOpener),
		NewPubspecLockVersionReader(p.fileOpener),
		NewPubspecVersionReader(p.fileOpener),
	}

	var flutterSDKVersions []VersionConstraint
	var dartSDKVersions []VersionConstraint
	for _, versionReader := range versionReaders {
		flutterSDKVersion, dartSDKVersion, err := versionReader.ReadSDKVersions(p.rootDir)
		if err != nil {
			return FlutterAndDartSDKVersions{}, err
		}
		if flutterSDKVersion != nil {
			flutterSDKVersions = append(flutterSDKVersions, *flutterSDKVersion)
		}
		if dartSDKVersion != nil {
			dartSDKVersions = append(dartSDKVersions, *dartSDKVersion)
		}
	}
	return FlutterAndDartSDKVersions{
		FlutterSDKVersions: flutterSDKVersions,
		DartSDKVersions:    dartSDKVersions,
	}, nil
}
