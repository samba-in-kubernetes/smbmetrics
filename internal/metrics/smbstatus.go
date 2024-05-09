// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
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

// SmbStatusJSON represents output of 'smbstatus --json'
type SmbStatusJSON struct {
	Timestamp   string                         `json:"timestamp"`
	Version     string                         `json:"version"`
	SmbConf     string                         `json:"smb_conf"`
	Sessions    map[string]SmbStatusSession    `json:"sessions"`
	TCons       map[string]SmbStatusTreeCon    `json:"tcons"`
	LockedFiles map[string]SmbStatusLockedFile `json:"locked_files"`
}

// SmbStatusProc represents a single entry from the output of 'smbstatus -p'
type SmbStatusProc struct {
	PID             string
	Username        string
	Group           string
	Machine         string
	ProtocolVersion string
	Encryption      string
	Signing         string
}

// SmbStatusLock represents a single entry from the output of 'smbstatus -L'
type SmbStatusLock struct {
	PID       string
	UserID    string
	DenyMode  string
	Access    string
	RW        string
	Oplock    string
	SharePath string
	Name      string
	Time      string
}

// SmbStatusJSON represents output of 'smbstatus -L --json'
type SmbStatusLocksJSON struct {
	Timestamp string                         `json:"timestamp"`
	Version   string                         `json:"version"`
	SmbConf   string                         `json:"smb_conf"`
	OpenFiles map[string]SmbStatusLockedFile `json:"open_files"`
}

// LocateSmbStatus finds the local executable of 'smbstatus' on host container.
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

// RunSmbStatusShares executes 'smbstatus -S --json' on host container
func RunSmbStatusShares() ([]SmbStatusTreeCon, error) {
	dat, err := executeSmbStatusCommand("-S --json")
	if err != nil {
		return []SmbStatusTreeCon{}, err
	}
	return parseSmbStatusSharesAsJSON(dat)
}

func parseSmbStatusSharesAsJSON(dat string) ([]SmbStatusTreeCon, error) {
	tcons := []SmbStatusTreeCon{}
	res, err := parseSmbStatusJSON(dat)
	if err != nil {
		return tcons, err
	}
	for _, share := range res.TCons {
		tcons = append(tcons, share)
	}
	return tcons, nil
}

// RunSmbStatusLocks executes 'smbstatus -L --json' on host container
func RunSmbStatusLockedFiles() ([]SmbStatusLockedFile, error) {
	dat, err := executeSmbStatusCommand("-L --json")
	if err != nil {
		return []SmbStatusLockedFile{}, err
	}
	return parseSmbStatusLocksAsJSON(dat)
}

func parseSmbStatusLocksAsJSON(dat string) ([]SmbStatusLockedFile, error) {
	lockedFiles := []SmbStatusLockedFile{}
	res, err := parseSmbStatusLocksJSON(dat)
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
	shares, err := RunSmbStatusShares()
	if err != nil {
		return map[string][]SmbStatusTreeCon{}, err
	}
	return makeSmbSharesMap(shares), nil
}

func makeSmbSharesMap(shares []SmbStatusTreeCon) map[string][]SmbStatusTreeCon {
	ret := map[string][]SmbStatusTreeCon{}
	for _, share := range shares {
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

func ParseTime(s string) (time.Time, error) {
	layouts := []string{
		time.ANSIC,
		time.UnixDate,
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	// samba's lib/util/time.c uses non standad layout...
	return time.Time{}, errors.New("unknow time format " + s)
}

// parseSmbStatusJSON parses to output of 'smbstatus --json' into internal
// representation.
func parseSmbStatusJSON(data string) (*SmbStatusJSON, error) {
	res := SmbStatusJSON{}
	err := json.Unmarshal([]byte(data), &res)
	return &res, err
}

// parseSmbStatusJSON parses to output of 'smbstatus -L --json' into internal
// representation.
func parseSmbStatusLocksJSON(data string) (*SmbStatusLocksJSON, error) {
	res := SmbStatusLocksJSON{}
	err := json.Unmarshal([]byte(data), &res)
	return &res, err
}
