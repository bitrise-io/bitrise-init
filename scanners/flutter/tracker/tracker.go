package tracker

import (
	"encoding/json"

	"github.com/bitrise-io/bitrise-init/scanners/flutter/flutterproject"
	"github.com/bitrise-io/go-utils/v2/analytics"
	"github.com/bitrise-io/go-utils/v2/log"
)

type FlutterTracker struct {
	tracker analytics.Tracker
	logger  log.Logger
}

func NewStepTracker(logger log.Logger) FlutterTracker {
	p := analytics.Properties{}
	return FlutterTracker{
		tracker: analytics.NewDefaultTracker(logger, p),
		logger:  logger,
	}
}

func (t *FlutterTracker) LogSDKVersions(versions flutterproject.FlutterAndDartSDKVersions) {
	p := analytics.Properties{}
	for _, flutterSDK := range versions.FlutterSDKVersions {
		key := "flutter_sdk_" + string(flutterSDK.Source)
		value := ""

		if flutterSDK.Version != nil {
			value = flutterSDK.Version.String()
		} else if flutterSDK.Constraint != nil {
			value = flutterSDK.Constraint.String()
		}

		if value != "" {
			p[key] = value
		}
	}

	for _, dartSDK := range versions.DartSDKVersions {
		key := "dart_sdk_" + string(dartSDK.Source)
		value := ""

		if dartSDK.Version != nil {
			value = dartSDK.Version.String()
		} else if dartSDK.Constraint != nil {
			value = dartSDK.Constraint.String()
		}

		if value != "" {
			p[key] = value
		}
	}

	//t.tracker.Enqueue("flutter_scanner_sdk_versions", p)
	t.debugPrint(p)
}

func (t *FlutterTracker) debugPrint(p analytics.Properties) {
	b, err := json.MarshalIndent(p, "", "  ")
	if err == nil {
		t.logger.Printf(string(b))
	}
}
