package cordova

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
// WidgetModel ...
type WidgetModel struct {
	ID       string `xml:"id,attr"`
	Version  string `xml:"version,attr"`
	XMLNS    string `xml:"xmlns,attr"`
	XMLNSCDV string `xml:"xmlns:cdv,attr"`
}

*/
func TestParseConfigXMLContent(t *testing.T) {
	widget, err := parseConfigXMLContent(testConfigXMLContent)
	require.NoError(t, err)
	require.Equal(t, "com.bitrise.cordovasample", widget.ID)
	require.Equal(t, "0.9.0", widget.Version)
	require.Equal(t, "http://www.w3.org/ns/widgets", widget.XMLNS)
	require.Equal(t, "http://cordova.apache.org/ns/1.0", widget.XMLNSCDV)
}

func TestConfigName(t *testing.T) {
	{
		name := ConfigName(false, false)
		require.Equal(t, "cordova-config", name)
	}

	{
		name := ConfigName(true, false)
		require.Equal(t, "cordova-ios-config", name)
	}

	{
		name := ConfigName(false, true)
		require.Equal(t, "cordova-android-config", name)
	}

	{
		name := ConfigName(true, true)
		require.Equal(t, "cordova-ios-android-config", name)
	}
}

const testConfigXMLContent = `<?xml version='1.0' encoding='utf-8'?>
<widget id="com.bitrise.cordovasample" version="0.9.0" xmlns="http://www.w3.org/ns/widgets" xmlns:cdv="http://cordova.apache.org/ns/1.0">
    <name>CordovaOnBitrise</name>
    <description>A sample Apache Cordova application that builds on Bitrise.</description>
    <content src="index.html" />
    <access origin="*" />
    <plugin name="cordova-plugin-whitelist" spec="1" />
    <allow-intent href="http://*/*" />
    <allow-intent href="https://*/*" />
    <allow-intent href="tel:*" />
    <allow-intent href="sms:*" />
    <allow-intent href="mailto:*" />
    <allow-intent href="geo:*" />
    <engine name="ios" spec="~4.3.1" />
    <platform name="android">
        <allow-intent href="market:*" />
    </platform>
    <platform name="ios">
        <allow-intent href="itms:*" />
        <allow-intent href="itms-apps:*" />
    </platform>
    <engine name="android" spec="~6.1.2" />
</widget>`
