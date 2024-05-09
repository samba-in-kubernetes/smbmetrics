// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//revive:disable line-length-limit
//nolint:revive,lll
var (
	outputSmbStatusJSON = `
{
	"timestamp": "2022-07-19T16:26:34.652845+0530",
	"version": "4.17.0pre1-GIT-130283cbae0",
	"smb_conf": "/usr/local/etc/samba/smb.conf",
	"tcons": {
	  "2464814757": {
		"service": "gluster-vol",
		"server_id": {
		  "pid": "4214",
		  "task_id": "0",
		  "vnn": "0",
		  "unique_id": "18344755514750214344"
		},
		"tcon_id": "2464814757",
		"session_id": "659628098",
		"machine": "192.168.122.155",
		"connected_at": "2022-07-19T16:20:37+0530",
		"encryption": {
		  "cipher": "",
		  "degree": "none"
		},
		"signing": {
		  "cipher": "",
		  "degree": "none"
		}
	  },
	  "2542351833": {
		"service": "gluster-vol",
		"server_id": {
		  "pid": "5299",
		  "task_id": "0",
		  "vnn": "1",
		  "unique_id": "13525406857402943822"
		},
		"tcon_id": "2542351833",
		"session_id": "687866031",
		"machine": "192.168.122.1",
		"connected_at": "2022-07-19T16:24:52+0530",
		"encryption": {
		  "cipher": "",
		  "degree": "none"
		},
		"signing": {
		  "cipher": "",
		  "degree": "none"
		}
	  }
	}
  }
  `

	outputSmbStatusJSON2 = `
  {
	"timestamp": "2023-06-07T11:49:05.528375+0000",
	"version": "4.17.8",
	"smb_conf": "/etc/samba/smb.conf",
	"tcons": {
	  "2295102631": {
		"service": "share1",
		"server_id": {
		  "pid": "355",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "4904044827525949167"
		},
		"tcon_id": "2295102631",
		"session_id": "2875824865",
		"machine": "::1",
		"connected_at": "2023-06-07T11:44:56.766022+00:00",
		"encryption": {
		  "cipher": "-",
		  "degree": "none"
		},
		"signing": {
		  "cipher": "-",
		  "degree": "none"
		}
	  }
	}
  }`

	outputSmbStatusAllJSON = `
  {
	"timestamp": "2022-04-15T18:25:15.364891+0200",
	"version": "4.17.0pre1-GIT-a0f12b9c80b",
	"smb_conf": "/opt/sambaTest/etc/smb.conf",
	"sessions": {
	  "3639217376": {
		"session_id": "3639217376",
		"server_id": {
		  "pid": "69650",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "10756714984493602300"
		},
		"uid": 1000,
		"gid": 1000,
		"username": "janger",
		"groupname": "janger",
		"remote_machine": "127.0.0.1",
		"hostname": "ipv4:127.0.0.1:59944",
		"session_dialect": "SMB3_11",
		"encryption": {
		  "cipher": "",
		  "degree": "none"
		},
		"signing": {
		  "cipher": "AES-128-GMAC",
		  "degree": "partial"
		}
	  }
	},
	"tcons": {
	  "3813255619": {
		"service": "gemeinsam",
		"server_id": {
		  "pid": "69650",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "10756714984493602300"
		},
		"tcon_id": "3813255619",
		"session_id": "3639217376",
		"machine": "127.0.0.1",
		"connected_at": "2022-04-15T17:30:37+0200",
		"encryption": {
		  "cipher": "AES-128-GMAC",
		  "degree": "full"
		},
		"signing": {
		  "cipher": "",
		  "degree": "none"
		}
	  }
	},
	"locked_files": {
	  "/home/janger/testfolder/hallo": {
		"service_path": "/home/janger/testfolder",
		"filename": "hallo",
		"fileid": {
		  "devid": 59,
		  "inode": 11404245,
		  "extid": 0
		},
		"num_pending_deletes": 0,
		"opens": {
		  "56839/2": {
			"server_id": {
			  "pid": "69650",
			  "task_id": "0",
			  "vnn": "4294967295",
			  "unique_id": "10756714984493602300"
			},
			"uid": 1000,
			"share_file_id": 2,
			"sharemode": {
			  "hex": "0x00000003",
			  "NONE": false,
			  "READ": true,
			  "WRITE": true,
			  "DELETE": false,
			  "text": "RW"
			},
			"access_mask": {
			  "hex": "0x00000003",
			  "READ_DATA": true,
			  "WRITE_DATA": true,
			  "APPEND_DATA": false,
			  "READ_EA": false,
			  "WRITE_EA": false,
			  "EXECUTE": false,
			  "READ_ATTRIBUTES": false,
			  "WRITE_ATTRIBUTES": false,
			  "DELETE_CHILD": false,
			  "DELETE": false,
			  "READ_CONTROL": false,
			  "WRITE_DAC": false,
			  "SYNCHRONIZE": false,
			  "ACCESS_SYSTEM_SECURITY": false,
			  "text": "RW"
			},
			"caching": {
			  "READ": false,
			  "WRITE": false,
			  "HANDLE": false,
			  "hex": "0x00000000",
			  "text": ""
			},
			"oplock": {
			  "EXCLUSIVE": false,
			  "BATCH": false,
			  "LEVEL_II": false,
			  "LEASE": false,
			  "text": "NONE"
			},
			"lease": {},
			"connected_at": "2022-04-15T17:30:38+0200"
		  }
		}
	  }
	}
  }

  `

	outputSmbStatusAllJSON2 = `
  {
	"timestamp": "2022-07-20T12:07:36.225955+0000",
	"version": "4.17.0pre1-UNKNOWN",
	"smb_conf": "/etc/samba/smb.conf",
	"sessions": {},
	"tcons": {
	  "348413079": {
		"service": "IPC$",
		"server_id": {
		  "pid": "101",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "7364797719700910696"
		},
		"tcon_id": "348413079",
		"session_id": "674813472",
		"machine": "127.0.0.1",
		"connected_at": "2022-07-20T12:04:15+0000",
		"encryption": {
		  "cipher": "",
		  "degree": "none"
		},
		"signing": {
		  "cipher": "",
		  "degree": "none"
		}
	  },
	  "1698398697": {
		"service": "samba-share",
		"server_id": {
		  "pid": "101",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "7364797719700910696"
		},
		"tcon_id": "1698398697",
		"session_id": "674813472",
		"machine": "127.0.0.1",
		"connected_at": "2022-07-20T12:04:15+0000",
		"encryption": {
		  "cipher": "",
		  "degree": "none"
		},
		"signing": {
		  "cipher": "",
		  "degree": "none"
		}
	  }
	},
	"locked_files": {
	  "/mnt/96dd85fd-6c60-409c-bc1c-15f98eb358ee/a/y": {
		"service_path": "/mnt/96dd85fd-6c60-409c-bc1c-15f98eb358ee",
		"filename": "a/y",
		"fileid": {
		  "devid": -397762331,
		  "inode": 13631494,
		  "extid": 0
		},
		"num_pending_deletes": 0,
		"opens": {
		  "101/61": {
			"server_id": {
			  "pid": "101",
			  "task_id": "0",
			  "vnn": "4294967295",
			  "unique_id": "7364797719700910696"
			},
			"uid": 1000,
			"share_file_id": 61,
			"sharemode": {
			  "hex": "0x00000007",
			  "NONE": false,
			  "READ": true,
			  "WRITE": true,
			  "DELETE": true,
			  "text": "RWD"
			},
			"access_mask": {
			  "hex": "0x0012019f",
			  "READ_DATA": true,
			  "WRITE_DATA": true,
			  "APPEND_DATA": true,
			  "READ_EA": true,
			  "WRITE_EA": true,
			  "EXECUTE": false,
			  "READ_ATTRIBUTES": true,
			  "WRITE_ATTRIBUTES": true,
			  "DELETE_CHILD": false,
			  "DELETE": false,
			  "READ_CONTROL": true,
			  "WRITE_DAC": false,
			  "SYNCHRONIZE": true,
			  "ACCESS_SYSTEM_SECURITY": false,
			  "text": "RW"
			},
			"caching": {
			  "READ": true,
			  "WRITE": true,
			  "HANDLE": true,
			  "hex": "0x00000007",
			  "text": "RWH"
			},
			"oplock": {
			  "EXCLUSIVE": false,
			  "BATCH": false,
			  "LEVEL_II": false,
			  "LEASE": true,
			  "text": "LEASE(RWH)"
			},
			"lease": {
			  "lease_key": "ac1cf117-4ac6-b543-8bc8-597d795e8546",
			  "hex": "0x00000007",
			  "READ": true,
			  "WRITE": true,
			  "HANDLE": true,
			  "text": "RWH"
			},
			"connected_at": "2022-07-20T12:06:29+0000"
		  }
		}
	  },
	  "/mnt/96dd85fd-6c60-409c-bc1c-15f98eb358ee/a/x": {
		"service_path": "/mnt/96dd85fd-6c60-409c-bc1c-15f98eb358ee",
		"filename": "a/x",
		"fileid": {
		  "devid": -397762331,
		  "inode": 13631493,
		  "extid": 0
		},
		"num_pending_deletes": 0,
		"opens": {
		  "101/65": {
			"server_id": {
			  "pid": "101",
			  "task_id": "0",
			  "vnn": "4294967295",
			  "unique_id": "7364797719700910696"
			},
			"uid": 1000,
			"share_file_id": 65,
			"sharemode": {
			  "hex": "0x00000007",
			  "NONE": false,
			  "READ": true,
			  "WRITE": true,
			  "DELETE": true,
			  "text": "RWD"
			},
			"access_mask": {
			  "hex": "0x0012019f",
			  "READ_DATA": true,
			  "WRITE_DATA": true,
			  "APPEND_DATA": true,
			  "READ_EA": true,
			  "WRITE_EA": true,
			  "EXECUTE": false,
			  "READ_ATTRIBUTES": true,
			  "WRITE_ATTRIBUTES": true,
			  "DELETE_CHILD": false,
			  "DELETE": false,
			  "READ_CONTROL": true,
			  "WRITE_DAC": false,
			  "SYNCHRONIZE": true,
			  "ACCESS_SYSTEM_SECURITY": false,
			  "text": "RW"
			},
			"caching": {
			  "READ": true,
			  "WRITE": true,
			  "HANDLE": true,
			  "hex": "0x00000007",
			  "text": "RWH"
			},
			"oplock": {
			  "EXCLUSIVE": false,
			  "BATCH": false,
			  "LEVEL_II": false,
			  "LEASE": true,
			  "text": "LEASE(RWH)"
			},
			"lease": {
			  "lease_key": "a110aee8-2bfa-2349-b148-22a018a2e061",
			  "hex": "0x00000007",
			  "READ": true,
			  "WRITE": true,
			  "HANDLE": true,
			  "text": "RWH"
			},
			"connected_at": "2022-07-20T12:06:42+0000"
		  }
		}
	  }
	}
  }
  `

	outputSmbStatusLocksJSON = `
  {
	"timestamp": "2024-04-14T14:53:34.901974+0300",
	"version": "4.21.0pre1-GIT-58a018fb7ad",
	"smb_conf": "//etc/samba/smb.conf",
	"open_files": {
	  "/A/A2/A6/r1": {
	    "service_path": "/",
	    "filename": "A/A2/A6/r1",
	    "fileid": {
	      "devid": 1,
	      "inode": 61,
	      "extid": 0
	    },
	    "num_pending_deletes": 0,
	    "opens": {
	      "1790/261": {
		"server_id": {
		  "pid": "1790",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "3607086338167075363"
		},
		"uid": 2222,
		"share_file_id": "261",
		"sharemode": {
		  "hex": "0x00000007",
		  "READ": true,
		  "WRITE": true,
		  "DELETE": true,
		  "text": "RWD"
		},
		"access_mask": {
		  "hex": "0x00120089",
		  "READ_DATA": true,
		  "WRITE_DATA": false,
		  "APPEND_DATA": false,
		  "READ_EA": true,
		  "WRITE_EA": false,
		  "EXECUTE": false,
		  "READ_ATTRIBUTES": true,
		  "WRITE_ATTRIBUTES": false,
		  "DELETE_CHILD": false,
		  "DELETE": false,
		  "READ_CONTROL": true,
		  "WRITE_DAC": false,
		  "SYNCHRONIZE": true,
		  "ACCESS_SYSTEM_SECURITY": false,
		  "text": "R"
		},
		"caching": {
		  "READ": true,
		  "WRITE": true,
		  "HANDLE": false,
		  "hex": "0x00000005",
		  "text": "RW"
		},
		"oplock": {
		  "EXCLUSIVE": true,
		  "BATCH": true,
		  "LEVEL_II": false,
		  "LEASE": false,
		  "text": "BATCH"
		},
		"lease": {},
		"opened_at": "2024-04-14T14:53:15.569085+03:00"
	      }
	    }
	  },
	  "/A/A1/r2": {
	    "service_path": "/",
	    "filename": "A/A1/r2",
	    "fileid": {
	      "devid": 2,
	      "inode": 52,
	      "extid": 0
	    },
	    "num_pending_deletes": 2,
	    "opens": {
	      "1790/267": {
		"server_id": {
		  "pid": "1790",
		  "task_id": "0",
		  "vnn": "4294967295",
		  "unique_id": "3607086338167075363"
		},
		"uid": 1111,
		"share_file_id": "222",
		"sharemode": {
		  "hex": "0x00000007",
		  "READ": true,
		  "WRITE": true,
		  "DELETE": true,
		  "text": "RWD"
		},
		"access_mask": {
		  "hex": "0x00120089",
		  "READ_DATA": true,
		  "WRITE_DATA": false,
		  "APPEND_DATA": false,
		  "READ_EA": true,
		  "WRITE_EA": false,
		  "EXECUTE": false,
		  "READ_ATTRIBUTES": true,
		  "WRITE_ATTRIBUTES": false,
		  "DELETE_CHILD": false,
		  "DELETE": false,
		  "READ_CONTROL": true,
		  "WRITE_DAC": false,
		  "SYNCHRONIZE": true,
		  "ACCESS_SYSTEM_SECURITY": false,
		  "text": "R"
		},
		"caching": {
		  "READ": true,
		  "WRITE": true,
		  "HANDLE": false,
		  "hex": "0x00000005",
		  "text": "RW"
		},
		"oplock": {
		  "EXCLUSIVE": true,
		  "BATCH": true,
		  "LEVEL_II": false,
		  "LEASE": false,
		  "text": "BATCH"
		},
		"lease": {},
		"opened_at": "2024-04-14T14:53:32.258325+03:00"
	      }
	    }
	  }
	}
      }
  `
)

