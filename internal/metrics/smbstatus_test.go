// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//revive:disable line-length-limit
//nolint:revive,lll
var (
	smbstatusOutput1 = `
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

	smbstatusOutput2 = `
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

	smbstatusOutput3 = `
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
	"open_files": {
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
			"share_file_id": "2",
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

	smbstatusOutput4 = `
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
	"open_files": {
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
			"share_file_id": "61",
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
			  "BATCH": true,
			  "LEVEL_II": false,
			  "LEASE": false,
			  "text": "BATCH"
			},
			"lease": {},
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
			"share_file_id": "65",
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
			  "BATCH": true,
			  "LEVEL_II": false,
			  "LEASE": true,
			  "text": "BATCH"
			},
			"lease": {},
			"connected_at": "2022-07-20T12:06:42+0000"
		  }
		}
	  }
	}
  }
  `

	smbstatusOutput5 = `
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
		  "EXCLUSIVE": false,
		  "BATCH": false,
		  "LEVEL_II": true,
		  "LEASE": false,
		  "text": "BATCH"
		},
		"opened_at": "2024-04-14T14:53:32.258325+03:00"
	      }
	    }
	  }
	}
      }
  `

	smbstatusOutput6 = `
{
  "timestamp": "2024-07-04T13:04:45.910759+0300",
  "version": "4.21.0pre1-GIT-a5c7776b2f9",
  "smb_conf": "//etc/samba/smb.conf",
  "sessions": {
    "2211197710": {
      "session_id": "2211197710",
      "server_id": {
        "pid": "34128",
        "task_id": "0",
        "vnn": "4294967295",
        "unique_id": "16540149349904229747"
      },
      "uid": 1111,
      "gid": 1111,
      "username": "testuser",
      "groupname": "testuser",
      "creation_time": "2024-07-04T12:44:14.413940+03:00",
      "expiration_time": "30828-09-14T05:48:05.477581+03:00",
      "auth_time": "2024-07-04T12:44:14.418769+03:00",
      "remote_machine": "192.168.122.83",
      "hostname": "ipv4:192.168.122.83:56892",
      "session_dialect": "SMB2_02",
      "client_guid": "00000000-0000-0000-0000-000000000000",
      "encryption": {
        "cipher": "-",
        "degree": "none"
      },
      "signing": {
        "cipher": "-",
        "degree": "none"
      },
      "channels": {
        "1": {
          "channel_id": "1",
          "creation_time": "2024-07-04T12:44:14.413940+03:00",
          "local_address": "ipv4:192.168.122.119:445",
          "remote_address": "ipv4:192.168.122.83:56892"
        }
      }
    },
    "664362732": {
      "session_id": "664362732",
      "server_id": {
        "pid": "34156",
        "task_id": "0",
        "vnn": "4294967295",
        "unique_id": "10342000937013985296"
      },
      "uid": 1111,
      "gid": 1111,
      "username": "testuser",
      "groupname": "testuser",
      "creation_time": "2024-07-04T12:48:20.117703+03:00",
      "expiration_time": "30828-09-14T05:48:05.477581+03:00",
      "auth_time": "2024-07-04T12:48:20.124245+03:00",
      "remote_machine": "192.168.122.235",
      "hostname": "ipv4:192.168.122.235:34092",
      "session_dialect": "SMB2_02",
      "client_guid": "00000000-0000-0000-0000-000000000000",
      "encryption": {
        "cipher": "-",
        "degree": "none"
      },
      "signing": {
        "cipher": "-",
        "degree": "none"
      },
      "channels": {
        "1": {
          "channel_id": "1",
          "creation_time": "2024-07-04T12:48:20.117703+03:00",
          "local_address": "ipv4:192.168.122.119:445",
          "remote_address": "ipv4:192.168.122.235:34092"
        }
      }
    }
  },
  "tcons": {
    "2459296875": {
      "service": "IPC$",
      "server_id": {
        "pid": "34156",
        "task_id": "0",
        "vnn": "4294967295",
        "unique_id": "10342000937013985296"
      },
      "tcon_id": "2459296875",
      "session_id": "664362732",
      "machine": "192.168.122.235",
      "connected_at": "2024-07-04T12:48:20.129126+03:00",
      "encryption": {
        "cipher": "-",
        "degree": "none"
      },
      "signing": {
        "cipher": "-",
        "degree": "none"
      }
    },
    "1357200611": {
      "service": "smbshare",
      "server_id": {
        "pid": "34156",
        "task_id": "0",
        "vnn": "4294967295",
        "unique_id": "10342000937013985296"
      },
      "tcon_id": "1357200611",
      "session_id": "664362732",
      "machine": "192.168.122.235",
      "connected_at": "2024-07-04T12:48:20.130103+03:00",
      "encryption": {
        "cipher": "-",
        "degree": "none"
      },
      "signing": {
        "cipher": "-",
        "degree": "none"
      }
    },
    "2373869966": {
      "service": "smbshare",
      "server_id": {
        "pid": "34128",
        "task_id": "0",
        "vnn": "4294967295",
        "unique_id": "16540149349904229747"
      },
      "tcon_id": "2373869966",
      "session_id": "2211197710",
      "machine": "192.168.122.83",
      "connected_at": "2024-07-04T12:44:14.424733+03:00",
      "encryption": {
        "cipher": "-",
        "degree": "none"
      },
      "signing": {
        "cipher": "-",
        "degree": "none"
      }
    },
    "2758374422": {
      "service": "IPC$",
      "server_id": {
        "pid": "34128",
        "task_id": "0",
        "vnn": "4294967295",
        "unique_id": "16540149349904229747"
      },
      "tcon_id": "2758374422",
      "session_id": "2211197710",
      "machine": "192.168.122.83",
      "connected_at": "2024-07-04T12:44:14.423739+03:00",
      "encryption": {
        "cipher": "-",
        "degree": "none"
      },
      "signing": {
        "cipher": "-",
        "degree": "none"
      }
    }
  },
  "open_files": {
    "/A/b": {
      "service_path": "/",
      "filename": "A/b",
      "fileid": {
        "devid": -2,
        "inode": 1099512325566,
        "extid": 0
      },
      "num_pending_deletes": 0,
      "opens": {
        "34156/114": {
          "server_id": {
            "pid": "34156",
            "task_id": "0",
            "vnn": "4294967295",
            "unique_id": "10342000937013985296"
          },
          "uid": 1111,
          "share_file_id": "114",
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
            "WRITE": false,
            "HANDLE": false,
            "hex": "0x00000001",
            "text": "R"
          },
          "oplock": {
            "EXCLUSIVE": false,
            "BATCH": false,
            "LEVEL_II": false,
            "LEASE": true,
            "text": "LEASE"
          },
          "lease": {
            "lease_key": "272e4282-36e6-11ef-8a34-309c2337f855",
            "hex": "0x00000005",
            "READ": true,
            "WRITE": true,
            "HANDLE": false,
            "text": "LEASE(RW)"
          },
          "opened_at": "2024-07-04T13:02:41.967466+03:00"
        },
        "34128/309": {
          "server_id": {
            "pid": "34128",
            "task_id": "0",
            "vnn": "4294967295",
            "unique_id": "16540149349904229747"
          },
          "uid": 1111,
          "share_file_id": "309",
          "sharemode": {
            "hex": "0x00000007",
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
            "WRITE": false,
            "HANDLE": false,
            "hex": "0x00000001",
            "text": "R"
          },
          "oplock": {
            "EXCLUSIVE": false,
            "BATCH": false,
            "LEVEL_II": true,
            "LEASE": false,
            "text": "LEVEL_II"
          },
          "lease": {},
          "opened_at": "2024-07-04T13:04:05.727421+03:00"
        }
      }
    },
    "/A/a": {
      "service_path": "/",
      "filename": "A/a",
      "fileid": {
        "devid": -2,
        "inode": 1099512188451,
        "extid": 0
      },
      "num_pending_deletes": 0,
      "opens": {
        "34156/110": {
          "server_id": {
            "pid": "34156",
            "task_id": "0",
            "vnn": "4294967295",
            "unique_id": "10342000937013985296"
          },
          "uid": 1111,
          "share_file_id": "110",
          "sharemode": {
            "hex": "0x00000007",
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
            "WRITE": false,
            "HANDLE": false,
            "hex": "0x00000001",
            "text": "R"
          },
          "oplock": {
            "EXCLUSIVE": false,
            "BATCH": false,
            "LEVEL_II": true,
            "LEASE": false,
            "text": "LEVEL_II"
          },
          "lease": {},
          "opened_at": "2024-07-04T13:02:18.430418+03:00"
        },
        "34128/303": {
          "server_id": {
            "pid": "34128",
            "task_id": "0",
            "vnn": "4294967295",
            "unique_id": "16540149349904229747"
          },
          "uid": 1111,
          "share_file_id": "303",
          "sharemode": {
            "hex": "0x00000007",
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
            "WRITE": false,
            "HANDLE": false,
            "hex": "0x00000001",
            "text": "R"
          },
          "oplock": {
            "EXCLUSIVE": false,
            "BATCH": false,
            "LEVEL_II": true,
            "LEASE": false,
            "text": "LEVEL_II"
          },
          "lease": {},
          "opened_at": "2024-07-04T13:02:04.739655+03:00"
        }
      }
    }
  }
}
`
)

//revive:enable line-length-limit

func TestParseSMBStatusTCons(t *testing.T) {
	dat, err := parseSMBStatus(smbstatusOutput1)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.TCons), 2)

	dat, err = parseSMBStatus(smbstatusOutput2)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.TCons), 1)

	shares, err := parseSMBStatusTreeCons(smbstatusOutput2)
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

func TestParseSMBStatusAll(t *testing.T) {
	dat, err := parseSMBStatus(smbstatusOutput3)
	assert.NoError(t, err)
	assert.Equal(t, len(dat.Sessions), 1)
	assert.Equal(t, len(dat.TCons), 1)
	assert.Equal(t, len(dat.OpenFiles), 1)

	dat2, err := parseSMBStatus(smbstatusOutput4)
	assert.NoError(t, err)
	assert.Equal(t, len(dat2.OpenFiles), 2)
}

func TestParseSMBStatusLocks(t *testing.T) {
	locks, err := parseSMBStatusLockedFiles(smbstatusOutput5)
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
	status, err := parseSMBStatus(smbstatusOutput6)
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
