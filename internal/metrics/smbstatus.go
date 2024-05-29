// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// SMBStatusServerID represents a server_id output field
type SMBStatusServerID struct {
	PID      string `json:"pid"`
	TaskID   string `json:"task_id"`
	VNN      string `json:"vnn"`
	UniqueID string `json:"unique_id"`
}

// SMBStatusEncryption represents a encryption output field
type SMBStatusEncryption struct {
	Cipher string `json:"cipher"`
	Degree string `json:"degree"`
}

// SMBStatusSigning represents a signing output field
type SMBStatusSigning struct {
	Cipher string `json:"cipher"`
	Degree string `json:"degree"`
}

// SMBStatusTreeCon represents a 'tcon' output field
type SMBStatusTreeCon struct {
	Service     string              `json:"service"`
	ServerID    SMBStatusServerID   `json:"server_id"`
	TConID      string              `json:"tcon_id"`
	SessionID   string              `json:"session_id"`
	Machine     string              `json:"machine"`
	ConnectedAt string              `json:"connected_at"`
	Encryption  SMBStatusEncryption `json:"encryption"`
	Signing     SMBStatusSigning    `json:"signing"`
}

// SMBStatusSession represents a session output field
type SMBStatusSession struct {
	SessionID      string              `json:"session_id"`
	ServerID       SMBStatusServerID   `json:"server_id"`
	UID            int                 `json:"uid"`
	GID            int                 `json:"gid"`
	Username       string              `json:"username"`
	Groupname      string              `json:"groupname"`
	CreationTime   string              `json:"creation_time"`
	ExpirationTime string              `json:"expiration_time"`
	AuthTime       string              `json:"auth_time"`
	RemoteMachine  string              `json:"remote_machine"`
	Hostname       string              `json:"hostname"`
	SessionDialect string              `json:"session_dialect"`
	ClientGUID     string              `json:"client_guid"`
	Encryption     SMBStatusEncryption `json:"encryption"`
	Signing        SMBStatusSigning    `json:"signing"`
}

// SMBStatusFileID represents a fileid output field
type SMBStatusFileID struct {
	DevID int64 `json:"devid"`
	Inode int64 `json:"inode"`
	Extid int32 `json:"extid"`
}

// SMBStatusLockedFile represents a locked-file output field
type SMBStatusLockedFile struct {
	ServicePath       string          `json:"service_path"`
	Filename          string          `json:"filename"`
	FileID            SMBStatusFileID `json:"fileid"`
	NumPendingDeletes int             `json:"num_pending_deletes"`
}

// SMBStatus represents output of 'smbstatus --json'
type SMBStatus struct {
	Timestamp   string                         `json:"timestamp"`
	Version     string                         `json:"version"`
	SmbConf     string                         `json:"smb_conf"`
	Sessions    map[string]SMBStatusSession    `json:"sessions"`
	TCons       map[string]SMBStatusTreeCon    `json:"tcons"`
	LockedFiles map[string]SMBStatusLockedFile `json:"locked_files"`
}

// SMBStatusLocks represents output of 'smbstatus -L --json'
type SMBStatusLocks struct {
	Timestamp string                         `json:"timestamp"`
	Version   string                         `json:"version"`
	SmbConf   string                         `json:"smb_conf"`
	OpenFiles map[string]SMBStatusLockedFile `json:"open_files"`
}

// LocateSMBStatus finds the local executable of 'smbstatus' on host.
func LocateSMBStatus() (string, error) {
	knowns := []string{
		"/usr/bin/smbstatus",
	}
	for _, loc := range knowns {
		fi, err := os.Stat(loc)
		if err != nil {
			continue
		}
		mode := fi.Mode()
		if !mode.IsRegular() {
			continue
		}
		if (mode & 0111) > 0 {
			return loc, nil
		}
	}
	return "", errors.New("failed to locate smbstatus")
}

// RunSMBStatusVersion executes 'smbstatus --version' on host container
func RunSMBStatusVersion() (string, error) {
	ver, err := executeSMBStatusCommand("--version")
	if err != nil {
		return "", err
	}
	return ver, nil
}

// RunSMBStatusShares executes 'smbstatus --shares --json' on host
func RunSMBStatusShares() ([]SMBStatusTreeCon, error) {
	dat, err := executeSMBStatusCommand("--shares --json")
	if err != nil {
		return []SMBStatusTreeCon{}, err
	}
	return parseSMBStatusTreeCons(dat)
}

func parseSMBStatusTreeCons(dat string) ([]SMBStatusTreeCon, error) {
	tcons := []SMBStatusTreeCon{}
	res, err := parseSMBStatus(dat)
	if err != nil {
		return tcons, err
	}
	for _, share := range res.TCons {
		tcons = append(tcons, share)
	}
	return tcons, nil
}

// RunSMBStatusLocks executes 'smbstatus --locks --json' on host
func RunSMBStatusLocks() ([]SMBStatusLockedFile, error) {
	dat, err := executeSMBStatusCommand("--locks --json")
	if err != nil {
		return []SMBStatusLockedFile{}, err
	}
	return parseSMBStatusLockedFiles(dat)
}

func parseSMBStatusLockedFiles(dat string) ([]SMBStatusLockedFile, error) {
	lockedFiles := []SMBStatusLockedFile{}
	res, err := parseSMBStatusLocks(dat)
	if err != nil {
		return lockedFiles, err
	}
	for _, lfile := range res.OpenFiles {
		lockedFiles = append(lockedFiles, lfile)
	}
	return lockedFiles, nil
}

// SMBStatusSharesByMachine converts the output of RunSMBStatusShares into map
// indexed by machine's host
func SMBStatusSharesByMachine() (map[string][]SMBStatusTreeCon, error) {
	tcons, err := RunSMBStatusShares()
	if err != nil {
		return map[string][]SMBStatusTreeCon{}, err
	}
	return makeSmbSharesMap(tcons), nil
}

func makeSmbSharesMap(tcons []SMBStatusTreeCon) map[string][]SMBStatusTreeCon {
	ret := map[string][]SMBStatusTreeCon{}
	for _, share := range tcons {
		ret[share.Machine] = append(ret[share.Machine], share)
	}
	return ret
}

func executeSMBStatusCommand(args ...string) (string, error) {
	loc, err := LocateSMBStatus()
	if err != nil {
		return "", err
	}
	return executeCommand(loc, args...)
}

func executeCommand(command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}
	res := strings.TrimSpace(string(out))
	return res, nil
}

// parseSMBStatus parses to output of 'smbstatus --json' into internal
// representation.
func parseSMBStatus(data string) (*SMBStatus, error) {
	res := NewSMBStatus()
	err := json.Unmarshal([]byte(data), res)
	return res, err
}

// parseSMBStatusLocks parses to output of 'smbstatus --locks --json' into
// internal representation.
func parseSMBStatusLocks(data string) (*SMBStatusLocks, error) {
	res := SMBStatusLocks{}
	err := json.Unmarshal([]byte(data), &res)
	return &res, err
}

// NewSMBStatus returns non-populated SMBStatus object
func NewSMBStatus() *SMBStatus {
	smbStatus := SMBStatus{
		Timestamp:   "",
		Version:     "",
		SmbConf:     "",
		Sessions:    map[string]SMBStatusSession{},
		TCons:       map[string]SMBStatusTreeCon{},
		LockedFiles: map[string]SMBStatusLockedFile{},
	}
	return &smbStatus
}