//revive:enable line-length-limit

func TestParseSmbStatusSharesJSON(t *testing.T) {
	dat, err := parseSmbStatusJSON(outputSmbStatusJSON)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.TCons), 2)

	dat, err = parseSmbStatusJSON(outputSmbStatusJSON2)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.TCons), 1)

	shares, err := parseSmbStatusSharesAsJSON(outputSmbStatusJSON2)
	assert.NoError(t, err)
	assert.Equal(t, len(shares), 1)
	share1 := shares[0]
	assert.Equal(t, share1.Service, "share1")
	assert.Equal(t, share1.ServerID.PID, "355")
	assert.Equal(t, share1.Machine, "::1")

	sharesMap := makeSmbSharesMap(shares)
	assert.Equal(t, len(sharesMap), 1)
	for machine, share := range sharesMap {
		sharesCount := len(share)
		assert.Equal(t, sharesCount, 1)
		assert.Equal(t, machine, "::1")
	}
}

func TestParseSmbStatusAllJSON(t *testing.T) {
	dat, err := parseSmbStatusJSON(outputSmbStatusAllJSON)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.Sessions), 1)
	assert.Equal(t, len(dat.TCons), 1)
	assert.Equal(t, len(dat.LockedFiles), 1)

	dat2, err := parseSmbStatusJSON(outputSmbStatusAllJSON2)
	assert.NoError(t, err)
	assert.Equal(t, len(dat2.LockedFiles), 2)
}

func TestParseSmbStatusLocksJSON(t *testing.T) {
	locks, err := parseSmbStatusLocksAsJSON(outputSmbStatusLocksJSON)
	assert.NoError(t, err)
	assert.Equal(t, len(locks), 2)
	lock1 := locks[0]
	assert.Equal(t, lock1.FileID.Inode, int64(61))
	assert.Equal(t, lock1.NumPendingDeletes, 0)
	lock2 := locks[1]
	assert.Equal(t, lock2.FileID.Inode, int64(52))
	assert.Equal(t, lock2.NumPendingDeletes, 2)
}
