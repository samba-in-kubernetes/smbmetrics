# Samba Metrics

Samba metrics exporter converts 'smbstatus' output into
[Prometheus](https://prometheus.io/)
[data-model](https://prometheus.io/docs/concepts/data_model/) metrics format.

## Build

```console
$ make build

$ make image-build
```


## Query metrics

When running (by privileged user) along-side active SMB server, `smbmetrics`
exports a set of gauge metrics over HTTP via port `9922`. Most metrics become
visible only when active SMB connections exists. When Samba is compiled and
run with profile-information enabled (`smb.conf` global section has
`smbd profiling level = on`), `smbmetrics` will also export various profile
stats as Prometheus metrics. Execute the following `curl` command on the same
machine where you run `smbmetrics` instance:

```console
$ curl --request GET "http://localhost:9922/metrics"
```

## Exported metrics

| Metric name               | Description                                      |
|---------------------------|--------------------------------------------------|
| `smb_metrics_status`      | Status and version of running process            |
| `smb_sessions_total`      | Number of active SMB sessions                    |
| `smb_tcon_total`          | Number of active SMB tree-connections            |
| `smb_users_total`         | Number of connected users                        |
| `smb_share_activity`      | Number of remote machines using each share       |
| `smb_share_byremote`      | Number of shares used by each remote machine     |


## Profile metrics

| Metric name               | Description                                      |
|---------------------------|--------------------------------------------------|
| `smb_smb2_request_total`  | Number of SMB2 requests                          |
| `smb_vfs_call_total`      | Number of calls to VFS layer                     |
| `smb_vfs_io_call_total`   | Number of I/O calls to VFS layer                 |

smb_vfs_call_total
## Example

The following example is from a setup with 2 shares and 2 users connected and
performing SMB file-system operations from 4 different machines:

```console
$ curl --request GET "http://localhost:9922/metrics"

# HELP smb_metrics_status Current metrics-collector status versions
# TYPE smb_metrics_status gauge
smb_metrics_status{commitid="092fe2bb0",ctdbvers="4.20.0-103",sambaimage="",sambavers="4.20.0-103",version="v0.2-28-g092fe2b"} 1
# HELP smb_sessions_total Number of currently active SMB sessions
# TYPE smb_sessions_total gauge
smb_sessions_total 8
# HELP smb_tcon_total Number of currently active SMB tree-connections
# TYPE smb_tcon_total gauge
smb_tcon_total 8
# HELP smb_users_total Number of currently active SMB users
# TYPE smb_users_total gauge
smb_users_total 2
# HELP smb_openfiles_total Number of currently open files
# TYPE smb_openfiles_total gauge
smb_openfiles_total 5
# HELP smb_openfiles_access_rw Number of open files with read-write access mode
# TYPE smb_openfiles_access_rw gauge
smb_openfiles_access_rw 4
# HELP smb_share_activity Number of remote machines currently using a share
# TYPE smb_share_activity gauge
smb_share_activity{service="smbshare1"} 4
smb_share_activity{service="smbshare2"} 2
# HELP smb_share_for_remote Number of shares served for remote machine
# TYPE smb_share_for_remote gauge
smb_share_byremote{machine="192.168.122.71"} 2
smb_share_byremote{machine="192.168.122.72"} 1
smb_share_byremote{machine="192.168.122.73"} 2
smb_share_byremote{machine="192.168.122.74"} 1
```

When running with profile enabled, we get also the following metrics:

```console
# HELP smb_smb2_request_total Total number of SMB2 requests
# TYPE smb_smb2_request_total gauge
smb_smb2_request_total{idle="0",inbytes="0",operation="break",outbytes="0",time="0"} 0
smb_smb2_request_total{idle="0",inbytes="0",operation="cancel",outbytes="0",time="0"} 0
smb_smb2_request_total{idle="0",inbytes="0",operation="lock",outbytes="0",time="0"} 0
smb_smb2_request_total{idle="0",inbytes="0",operation="logoff",outbytes="0",time="0"} 0
smb_smb2_request_total{idle="0",inbytes="0",operation="notify",outbytes="0",time="0"} 0
smb_smb2_request_total{idle="0",inbytes="0",operation="tdis",outbytes="0",time="0"} 0
smb_smb2_request_total{idle="0",inbytes="1156",operation="keepalive",outbytes="1156",time="98"} 17
smb_smb2_request_total{idle="0",inbytes="12995",operation="read",outbytes="232793072",time="102307"} 115
smb_smb2_request_total{idle="0",inbytes="16852",operation="getinfo",outbytes="27754",time="26869"} 162
smb_smb2_request_total{idle="0",inbytes="23069904",operation="write",outbytes="880",time="30958"} 11
smb_smb2_request_total{idle="0",inbytes="240",operation="negprot",outbytes="268",time="4049881"} 1
smb_smb2_request_total{idle="0",inbytes="240",operation="tcon",outbytes="160",time="97672"} 2
smb_smb2_request_total{idle="0",inbytes="26312",operation="close",outbytes="37712",time="101378"} 299
smb_smb2_request_total{idle="0",inbytes="295",operation="ioctl",outbytes="337",time="186"} 2
smb_smb2_request_total{idle="0",inbytes="3136",operation="setinfo",outbytes="1740",time="56382"} 26
smb_smb2_request_total{idle="0",inbytes="430",operation="sesssetup",outbytes="264",time="15242"} 2
smb_smb2_request_total{idle="0",inbytes="74352",operation="create",outbytes="82692",time="705058"} 363
smb_smb2_request_total{idle="0",inbytes="968",operation="flush",outbytes="748",time="1609"} 11
smb_smb2_request_total{idle="0",inbytes="980",operation="find",outbytes="3629",time="19929"} 10
# HELP smb_vfs_call_total Total number of calls to underlying VFS layer
# TYPE smb_vfs_call_total gauge
smb_vfs_call_total{operation="chdir",time="9045"} 209
smb_vfs_call_total{operation="chmod",time="0"} 0
smb_vfs_call_total{operation="close",time="12056"} 873
smb_vfs_call_total{operation="closedir",time="3738"} 83
smb_vfs_call_total{operation="createfile",time="0"} 0
smb_vfs_call_total{operation="fallocate",time="0"} 0
smb_vfs_call_total{operation="fchmod",time="603"} 11
smb_vfs_call_total{operation="fchown",time="0"} 0
smb_vfs_call_total{operation="fdopendir",time="1839"} 83
smb_vfs_call_total{operation="fntimes",time="14315"} 22
smb_vfs_call_total{operation="fstat",time="42921"} 1790
smb_vfs_call_total{operation="fstatat",time="0"} 0
smb_vfs_call_total{operation="ftruncate",time="0"} 0
smb_vfs_call_total{operation="getwd",time="7"} 2
smb_vfs_call_total{operation="lchown",time="0"} 0
smb_vfs_call_total{operation="linkat",time="0"} 0
smb_vfs_call_total{operation="lseek",time="0"} 0
smb_vfs_call_total{operation="lstat",time="9"} 1
smb_vfs_call_total{operation="mkdirat",time="2614"} 2
smb_vfs_call_total{operation="mknodat",time="0"} 0
smb_vfs_call_total{operation="open",time="0"} 0
smb_vfs_call_total{operation="openat",time="138723"} 1034
smb_vfs_call_total{operation="opendir",time="0"} 0
smb_vfs_call_total{operation="readdir",time="305110"} 534
smb_vfs_call_total{operation="readlinkat",time="0"} 0
smb_vfs_call_total{operation="realpath",time="17"} 4
smb_vfs_call_total{operation="renameat",time="22538"} 5
smb_vfs_call_total{operation="rewinddir",time="0"} 0
smb_vfs_call_total{operation="stat",time="42770"} 479
smb_vfs_call_total{operation="symlinkat",time="0"} 0
smb_vfs_call_total{operation="unlinkat",time="7843"} 10
# HELP smb_vfs_io_call_total Total number of I/O calls to underlying VFS layer
# TYPE smb_vfs_io_call_total gauge
smb_vfs_io_call_total{bytes="0",idle="0",operation="pread",time="0"} 0
smb_vfs_io_call_total{bytes="0",idle="0",operation="pwrite",time="0"} 0
smb_vfs_io_call_total{bytes="0",idle="14",operation="asys_fsync",time="1416"} 11
smb_vfs_io_call_total{bytes="23068672",idle="23",operation="asys_pwrite",time="29534"} 11
smb_vfs_io_call_total{bytes="232783872",idle="267",operation="asys_pread",time="99338"} 115
```
