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

// SmbStatusShare represents a single entry from the output of 'smbstatus -S'
type SmbStatusShare struct {
	Service     string              `json:"service"`
	ServerID    SmbStatusServerID   `json:"server_id"`
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
	RemoteMachine  string              `json:"remote_machine"`
	Hostname       string              `json:"hostname"`
	SessionDialect string              `json:"session_dialect"`
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
	TCons       map[string]SmbStatusShare      `json:"tcons"`
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

// RunSmbStatusShares executes 'smbstatus -S' on host container
func RunSmbStatusShares() ([]SmbStatusShare, error) {
	// Case 1: using new json output
	dat, err := executeSmbStatusCommand("-S --json")
	if err == nil {
		return parseSmbStatusSharesAsJSON(dat)
	}
	// Case 2: fallback to raw-text output
	dat, err = executeSmbStatusCommand("-S")
	if err == nil {
		return parseSmbStatusShares(dat)
	}
	return []SmbStatusShare{}, err
}

func parseSmbStatusSharesAsJSON(dat string) ([]SmbStatusShare, error) {
	shares := []SmbStatusShare{}
	res, err := parseSmbStatusJSON(dat)
	if err != nil {
		return shares, err
	}
	for _, share := range res.TCons {
		shares = append(shares, share)
	}
	return shares, nil
}

// RunSmbStatusLocks executes 'smbstatus -L' on host container
func RunSmbStatusLocks() ([]SmbStatusLock, error) {
	dat, err := executeSmbStatusCommand("-L")
	if err != nil {
		return []SmbStatusLock{}, err
	}
	return parseSmbStatusLocks(dat)
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

// RunSmbStatusProcs executes 'smbstatus -p' on host container
func RunSmbStatusProcs() ([]SmbStatusProc, error) {
	dat, err := executeSmbStatusCommand("-p")
	if err != nil {
		return []SmbStatusProc{}, err
	}
	return parseSmbStatusProcs(dat)
}

// SmbStatusSharesByMachine converts the output of RunSmbStatusShares into map
// indexed by machine's host
func SmbStatusSharesByMachine() (map[string][]SmbStatusShare, error) {
	shares, err := RunSmbStatusShares()
	if err != nil {
		return map[string][]SmbStatusShare{}, err
	}
	return makeSmbSharesMap(shares), nil
}

func makeSmbSharesMap(shares []SmbStatusShare) map[string][]SmbStatusShare {
	ret := map[string][]SmbStatusShare{}
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

// parseSmbStatusShares parses to output of 'smbstatus -S' into internal
// representation.
func parseSmbStatusShares(data string) ([]SmbStatusShare, error) {
	shares := []SmbStatusShare{}
	serviceIndex := 0
	pidIndex := 0
	machineIndex := 0
	connectedAtIndex := 0
	encryptionIndex := 0
	signingIndex := 0
	hasDashLine := false
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		ln := strings.TrimSpace(line)
		// Ignore empty and coment lines
		if len(ln) == 0 || ln[0] == '#' {
			continue
		}
		// Detect the all-dash line
		if strings.HasPrefix(ln, "------") {
			hasDashLine = true
			continue
		}
		// Parse header line into index of data
		if strings.HasPrefix(ln, "Service") {
			serviceIndex = strings.Index(ln, "Service")
			pidIndex = strings.Index(ln, "pid")
			machineIndex = strings.Index(ln, "Machine")
			connectedAtIndex = strings.Index(ln, "Connected at")
			encryptionIndex = strings.Index(ln, "Encryption")
			signingIndex = strings.Index(ln, "Signing")
			continue
		}
		// Ignore lines before header
		if !hasDashLine {
			continue
		}
		// Parse data into internal repr
		share := SmbStatusShare{}
		share.Service = parseSubstr(ln, serviceIndex)
		share.ServerID.PID = parseSubstr(ln, pidIndex)
		share.Machine = parseSubstr(ln, machineIndex)
		share.ConnectedAt = parseSubstr2(ln, connectedAtIndex, encryptionIndex)
		share.Encryption.Cipher = parseSubstr(ln, encryptionIndex)
		share.Signing.Cipher = parseSubstr(ln, signingIndex)

		// Ignore "IPC$"
		if share.Service == "IPC$" {
			continue
		}

		shares = append(shares, share)
	}
	return shares, nil
}

// parseSmbStatusProcs parses to output of 'smbstatus -p' into internal
// representation.
func parseSmbStatusProcs(data string) ([]SmbStatusProc, error) {
	procs := []SmbStatusProc{}
	pidIndex := 0
	usernameIndex := 0
	groupIndex := 0
	machineIndex := 0
	protocolVersionIndex := 0
	encryptionIndex := 0
	signingIndex := 0
	hasDashLine := false
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		ln := strings.TrimSpace(line)
		// Ignore empty and coment lines
		if len(ln) == 0 || ln[0] == '#' {
			continue
		}
		// Detect the all-dash line
		if strings.HasPrefix(ln, "------") {
			hasDashLine = true
			continue
		}
		// Parse header line into index of data
		if strings.HasPrefix(ln, "PID") {
			pidIndex = strings.Index(ln, "PID")
			usernameIndex = strings.Index(ln, "Username")
			groupIndex = strings.Index(ln, "Group")
			machineIndex = strings.Index(ln, "Machine")
			protocolVersionIndex = strings.Index(ln, "Protocol Version")
			encryptionIndex = strings.Index(ln, "Encryption")
			signingIndex = strings.Index(ln, "Signing")
			continue
		}
		// Ignore lines before header
		if !hasDashLine {
			continue
		}
		// Parse data into internal repr
		proc := SmbStatusProc{}
		proc.PID = parseSubstr(ln, pidIndex)
		proc.Username = parseSubstr(ln, usernameIndex)
		proc.Group = parseSubstr(ln, groupIndex)
		proc.Machine = parseSubstr(ln, machineIndex)
		proc.ProtocolVersion = parseSubstr(ln, protocolVersionIndex)
		proc.Encryption = parseSubstr(ln, encryptionIndex)
		proc.Signing = parseSubstr(ln, signingIndex)
		procs = append(procs, proc)
	}
	return procs, nil
}

// parseSmbStatusLocks parses to output of 'smbstatus -L' into internal
// representation.
func parseSmbStatusLocks(data string) ([]SmbStatusLock, error) {
	locks := []SmbStatusLock{}
	pidIndex := 0
	userIndex := 0
	denyModeIndex := 0
	accessIndex := 0
	rwIndex := 0
	oplockIndex := 0
	sharePathIndex := 0
	hasDashLine := false
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		ln := strings.TrimSpace(line)
		// Ignore empty and coment lines
		if len(ln) == 0 || ln[0] == '#' {
			continue
		}
		// Detect the all-dash line
		if strings.HasPrefix(ln, "------") {
			hasDashLine = true
			continue
		}
		// Ignore generic-info line
		if strings.HasPrefix(ln, "Locked files") {
			continue
		}
		// Parse header line into index of data
		if strings.HasPrefix(ln, "Pid") {
			pidIndex = strings.Index(ln, "Pid")
			userIndex = strings.Index(ln, "User")
			denyModeIndex = strings.Index(ln, "DenyMode")
			accessIndex = strings.Index(ln, "Access")
			rwIndex = strings.Index(ln, "R/W")
			oplockIndex = strings.Index(ln, "Oplock")
			sharePathIndex = strings.Index(ln, "SharePath")
			continue
		}
		// Ignore lines before header
		if !hasDashLine {
			continue
		}
		// Parse data into internal repr
		lock := SmbStatusLock{}
		lock.PID = parseSubstr(ln, pidIndex)
		lock.UserID = parseSubstr(ln, userIndex)
		lock.DenyMode = parseSubstr(ln, denyModeIndex)
		lock.Access = parseSubstr(ln, accessIndex)
		lock.RW = parseSubstr(ln, rwIndex)
		lock.Oplock = parseSubstr(ln, oplockIndex)
		lock.SharePath = parseSubstr(ln, sharePathIndex)
		locks = append(locks, lock)
	}
	return locks, nil
}

func parseSubstr(s string, startIndex int) string {
	sub := strings.TrimSpace(s[startIndex:])
	fields := strings.Fields(sub)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func parseSubstr2(s string, startIndex, endIndex int) string {
	return strings.TrimSpace(s[startIndex:endIndex])
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
