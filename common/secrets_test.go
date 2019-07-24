package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"

	"github.com/bitrise-io/go-utils/pathutil"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/stretchr/testify/require"
)

func TestLoadSecret(t *testing.T) {
	t.Log("Read from env var")
	{
		revokeFn, err := envutil.RevokableSetenv("SECRET_KEY_1", "secret value 1")
		defer func() {
			require.NoError(t, revokeFn())
		}()
		require.NoError(t, err)

		secretValue, err := LoadSecret("SECRET_KEY_1")
		require.NoError(t, err)
		require.Equal(t, "secret value 1", secretValue)
	}

	t.Log("Read secret from file")
	secretsDirPath, err := pathutil.NormalizedOSTempDirPath("LoadSecret")
	require.NoError(t, err)

	revokeFn, err := envutil.RevokableSetenv("SECRETS_CONFIG_DIR_PATH", secretsDirPath)
	defer func() {
		require.NoError(t, revokeFn())
	}()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(secretsDirPath))
	}()
	subdirPath := filepath.Join(secretsDirPath, "subdirtest")
	require.NoError(t, pathutil.EnsureDirExist(subdirPath))

	t.Log("Read from file directly in SECRETS_CONFIG_DIR_PATH")
	{
		filePth := filepath.Join(secretsDirPath, "SECRET_KEY_2")
		require.NoError(t, fileutil.WriteStringToFile(filePth, "secret value 2"))

		secretValue, err := LoadSecret("SECRET_KEY_2")
		require.NoError(t, err)
		require.Equal(t, "secret value 2", secretValue)
	}

	t.Log("Read from file in subdir of SECRETS_CONFIG_DIR_PATH")
	{
		filePth := filepath.Join(subdirPath, "SECRET_KEY_3")
		require.NoError(t, fileutil.WriteStringToFile(filePth, "secret value 3"))

		secretValue, err := LoadSecret("SECRET_KEY_3")
		require.NoError(t, err)
		require.Equal(t, "secret value 3", secretValue)
	}

	t.Log("Not found")
	{
		secretValue, err := LoadSecret("SECRET_KEY_NOT_FOUND")
		require.Error(t, err, "Secret (SECRET_KEY_NOT_FOUND) not found")
		require.Equal(t, "", secretValue)
	}
	t.Log("File exists but with empty content - should be Not found")
	{
		emptyfilePth := filepath.Join(secretsDirPath, "SECRET_KEY_EMPTYFILE")
		require.NoError(t, fileutil.WriteStringToFile(emptyfilePth, ""))

		secretValue, err := LoadSecret("SECRET_KEY_EMPTYFILE")
		require.Error(t, err, "Secret (SECRET_KEY_EMPTYFILE) not found")
		require.Equal(t, "", secretValue)
	}
	t.Log("Multiple files in subdirs with the same filename - should not read any")
	{
		subdirPath2 := filepath.Join(secretsDirPath, "subdirtest2")
		require.NoError(t, pathutil.EnsureDirExist(subdirPath2))

		filePth := filepath.Join(subdirPath2, "SECRET_KEY_3")
		require.NoError(t, fileutil.WriteStringToFile(filePth, "secret value 3"))

		secretValue, err := LoadSecret("SECRET_KEY_3")
		require.Error(t, err, "Secret (SECRET_KEY_3) not found")
		require.Equal(t, "", secretValue)
	}
}

func TestLoadSecrets(t *testing.T) {
	t.Log("Read from env var")
	{
		{
			revokeFn, err := envutil.RevokableSetenv("SECRET_KEY_1", "secret value 1")
			defer func() {
				require.NoError(t, revokeFn())
			}()
			require.NoError(t, err)
		}
		{
			revokeFn2, err := envutil.RevokableSetenv("SECRET_KEY_2", "secret value 2")
			defer func() {
				require.NoError(t, revokeFn2())
			}()
			require.NoError(t, err)
		}

		secretValues, err := LoadSecrets([]string{"SECRET_KEY_1", "SECRET_KEY_2"})
		require.NoError(t, err)
		require.Equal(t, map[string]string{"SECRET_KEY_1": "secret value 1", "SECRET_KEY_2": "secret value 2"}, secretValues)
	}
}
