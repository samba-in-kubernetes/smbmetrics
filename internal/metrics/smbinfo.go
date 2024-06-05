// SPDX-License-Identifier: Apache-2.0

package metrics

// SMBInfo provides a bridge layer between raw smbstatus info and exported
// metric counters. It also implements the more complex logic which requires in
// memory re-mapping of the low-level information (e.g., stats by machine/user).
type SMBInfo struct {
	smbstat *SMBStatus
}

func NewSMBInfo() *SMBInfo {
	return &SMBInfo{smbstat: NewSMBStatus()}
}

func NewUpdatedSMBInfo() (*SMBInfo, error) {
	smbinfo := NewSMBInfo()
	err := smbinfo.Update()
	return smbinfo, err
}

func (smbinfo *SMBInfo) Update() error {
	smbstat, err := RunSMBStatus()
	if err != nil {
		return err
	}
	smbinfo.smbstat = smbstat
	return nil
}

func (smbinfo *SMBInfo) TotalSessions() int {
	return len(smbinfo.smbstat.Sessions)
}

func (smbinfo *SMBInfo) TotalTreeCons() int {
	return len(smbinfo.smbstat.TCons)
}

func (smbinfo *SMBInfo) TotalOpenFiles() int {
	return len(smbinfo.smbstat.OpenFiles)
}

func (smbinfo *SMBInfo) TotalConnectedUsers() int {
	users := map[string]bool{}
	for _, session := range smbinfo.smbstat.Sessions {
		username := session.Username
		if len(username) > 0 {
			users[username] = true
		}
	}
	return len(users)
}

func (smbinfo *SMBInfo) MapMachineToSessions() map[string][]*SMBStatusSession {
	ret := map[string][]*SMBStatusSession{}
	for _, session := range smbinfo.smbstat.Sessions {
		machineID := session.RemoteMachine
		sessionRef := &session
		ret[machineID] = append(ret[machineID], sessionRef)
	}
	return ret
}

func (smbinfo *SMBInfo) MapServiceToTreeCons() map[string][]*SMBStatusTreeCon {
	ret := map[string][]*SMBStatusTreeCon{}
	for _, tcon := range smbinfo.smbstat.TCons {
		serviceID := tcon.Service
		tconRef := &tcon
		ret[serviceID] = append(ret[serviceID], tconRef)
	}
	return ret
}

func (smbinfo *SMBInfo) MapMachineToTreeCons() map[string][]*SMBStatusTreeCon {
	ret := map[string][]*SMBStatusTreeCon{}
	for _, tcon := range smbinfo.smbstat.TCons {
		machineID := tcon.Machine
		tconRef := &tcon
		ret[machineID] = append(ret[machineID], tconRef)
	}
	return ret
}

func (smbinfo *SMBInfo) MapServiceToMachines() map[string]map[string]int {
	ret := map[string]map[string]int{}
	for _, tcon := range smbinfo.smbstat.TCons {
		serviceID := tcon.Service
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
