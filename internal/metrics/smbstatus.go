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
	Timestamp string                      `json:"timestamp"`
	Version   string                      `json:"version"`
	SmbConf   string                      `json:"smb_conf"`
	Sessions  map[string]SMBStatusSession `json:"sessions"`
	TCons     map[string]SMBStatusTreeCon `json:"tcons"`
}

// SMBStatusLocks represents output of 'smbstatus -L --json'
type SMBStatusLocks struct {
	Timestamp string                       `json:"timestamp"`
	Version   string                       `json:"version"`
	SmbConf   string                       `json:"smb_conf"`
	OpenFiles map[string]SMBStatusOpenFile `json:"open_files"`
}

// SMBProfileEntry represents basic profile entry of 'smbstatus --profile'
type SMBProfileEntry struct {
	Count int `json:"count"`
	Time  int `json:"time"`
}

// SMBProfileSyscalls represents 'SMBD loop' entries of 'smbstatus --profile'
type SMBProfileLoop struct {
	Connect       SMBProfileEntry `json:"connect"`
	Disconnect    SMBProfileEntry `json:"disconnect"`
	Idle          SMBProfileEntry `json:"idle"`
	CPUUser       SMBProfileEntry `json:"cpu_user"`
	CPUSystem     SMBProfileEntry `json:"cpu_system"`
	Request       SMBProfileEntry `json:"request"`
	PushSecCtx    SMBProfileEntry `json:"push_sec_ctx"`
	SetSecCtx     SMBProfileEntry `json:"set_sec_ctx"`
	SetRootSecCtx SMBProfileEntry `json:"set_root_sec_ctx"`
	PopSecCtx     SMBProfileEntry `json:"pop_sec_ctx"`
}

// SMBProfileIOEntry represents async-io profile entry of 'smbstatus --profile'
type SMBProfileIOEntry struct {
	SMBProfileEntry
	Idle  int `json:"idle"`
	Bytes int `json:"bytes"`
}

// SMBProfileSyscalls represents 'System Calls' entries of 'smbstatus --profile'
type SMBProfileSyscalls struct {
	Opendir        SMBProfileEntry   `json:"syscall_opendir"`
	FDOpendir      SMBProfileEntry   `json:"syscall_fdopendir"`
	Readdir        SMBProfileEntry   `json:"syscall_readdir"`
	Rewinddir      SMBProfileEntry   `json:"syscall_rewinddir"`
	Mkdirat        SMBProfileEntry   `json:"syscall_mkdirat"`
	Closedir       SMBProfileEntry   `json:"syscall_closedir"`
	Open           SMBProfileEntry   `json:"syscall_open"`
	OpenAt         SMBProfileEntry   `json:"syscall_openat"`
	CreateFile     SMBProfileEntry   `json:"syscall_createfile"`
	Close          SMBProfileEntry   `json:"syscall_close"`
	PRead          SMBProfileIOEntry `json:"syscall_pread"`
	AsysPRead      SMBProfileIOEntry `json:"syscall_asys_pread"`
	PWrite         SMBProfileIOEntry `json:"syscall_pwrite"`
	AsysPWrite     SMBProfileIOEntry `json:"syscall_asys_pwrite"`
	Lseek          SMBProfileEntry   `json:"syscall_lseek"`
	SendFile       SMBProfileIOEntry `json:"syscall_sendfile"`
	RecvFile       SMBProfileIOEntry `json:"syscall_recvfile"`
	RenameAt       SMBProfileEntry   `json:"syscall_renameat"`
	AsysFSync      SMBProfileIOEntry `json:"syscall_asys_fsync"`
	Stat           SMBProfileEntry   `json:"syscall_stat"`
	FStat          SMBProfileEntry   `json:"syscall_fstat"`
	LStat          SMBProfileEntry   `json:"syscall_lstat"`
	FStatAt        SMBProfileEntry   `json:"syscall_fstatat"`
	GetAllocSize   SMBProfileEntry   `json:"syscall_get_alloc_size"`
	UnlinkAt       SMBProfileEntry   `json:"syscall_unlinkat"`
	Chmod          SMBProfileEntry   `json:"syscall_chmod"`
	FChmod         SMBProfileEntry   `json:"syscall_fchmod"`
	FChown         SMBProfileEntry   `json:"syscall_fchown"`
	LChown         SMBProfileEntry   `json:"syscall_lchown"`
	Chdir          SMBProfileEntry   `json:"syscall_chdir"`
	GetWD          SMBProfileEntry   `json:"syscall_getwd"`
	Fntimes        SMBProfileEntry   `json:"syscall_fntimes"`
	FTruncate      SMBProfileEntry   `json:"syscall_ftruncate"`
	FAllocate      SMBProfileEntry   `json:"syscall_fallocate"`
	ReadLinkAt     SMBProfileEntry   `json:"syscall_readlinkat"`
	SymLinkAt      SMBProfileEntry   `json:"syscall_symlinkat"`
	LinkAt         SMBProfileEntry   `json:"syscall_linkat"`
	MknodAt        SMBProfileEntry   `json:"syscall_mknodat"`
	RealPath       SMBProfileEntry   `json:"syscall_realpath"`
	GetQuota       SMBProfileEntry   `json:"syscall_get_quota"`
	SetQuota       SMBProfileEntry   `json:"syscall_set_quota"`
	AsysGetXattrAt SMBProfileIOEntry `json:"syscall_asys_getxattrat"`
}

// SMBStatusProfile represents single call entry of 'smbstatus --profile'
type SMBProfileCallEntry struct {
	SMBProfileEntry
	Idle     int `json:"idle"`
	Inbytes  int `json:"inbytes"`
	Outbytes int `json:"outbytes"`
}

