// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"github.com/go-logr/logr"
)

// SMBInfo provides a bridge layer between raw smbstatus info and exported
// metric counters. It also implements the more complex logic which requires in
// memory re-mapping of the low-level information (e.g., stats by machine/user).
type SMBInfo struct {
	tconsStatus    *SMBStatus
	sessionsStatus *SMBStatus
	log            logr.Logger
}

func NewSMBInfo(log logr.Logger) *SMBInfo {
	return &SMBInfo{
		tconsStatus:    NewSMBStatus(),
		sessionsStatus: NewSMBStatus(),
		log:            log,
	}
}

func NewUpdatedSMBInfo(log logr.Logger) (*SMBInfo, error) {
	smbinfo := NewSMBInfo(log)
	err := smbinfo.Update()
	return smbinfo, err
}

func (smbinfo *SMBInfo) Update() error {
	tconsStatus, err := RunSMBStatusShares()
	if err != nil {
		smbinfo.log.Error(err, "smbsstatus --shares failed")
		return err
	}
	sessionsStatus, err := RunSMBStatusProcesses()
	if err != nil {
		smbinfo.log.Error(err, "smbsstatus --processes failed")
		return err
	}
	smbinfo.tconsStatus = tconsStatus
	smbinfo.sessionsStatus = sessionsStatus
	return nil
}

func (smbinfo *SMBInfo) TotalSessions() int {
	return len(smbinfo.sessionsStatus.Sessions)
}

func (smbinfo *SMBInfo) TotalTreeCons() int {
	total := 0
	for _, tcon := range smbinfo.tconsStatus.TCons {
		serviceID := tcon.Service
		if isInternalServiceID(serviceID) {
			continue
		}
		total++
	}
	return total
}

func (smbinfo *SMBInfo) TotalConnectedUsers() int {
	users := map[string]bool{}
	for _, session := range smbinfo.sessionsStatus.Sessions {
		username := session.Username
		if len(username) > 0 {
			users[username] = true
		}
	}
	return len(users)
}

func (smbinfo *SMBInfo) MapMachineToSessions() map[string][]*SMBStatusSession {
	ret := map[string][]*SMBStatusSession{}
	for _, session := range smbinfo.sessionsStatus.Sessions {
		machineID := session.RemoteMachine
		sessionRef := &session
		ret[machineID] = append(ret[machineID], sessionRef)
	}
	return ret
}

func (smbinfo *SMBInfo) MapServiceToTreeCons() map[string][]*SMBStatusTreeCon {
	ret := map[string][]*SMBStatusTreeCon{}
	for _, tcon := range smbinfo.tconsStatus.TCons {
		serviceID := tcon.Service
		if isInternalServiceID(serviceID) {
			continue
		}
		tconRef := &tcon
		ret[serviceID] = append(ret[serviceID], tconRef)
	}
	return ret
}

func (smbinfo *SMBInfo) MapMachineToTreeCons() map[string][]*SMBStatusTreeCon {
	ret := map[string][]*SMBStatusTreeCon{}
	for _, tcon := range smbinfo.tconsStatus.TCons {
		serviceID := tcon.Service
		if isInternalServiceID(serviceID) {
			continue
		}
		machineID := tcon.Machine
		tconRef := &tcon
		ret[machineID] = append(ret[machineID], tconRef)
	}
	return ret
}

func (smbinfo *SMBInfo) MapServiceToMachines() map[string]map[string]int {
	ret := map[string]map[string]int{}
	for _, tcon := range smbinfo.tconsStatus.TCons {
		serviceID := tcon.Service
		if isInternalServiceID(serviceID) {
			continue
		}
		machineID := tcon.Machine
		sub, found := ret[serviceID]
		if !found {
			ret[serviceID] = map[string]int{machineID: 1}
		} else {
			sub[machineID]++
		}
	}
	return ret
}

func (smbinfo *SMBInfo) MapMachineToServies() map[string]map[string]int {
	ret := map[string]map[string]int{}
	for _, tcon := range smbinfo.tconsStatus.TCons {
		serviceID := tcon.Service
		if isInternalServiceID(serviceID) {
			continue
		}
		machineID := tcon.Machine
		sub, found := ret[machineID]
		if !found {
			ret[machineID] = map[string]int{serviceID: 1}
		} else {
			sub[serviceID]++
		}
	}
	return ret
}

func isInternalServiceID(serviceID string) bool {
	return serviceID == "IPC$"
}

// SMBProfileInfo provides a bridge layer between raw smbstatus profile info and
// exported metric counters.
type SMBProfileInfo struct {
	profileStatus *SMBProfile
	log           logr.Logger
}

func NewSMBProfileInfo(log logr.Logger) *SMBProfileInfo {
	return &SMBProfileInfo{
		profileStatus: NewSMBProfile(),
		log:           log,
	}
}

func NewUpdatedSMBProfileInfo(log logr.Logger) (*SMBProfileInfo, error) {
	smbProfileInfo := NewSMBProfileInfo(log)
	err := smbProfileInfo.Update()
	return smbProfileInfo, err
}

func (smbProfileInfo *SMBProfileInfo) Update() error {
	profiuleStatus, err := RunSMBStatusProfile()
	if err != nil {
		smbProfileInfo.log.Error(err, "smbsstatus --profile failed")
		return err
	}
	smbProfileInfo.profileStatus = profiuleStatus
	return nil
}
