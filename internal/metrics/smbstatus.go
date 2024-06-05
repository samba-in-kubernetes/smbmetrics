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

// SMBStatusOpenFile represents a open-file output field
type SMBStatusOpenFile struct {
	ServicePath       string                       `json:"service_path"`
	Filename          string                       `json:"filename"`
	FileID            SMBStatusFileID              `json:"fileid"`
	NumPendingDeletes int                          `json:"num_pending_deletes"`
	Opens             map[string]SMBStatusOpenInfo `json:"opens"`
}

// SMBStatusOpenInfo represents a single entry open_files/opens output field
type SMBStatusOpenInfo struct {
	UID         int                     `json:"uid"`
	ShareFileID string                  `json:"share_file_id"`
	OpenedAt    string                  `json:"opened_at"`
	ShareMode   SMBStatusOpenShareMode  `json:"sharemode"`
	AccessMask  SMBStatusOpenAccessMask `json:"access_mask"`
	OpLock      SMBStatusOpenOpLock     `json:"oplock"`
	Lease       SMBStatusOpenLease      `json:"lease"`
}

// SMBStatusOpenShareMode represents a open file share-mode entry
type SMBStatusOpenShareMode struct {
	Read   bool   `json:"READ"`
	Write  bool   `json:"WRITE"`
	Delete bool   `json:"DELETE"`
	Text   string `json:"text"`
	Hex    string `json:"hex"`
}

// SMBStatusOpenShareMode represents a open file access-mask entry
type SMBStatusOpenAccessMask struct {
	ReadData             bool   `json:"READ_DATA"`
	WriteData            bool   `json:"WRITE_DATA"`
	AppendData           bool   `json:"APPEND_DATA"`
	ReadEA               bool   `json:"READ_EA"`
	WriteEA              bool   `json:"WRITE_EA"`
	Execute              bool   `json:"EXECUTE"`
	ReadAttributes       bool   `json:"READ_ATTRIBUTES"`
	WriteAttributes      bool   `json:"WRITE_ATTRIBUTES"`
	DeleteChild          bool   `json:"DELETE_CHILD"`
	Delete               bool   `json:"DELETE"`
	ReadControl          bool   `json:"READ_CONTROL"`
	WriteDAC             bool   `json:"WRITE_DAC"`
	Synchronize          bool   `json:"SYNCHRONIZE"`
	AccessSystemSecurity bool   `json:"ACCESS_SYSTEM_SECURITY"`
	Text                 string `json:"text"`
	Hex                  string `json:"hex"`
}

// SMBStatusOpenOplock represents an open file operation-lock entry
type SMBStatusOpenOpLock struct {
	Exclusive bool   `json:"EXCLUSIVE"`
	Batch     bool   `json:"BATCH"`
	LevelII   bool   `json:"LEVEL_II"`
	Lease     bool   `json:"LEASE"`
	Text      string `json:"text"`
}

// SMBStatusOpenLease represents an open file lease entry
type SMBStatusOpenLease struct {
	LeaseKey string `json:"lease_key"`
	Read     bool   `json:"READ"`
	Write    bool   `json:"WRITE"`
	Handle   bool   `json:"HANDLE"`
	Text     string `json:"text"`
	Hex      string `json:"hex"`
}

// SMBStatus represents output of 'smbstatus --json'
type SMBStatus struct {
	Timestamp string                       `json:"timestamp"`
	Version   string                       `json:"version"`
	SmbConf   string                       `json:"smb_conf"`
	Sessions  map[string]SMBStatusSession  `json:"sessions"`
	TCons     map[string]SMBStatusTreeCon  `json:"tcons"`
	OpenFiles map[string]SMBStatusOpenFile `json:"open_files"`
}

// SMBStatusLocks represents output of 'smbstatus -L --json'
type SMBStatusLocks struct {
	Timestamp string                       `json:"timestamp"`
	Version   string                       `json:"version"`
	SmbConf   string                       `json:"smb_conf"`
	OpenFiles map[string]SMBStatusOpenFile `json:"open_files"`
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

// RunSMBStatus executes 'smbstatus --json' on host machine
func RunSMBStatus() (*SMBStatus, error) {
	dat, err := executeSMBStatusCommand("--json")
	if err != nil {
		return &SMBStatus{}, err
	}
	return parseSMBStatus(dat)
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
func RunSMBStatusLocks() ([]SMBStatusOpenFile, error) {
	dat, err := executeSMBStatusCommand("--locks --json")
	if err != nil {
		return []SMBStatusOpenFile{}, err
	}
	return parseSMBStatusLockedFiles(dat)
}

func parseSMBStatusLockedFiles(dat string) ([]SMBStatusOpenFile, error) {
	lockedFiles := []SMBStatusOpenFile{}
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
		Timestamp: "",
		Version:   "",
		SmbConf:   "",
		Sessions:  map[string]SMBStatusSession{},
		TCons:     map[string]SMBStatusTreeCon{},
		OpenFiles: map[string]SMBStatusOpenFile{},
	}
	return &smbStatus
}
