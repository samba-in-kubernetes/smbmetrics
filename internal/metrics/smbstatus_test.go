// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, filename string) string {
	file, err := os.Open("testdata/" + filename)
	assert.NoError(t, err)
	defer file.Close()
	data, err := io.ReadAll(file)
	assert.NoError(t, err)
	return string(data)
}

func TestParseSMBStatusTCons(t *testing.T) {
	testdata1 := readTestData(t, "smbstatus-simple1.json")
	dat, err := parseSMBStatus(testdata1)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.TCons), 2)

	testdata2 := readTestData(t, "smbstatus-simple2.json")
	dat, err = parseSMBStatus(testdata2)
	assert.Equal(t, len(dat.TCons), 1)
	assert.NoError(t, err)
	tcons := dat.ListTreeCons()
	assert.Equal(t, len(tcons), 1)
	tcon1 := tcons[0]
	assert.Equal(t, tcon1.Service, "share1")
	assert.Equal(t, tcon1.ServerID.PID, "355")
	assert.Equal(t, tcon1.Machine, "::1")

	sharesMap := makeSmbSharesMap(tcons)
	assert.Equal(t, len(sharesMap), 1)
	for machine, share := range sharesMap {
		sharesCount := len(share)
		assert.Equal(t, sharesCount, 1)
		assert.Equal(t, machine, "::1")
	}
}

func TestParseSMBStatusAll(t *testing.T) {
	testdata1 := readTestData(t, "smbstatus-all1.json")
	dat, err := parseSMBStatus(testdata1)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.Sessions), 1)
	assert.Equal(t, len(dat.TCons), 1)

	testdata2 := readTestData(t, "smbstatus-all2.json")
	dat2, err := parseSMBStatusLocks(testdata2)
	assert.NoError(t, err)
	assert.Equal(t, len(dat2.OpenFiles), 2)
}

func TestParseSMBStatusLocks(t *testing.T) {
	testdata := readTestData(t, "smbstatus-locks.json")
	locks, err := parseSMBStatusLockedFiles(testdata)
	assert.NoError(t, err)
	assert.Equal(t, len(locks), 2)
	lock1 := locks[0]
	assert.Equal(t, lock1.FileID.Inode, int64(61))
	assert.Equal(t, lock1.NumPendingDeletes, 0)
	lock2 := locks[1]
	assert.Equal(t, lock2.FileID.Inode, int64(52))
	assert.Equal(t, lock2.NumPendingDeletes, 2)
}

func TestParseSMBStatusOpenFiles(t *testing.T) {
	testdata := readTestData(t, "smbstatus-openfiles.json")
	status, err := parseSMBStatusLocks(testdata)
	assert.NoError(t, err)
	assert.Equal(t, len(status.OpenFiles), 2)
	openFileAa := status.OpenFiles["/A/a"]
	assert.Equal(t, len(openFileAa.Opens), 2)
	for _, open := range openFileAa.Opens {
		oplock := open.OpLock
		lease := open.Lease
		assert.Equal(t, oplock.Batch, false)
		assert.Equal(t, oplock.LevelII, true)
		assert.Equal(t, oplock.Text, "LEVEL_II")
		assert.Equal(t, oplock.Exclusive, false)
		assert.Equal(t, lease.Handle, false)
		assert.Equal(t, lease.Read, false)
		assert.Equal(t, lease.Write, false)
		assert.Equal(t, lease.Text, "")
	}
	openFileAb := status.OpenFiles["/A/b"]
	assert.Equal(t, len(openFileAb.Opens), 2)
	for _, open := range openFileAb.Opens {
		oplock := open.OpLock
		lease := open.Lease
		if oplock.Lease {
			assert.Equal(t, oplock.Batch, false)
			assert.Equal(t, oplock.LevelII, false)
			assert.Equal(t, oplock.Text, "LEASE")
			assert.Equal(t, oplock.Exclusive, false)
			assert.Equal(t, lease.Handle, false)
			assert.Equal(t, lease.Read, true)
			assert.Equal(t, lease.Write, true)
			assert.Equal(t, lease.Text, "LEASE(RW)")
		} else {
			assert.Equal(t, oplock.Batch, false)
			assert.Equal(t, oplock.LevelII, true)
			assert.Equal(t, oplock.Text, "LEVEL_II")
			assert.Equal(t, oplock.Exclusive, false)
			assert.Equal(t, lease.Handle, false)
			assert.Equal(t, lease.Read, false)
			assert.Equal(t, lease.Write, false)
			assert.Equal(t, lease.Text, "")
		}
	}
}

func TestParseSMBStatusProfile(t *testing.T) {
	testdata := readTestData(t, "smbstatus-profile.json")
	profile, err := parseSMBProfile(testdata)
	assert.NoError(t, err)
	assert.Equal(t, profile.SmbdLoop.Connect.Count, 1)
	assert.Equal(t, profile.SmbdLoop.CPUSystem.Time, 667333)
	assert.Equal(t, profile.SmbdLoop.Request.Count, 1292)
	assert.Equal(t, profile.SystemCalls.AsysPRead.Count, 6)
	assert.Equal(t, profile.SystemCalls.AsysPWrite.Bytes, 90177536)
	assert.Equal(t, profile.SystemCalls.AsysFSync.Count, 47)
	assert.Equal(t, profile.SMB2Calls.Read.Outbytes, 10486240)
	assert.Equal(t, profile.SMB2Calls.Write.Inbytes, 90180784)
}

func TestParseSMBStatusProfileNoData(t *testing.T) {
	testdata := readTestData(t, "smbstatus-nodata.json")
	profile, err := parseSMBProfile(testdata)
	assert.NoError(t, err)
	assert.Nil(t, profile.SmbdLoop)
	assert.Nil(t, profile.SystemCalls)
	assert.Nil(t, profile.SMB2Calls)
}
