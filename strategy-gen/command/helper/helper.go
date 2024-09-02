package helper

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/tsinghua-cel/strategy-gen/command"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

// FormatKV formats key value pairs:
//
// Key = Value
//
// Key = <none>
func FormatKV(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	columnConf.Glue = " = "

	return columnize.Format(in, columnConf)
}

// Creates a file at path and with perms level permissions.
// If file already exists, owner and permissions are
// verified, and the file is overwritten.
func SaveFileSafe(path string, data []byte, perms fs.FileMode) error {
	info, err := os.Stat(path)
	// check if an error occurred other than path not exists
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if FileExists(path) {
		// verify that existing file's owner and permissions are safe
		if err := verifyFileOwnerAndPermissions(path, info, perms); err != nil {
			return err
		}
	}

	// create or overwrite the file
	return os.WriteFile(path, data, perms)
}

// Checks if the file at the specified path exists
func FileExists(filePath string) bool {
	// Check if path is empty
	if filePath == "" {
		return false
	}

	// Grab the absolute filepath
	pathAbs, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}

	// Check if the file exists, and that it's actually a file if there is a hit
	if fileInfo, statErr := os.Stat(pathAbs); os.IsNotExist(statErr) || (fileInfo != nil && fileInfo.IsDir()) {
		return false
	}

	return true
}

// Verifies that the file owner is the current user,
// or the file owner is in the same group as current user
// and permissions are set correctly by the owner.
func verifyFileOwnerAndPermissions(path string, info fs.FileInfo, expectedPerms fs.FileMode) error {
	// get stats
	stat, ok := info.Sys().(*syscall.Stat_t)
	if stat == nil || !ok {
		return fmt.Errorf("failed to get stats of %s", path)
	}

	// get current user
	currUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user")
	}

	// get user id of the owner
	ownerUID := strconv.FormatUint(uint64(stat.Uid), 10)
	if currUser.Uid == ownerUID {
		return nil
	}

	// get group id of the owner
	ownerGID := strconv.FormatUint(uint64(stat.Gid), 10)
	if currUser.Gid != ownerGID {
		return fmt.Errorf("file/directory created by a user from a different group: %s", path)
	}

	// check if permissions are set correctly by the owner
	if info.Mode() != expectedPerms {
		return fmt.Errorf("permissions of the file/directory '%s' are set incorrectly by another user", path)
	}

	return nil
}

// RegisterJSONOutputFlag registers the --json output setting for all child commands
func RegisterJSONOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(
		command.JSONOutputFlag,
		false,
		"get all outputs in json format (default false)",
	)
}
