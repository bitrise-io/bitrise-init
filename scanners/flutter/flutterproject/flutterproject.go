package flutterproject

type FlutterAndDartSDKVersions struct {
	FlutterSDKVersions []VersionConstraint
	DartSDKVersions    []VersionConstraint
}

type SDKVersionsReader interface {
	ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error)
}

type Project struct {
	rootDir     string
	fileOpener  FileOpener
	sdkVersions FlutterAndDartSDKVersions
}

func New(rootDir string, fileOpener FileOpener) (*Project, error) {
	sdkVersions, err := readFlutterAndDartSDKVersions(rootDir, fileOpener)
	if err != nil {
		return nil, err
	}

	return &Project{
		rootDir:     rootDir,
		fileOpener:  fileOpener,
		sdkVersions: sdkVersions,
	}, nil
}

func (p Project) FlutterAndDartSDKVersions() FlutterAndDartSDKVersions {
	return p.sdkVersions
}

func (p Project) FlutterSDKVersionToUse() string {
	sdkVersions := p.FlutterAndDartSDKVersions()
	for _, flutterSDKVersion := range sdkVersions.FlutterSDKVersions {
		if flutterSDKVersion.Version != nil {
			return flutterSDKVersion.Version.String()
		} else if flutterSDKVersion.Constraint != nil {

		}
	}
	return ""
}

func readFlutterAndDartSDKVersions(rootDir string, fileOpener FileOpener) (FlutterAndDartSDKVersions, error) {
	versionReaders := []SDKVersionsReader{
		NewFVMVersionReader(fileOpener),
		NewASDFVersionReader(fileOpener),
		NewPubspecLockVersionReader(fileOpener),
		NewPubspecVersionReader(fileOpener),
	}

	var flutterSDKVersions []VersionConstraint
	var dartSDKVersions []VersionConstraint
	for _, versionReader := range versionReaders {
		flutterSDKVersion, dartSDKVersion, err := versionReader.ReadSDKVersions(rootDir)
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