// SMBProfileSMB2Calls represents 'SMB2 Calls' entries of 'smbstatus --profile'
type SMBProfileSMB2Calls struct {
	NegProt   SMBProfileCallEntry `json:"smb2_negprot"`
	SessSetup SMBProfileCallEntry `json:"smb2_sesssetup"`
	LogOff    SMBProfileCallEntry `json:"smb2_logoff"`
	Tcon      SMBProfileCallEntry `json:"smb2_tcon"`
	Tdis      SMBProfileCallEntry `json:"smb2_tdis"`
	Create    SMBProfileCallEntry `json:"smb2_create"`
	Close     SMBProfileCallEntry `json:"smb2_close"`
	Flush     SMBProfileCallEntry `json:"smb2_flush"`
	Read      SMBProfileCallEntry `json:"smb2_read"`
	Write     SMBProfileCallEntry `json:"smb2_write"`
	Lock      SMBProfileCallEntry `json:"smb2_lock"`
	Ioctl     SMBProfileCallEntry `json:"smb2_ioctl"`
	Cancel    SMBProfileCallEntry `json:"smb2_cancel"`
	KeepAlive SMBProfileCallEntry `json:"smb2_keepalive"`
	Find      SMBProfileCallEntry `json:"smb2_find"`
	Notify    SMBProfileCallEntry `json:"smb2_notify"`
	GetInfo   SMBProfileCallEntry `json:"smb2_getinfo"`
	SetInfo   SMBProfileCallEntry `json:"smb2_setinfo"`
	Break     SMBProfileCallEntry `json:"smb2_break"`
}

// SMBProfileShare represents per-share profile information
type SMBProfileShare struct {
	SystemCalls *SMBProfileSyscalls  `json:"System Calls"`
	SMB2Calls   *SMBProfileSMB2Calls `json:"SMB2 Calls"`
}

// SMBProfile represents (a subset of the) output of 'smbstatus --profile'
type SMBProfile struct {
	Timestamp   string                      `json:"timestamp"`
	Version     string                      `json:"version"`
	SmbConf     string                      `json:"smb_conf"`
	SmbdLoop    *SMBProfileLoop             `json:"SMBD loop"`
	SystemCalls *SMBProfileSyscalls         `json:"System Calls"`
	SMB2Calls   *SMBProfileSMB2Calls        `json:"SMB2 Calls"`
	Extended    map[string]*SMBProfileShare `json:"Extended Profile"`
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

// RunSMBStatusShares executes 'smbstatus --processes --json' on host
func RunSMBStatusProcesses() (*SMBStatus, error) {
	dat, err := executeSMBStatusCommand("--processes", "--json")
	if err != nil {
		return &SMBStatus{}, err
	}
	return parseSMBStatus(dat)
}

// RunSMBStatusShares executes 'smbstatus --shares --json' on host
func RunSMBStatusShares() (*SMBStatus, error) {
	dat, err := executeSMBStatusCommand("--shares", "--json")
	if err != nil {
		return &SMBStatus{}, err
	}
	return parseSMBStatus(dat)
}

// RunSMBStatusLocks executes 'smbstatus --locks --json' on host
func RunSMBStatusLocks() ([]SMBStatusOpenFile, error) {
	dat, err := executeSMBStatusCommand("--locks", "--json")
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

// RunSMBStatusProfile executes 'smbstatus --profile --json' on host
func RunSMBStatusProfile() (*SMBProfile, error) {
	dat, err := executeSMBStatusCommand("--profile", "--json")
	if err != nil {
		return &SMBProfile{}, err
	}
	return parseSMBProfile(dat)
}

// SMBStatusSharesByMachine converts the output of RunSMBStatusShares into map
// indexed by machine's host
func SMBStatusSharesByMachine() (map[string][]SMBStatusTreeCon, error) {
	smbstat, err := RunSMBStatusShares()
	if err != nil {
		return map[string][]SMBStatusTreeCon{}, err
	}
	return makeSmbSharesMap(smbstat.ListTreeCons()), nil
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
	}
	return &smbStatus
}

// parseSMBProfile parses to output of 'smbstatus --json --profile' into
// internal representation.
func parseSMBProfile(data string) (*SMBProfile, error) {
	res := NewSMBProfile()
	err := json.Unmarshal([]byte(data), res)
	return res, err
}

// NewSMBProfile returns non-populated SMBStatusProfile object
func NewSMBProfile() *SMBProfile {
	smbStatusProfile := SMBProfile{
		Timestamp: "",
		Version:   "",
		SmbConf:   "",
	}
	return &smbStatusProfile
}

// ListSessions returns a slice for mapped sessions
func (smbstat *SMBStatus) ListSessions() []SMBStatusSession {
	sessions := []SMBStatusSession{}
	for _, session := range smbstat.Sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// ListTreeCons returns a slice for mapped tree-connection
func (smbstat *SMBStatus) ListTreeCons() []SMBStatusTreeCon {
	tcons := []SMBStatusTreeCon{}
	for _, share := range smbstat.TCons {
		tcons = append(tcons, share)
	}
	return tcons
}

// ParseExtendedProfileKey parse the extended profile key into a pair of
// share-name and client-ip as string. Returns a pair of empty strings in case
// of parse failure.
func ParseExtendedProfileKey(key string) (shareName, clientIP string) {
	shareName = ""
	clientIP = ""
	sp := strings.Split(key, ":")
	if len(sp) != 2 {
		return
	}
	shareName = sp[0]
	sp = strings.Split(sp[1], "[")
	if len(sp) != 2 {
		return
	}
	clientIP = strings.Trim(sp[1], "[]")
	return
}
