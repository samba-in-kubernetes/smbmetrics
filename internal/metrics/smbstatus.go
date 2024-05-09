// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// SmbStatusServerID represents a server_id output field
type SmbStatusServerID struct {
	PID      string `json:"pid"`
	TaskID   string `json:"task_id"`
	VNN      string `json:"vnn"`
	UniqueID string `json:"unique_id"`
}

// SmbStatusEncryption represents a encryption output field
type SmbStatusEncryption struct {
	Cipher string `json:"cipher"`
	Degree string `json:"degree"`
}

// SmbStatusSigning represents a signing output field
type SmbStatusSigning struct {
	Cipher string `json:"cipher"`
	Degree string `json:"degree"`
}

// SmbStatusTreeCon represents a 'tcon' output field
type SmbStatusTreeCon struct {
	Service     string              `json:"service"`
	ServerID    SmbStatusServerID   `json:"server_id"`
	TConID      string              `json:"tcon_id"`
	SessionID   string              `json:"session_id"`
	Machine     string              `json:"machine"`
	ConnectedAt string              `json:"connected_at"`
	Encryption  SmbStatusEncryption `json:"encryption"`
	Signing     SmbStatusSigning    `json:"signing"`
}

// SmbStatusSession represents a session output field
type SmbStatusSession struct {
	SessionID      string              `json:"session_id"`
	ServerID       SmbStatusServerID   `json:"server_id"`
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
	Encryption     SmbStatusEncryption `json:"encryption"`
	Signing        SmbStatusSigning    `json:"signing"`
}

// SmbStatusFileID represents a fileid output field
type SmbStatusFileID struct {
	DevID int64 `json:"devid"`
	Inode int64 `json:"inode"`
	Extid int32 `json:"extid"`
}

// SmbStatusLockedFile represents a locked-file output field
type SmbStatusLockedFile struct {
	ServicePath       string          `json:"service_path"`
	Filename          string          `json:"filename"`
	FileID            SmbStatusFileID `json:"fileid"`
	NumPendingDeletes int             `json:"num_pending_deletes"`
}

// SmbStatus represents output of 'smbstatus --json'
type SmbStatus struct {
	Timestamp   string                         `json:"timestamp"`
	Version     string                         `json:"version"`
	SmbConf     string                         `json:"smb_conf"`
	Sessions    map[string]SmbStatusSession    `json:"sessions"`
	TCons       map[string]SmbStatusTreeCon    `json:"tcons"`
	LockedFiles map[string]SmbStatusLockedFile `json:"locked_files"`
}

// SmbStatusLocks represents output of 'smbstatus -L --json'
type SmbStatusLocks struct {
	Timestamp string                         `json:"timestamp"`
	Version   string                         `json:"version"`
	SmbConf   string                         `json:"smb_conf"`
	OpenFiles map[string]SmbStatusLockedFile `json:"open_files"`
}

// LocateSmbStatus finds the local executable of 'smbstatus' on host.
func LocateSmbStatus() (string, error) {
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

// RunSmbStatusVersion executes 'smbstatus --version' on host container
func RunSmbStatusVersion() (string, error) {
	ver, err := executeSmbStatusCommand("--version")
	if err != nil {
		return "", err
	}
	return ver, nil
}

// RunSmbStatusShares executes 'smbstatus --shares --json' on host
func RunSmbStatusShares() ([]SmbStatusTreeCon, error) {
	dat, err := executeSmbStatusCommand("--shares --json")
	if err != nil {
		return []SmbStatusTreeCon{}, err
	}
	return parseSmbStatusTreeCons(dat)
}

func parseSmbStatusTreeCons(dat string) ([]SmbStatusTreeCon, error) {
	tcons := []SmbStatusTreeCon{}
	res, err := parseSmbStatus(dat)
	if err != nil {
		return tcons, err
	}
	for _, share := range res.TCons {
		tcons = append(tcons, share)
	}
	return tcons, nil
}

// RunSmbStatusLocks executes 'smbstatus --locks --json' on host
func RunSmbStatusLocks() ([]SmbStatusLockedFile, error) {
	dat, err := executeSmbStatusCommand("--locks --json")
	if err != nil {
		return []SmbStatusLockedFile{}, err
	}
	return parseSmbStatusLockedFiles(dat)
}

func parseSmbStatusLockedFiles(dat string) ([]SmbStatusLockedFile, error) {
	lockedFiles := []SmbStatusLockedFile{}
	res, err := parseSmbStatusLocks(dat)
	if err != nil {
		return lockedFiles, err
	}
	for _, lfile := range res.OpenFiles {
		lockedFiles = append(lockedFiles, lfile)
	}
	return lockedFiles, nil
}

// SmbStatusSharesByMachine converts the output of RunSmbStatusShares into map
// indexed by machine's host
func SmbStatusSharesByMachine() (map[string][]SmbStatusTreeCon, error) {
	tcons, err := RunSmbStatusShares()
	if err != nil {
		return map[string][]SmbStatusTreeCon{}, err
	}
	return makeSmbSharesMap(tcons), nil
}

func makeSmbSharesMap(tcons []SmbStatusTreeCon) map[string][]SmbStatusTreeCon {
	ret := map[string][]SmbStatusTreeCon{}
	for _, share := range tcons {
		ret[share.Machine] = append(ret[share.Machine], share)
	}
	return ret
}

func executeSmbStatusCommand(args ...string) (string, error) {
	loc, err := LocateSmbStatus()
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

// parseSmbStatus parses to output of 'smbstatus --json' into internal
// representation.
func parseSmbStatus(data string) (*SmbStatus, error) {
	res := SmbStatus{}
	err := json.Unmarshal([]byte(data), &res)
	return &res, err
}

// parseSmbStatusLocks parses to output of 'smbstatus --locks --json' into
// internal representation.
func parseSmbStatusLocks(data string) (*SmbStatusLocks, error) {
	res := SmbStatusLocks{}
	err := json.Unmarshal([]byte(data), &res)
	return &res, err
}
