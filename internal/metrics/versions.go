// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"context"
	"errors"
)

var (
	defaultVersions Versions
)

type Versions struct {
	Version      string
	CommitID     string
	SambaImage   string
	SambaVersion string
	CtdbVersion  string
}

// UpdateDefaultVersions assigns defaults upon init
func UpdateDefaultVersions(version, commitid string) {
	defaultVersions.Version = version
	defaultVersions.CommitID = commitid
}

// ResolveVersions is a best-effort to resolve current pod's versions info
func ResolveVersions(clnt *kclient) (Versions, error) {
	var imgErr, smbVersErr, ctdbVersErr error
	vers := Versions{
		Version:  defaultVersions.Version,
		CommitID: defaultVersions.CommitID,
	}
	if clnt != nil {
		vers.SambaImage, imgErr = resolveSambaImage(clnt)
	}
	vers.SambaVersion, smbVersErr = resolveSambaVersion()
	vers.CtdbVersion, ctdbVersErr = resolveCtdbVersion()
	return vers, errors.Join(imgErr, smbVersErr, ctdbVersErr)
}

func resolveSambaImage(clnt *kclient) (string, error) {
	pod, err := GetSelfPod(context.TODO(), clnt)
	if err != nil {
		return "", err
	}
	for _, cont := range pod.Spec.Containers {
		if cont.Name == "samba" {
			return cont.Image, nil
		}
	}
	return "", nil
}

func resolveSambaVersion() (string, error) {
	return executeRpmQCommand("samba")
}

func resolveCtdbVersion() (string, error) {
	return executeRpmQCommand("ctdb")
}

func executeRpmQCommand(name string) (string, error) {
	return executeCommand("rpm", "-q", name)
}
