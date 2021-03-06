package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	ver "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

const examplePluginGitURL = "https://github.com/bitrise-core/bitrise-plugins-example.git"
const analyticsPluginBinURL = "https://github.com/bitrise-core/bitrise-plugins-analytics/releases/download/0.9.1/analytics-Darwin-x86_64"

func TestValidateVersion(t *testing.T) {
	t.Log("required min - pass")
	{
		requiredMin, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		current, err := ver.NewVersion("1.0.1")
		require.NoError(t, err)

		err = validateVersion(*current, *requiredMin, nil)
		require.NoError(t, err)
	}

	t.Log("required min - fail")
	{
		requiredMin, err := ver.NewVersion("1.0.2")
		require.NoError(t, err)

		current, err := ver.NewVersion("1.0.1")
		require.NoError(t, err)

		err = validateVersion(*current, *requiredMin, nil)
		require.Error(t, err)
	}

	t.Log("required min + required max - pass")
	{
		requiredMin, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		requiredMax, err := ver.NewVersion("1.0.2")
		require.NoError(t, err)

		current, err := ver.NewVersion("1.0.1")
		require.NoError(t, err)

		err = validateVersion(*current, *requiredMin, requiredMax)
		require.NoError(t, err)
	}

	t.Log("required min + required max - pass")
	{
		requiredMin, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		requiredMax, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		current, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		err = validateVersion(*current, *requiredMin, requiredMax)
		require.NoError(t, err)
	}

	t.Log("required min + required max - min fail")
	{
		requiredMin, err := ver.NewVersion("1.0.1")
		require.NoError(t, err)

		requiredMax, err := ver.NewVersion("1.0.2")
		require.NoError(t, err)

		current, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		err = validateVersion(*current, *requiredMin, requiredMax)
		require.Error(t, err)
	}

	t.Log("required min + required max - max fail")
	{
		requiredMin, err := ver.NewVersion("1.0.0")
		require.NoError(t, err)

		requiredMax, err := ver.NewVersion("1.0.1")
		require.NoError(t, err)

		current, err := ver.NewVersion("1.0.2")
		require.NoError(t, err)

		err = validateVersion(*current, *requiredMin, requiredMax)
		require.Error(t, err)
	}
}

func TestValidateRequirements(t *testing.T) {
	bitriseVersion, err := ver.NewVersion("1.0.0")
	require.NoError(t, err)

	envmanVersion, err := ver.NewVersion("1.0.0")
	require.NoError(t, err)

	stepmanVersion, err := ver.NewVersion("1.0.0")
	require.NoError(t, err)

	currentVersionMap := map[string]ver.Version{
		"bitrise": *bitriseVersion,
		"envman":  *envmanVersion,
		"stepman": *stepmanVersion,
	}

	t.Log("valid requirements")
	{
		requirements := []Requirement{
			Requirement{
				Tool:       "bitrise",
				MinVersion: "1.0.0",
				MaxVersion: "1.0.0",
			},
			Requirement{
				Tool:       "envman",
				MinVersion: "0.9.0",
				MaxVersion: "1.1.0",
			},
			Requirement{
				Tool:       "stepman",
				MinVersion: "1.0.0",
				MaxVersion: "1.0.0",
			},
		}

		err := validateRequirements(requirements, currentVersionMap)
		require.NoError(t, err)
	}

	t.Log("invalid requirements")
	{
		requirements := []Requirement{
			Requirement{
				Tool:       "bitrise",
				MinVersion: "1.0.0",
				MaxVersion: "1.0.0",
			},
			Requirement{
				Tool:       "envman",
				MinVersion: "1.1.0",
				MaxVersion: "1.1.0",
			},
			Requirement{
				Tool:       "stepman",
				MinVersion: "1.0.0",
				MaxVersion: "1.0.0",
			},
		}

		err := validateRequirements(requirements, currentVersionMap)
		require.Error(t, err)
	}
}

func TestClonePluginSrc(t *testing.T) {
	t.Log("example plugin - latest version")
	{
		pluginSource := examplePluginGitURL
		versionTag := ""
		destinationDir, err := pathutil.NormalizedOSTempDirPath("TestClonePluginSrc")
		require.NoError(t, err)

		exist, err := pathutil.IsPathExists(destinationDir)
		require.NoError(t, err)
		if exist {
			err := os.RemoveAll(destinationDir)
			require.NoError(t, err)
		}

		version, hash, err := clonePluginSrc(pluginSource, versionTag, destinationDir)
		require.NoError(t, err)
		require.NotNil(t, version)
		require.NotEmpty(t, hash)

		exist, err = pathutil.IsPathExists(destinationDir)
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}

	t.Log("example plugin - 0.9.0 version")
	{
		pluginSource := examplePluginGitURL
		versionTag := "0.9.0"
		destinationDir, err := pathutil.NormalizedOSTempDirPath("TestClonePluginSrc")
		require.NoError(t, err)

		exist, err := pathutil.IsPathExists(destinationDir)
		require.NoError(t, err)
		if exist {
			err := os.RemoveAll(destinationDir)
			require.NoError(t, err)
		}

		version, hash, err := clonePluginSrc(pluginSource, versionTag, destinationDir)
		require.NoError(t, err)
		require.NotNil(t, version)
		require.Equal(t, "0.9.0", version.String())
		require.NotEmpty(t, hash)

		exist, err = pathutil.IsPathExists(destinationDir)
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}
}

func TestDownloadPluginBin(t *testing.T) {
	t.Log("example plugin bin - ")
	{
		pluginBinURL := analyticsPluginBinURL
		destinationDir, err := pathutil.NormalizedOSTempDirPath("TestDownloadPluginBin")
		require.NoError(t, err)

		exist, err := pathutil.IsPathExists(destinationDir)
		require.NoError(t, err)
		if exist {
			err := os.RemoveAll(destinationDir)
			require.NoError(t, err)
		}

		require.NoError(t, os.MkdirAll(destinationDir, 0777))

		destinationPth := filepath.Join(destinationDir, "example")

		require.NoError(t, downloadPluginBin(pluginBinURL, destinationPth))

		exist, err = pathutil.IsPathExists(destinationPth)
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}
}
