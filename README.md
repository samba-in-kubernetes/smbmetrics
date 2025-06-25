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


## Profile metrics (per operation)

| Metric name                                  | Description                                              |
|----------------------------------------------|----------------------------------------------------------|
| `smb_smb2_request_total`                     | Total number of SMB2 requests                            |
| `smb_smb2_request_inbytes`                   | Bytes received for SMB2 requests                         |
| `smb_smb2_request_outbytes`                  | Bytes replied for SMB2 requests                          |
| `smb_smb2_request_duration_microseconds_sum` | Execution time in microseconds of SMB2 requests          |
| `smb_vfs_total`                              | Total number of calls to underlying VFS layer            |
| `smb_vfs_io_total`                           | Total number of I/O calls to underlying VFS layer        |
| `smb_vfs_io_bytes`                           | Number of bytes transferred via underlying VFS I/O layer |
| `smb_vfs_io_duration_microseconds_sum`       | Execution time in microseconds of VFS I/O requests       |


## Example

The following example is from a setup with 2 shares and 2 users connected and
performing SMB file-system operations from 4 different machines:

```console
$ curl --request GET "http://localhost:9922/metrics"

# HELP smb_metrics_status Current metrics-collector status versions
# TYPE smb_metrics_status gauge
smb_metrics_status{commitid="092fe2bb0",ctdbvers="4.20.0-103",netbiosname="cluster1",sambaimage="",sambavers="4.20.0-103",version="v0.2-28-g092fe2b"} 1
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
smb_smb2_request_total{operation="break"} 0
smb_smb2_request_total{operation="cancel"} 0
smb_smb2_request_total{operation="close"} 1347
smb_smb2_request_total{operation="create"} 3378
smb_smb2_request_total{operation="find"} 394
smb_smb2_request_total{operation="flush"} 15
smb_smb2_request_total{operation="getinfo"} 653
smb_smb2_request_total{operation="ioctl"} 103
smb_smb2_request_total{operation="keepalive"} 1
smb_smb2_request_total{operation="lock"} 0
smb_smb2_request_total{operation="logoff"} 0
smb_smb2_request_total{operation="negprot"} 1
smb_smb2_request_total{operation="notify"} 0
smb_smb2_request_total{operation="read"} 228
smb_smb2_request_total{operation="sesssetup"} 2
smb_smb2_request_total{operation="setinfo"} 109
smb_smb2_request_total{operation="tcon"} 2
smb_smb2_request_total{operation="tdis"} 0
smb_smb2_request_total{operation="write"} 145
# HELP smb_smb2_request_inbytes Bytes received for SMB2 requests
# TYPE smb_smb2_request_inbytes gauge
smb_smb2_request_inbytes{operation="break"} 0
smb_smb2_request_inbytes{operation="cancel"} 0
smb_smb2_request_inbytes{operation="close"} 118536
smb_smb2_request_inbytes{operation="create"} 716288
smb_smb2_request_inbytes{operation="find"} 38612
smb_smb2_request_inbytes{operation="flush"} 1320
smb_smb2_request_inbytes{operation="getinfo"} 67916
smb_smb2_request_inbytes{operation="ioctl"} 15313
smb_smb2_request_inbytes{operation="keepalive"} 68
smb_smb2_request_inbytes{operation="lock"} 0
smb_smb2_request_inbytes{operation="logoff"} 0
smb_smb2_request_inbytes{operation="negprot"} 240
smb_smb2_request_inbytes{operation="notify"} 0
smb_smb2_request_inbytes{operation="read"} 25764
smb_smb2_request_inbytes{operation="sesssetup"} 430
smb_smb2_request_inbytes{operation="setinfo"} 15528
smb_smb2_request_inbytes{operation="tcon"} 240
smb_smb2_request_inbytes{operation="tdis"} 0
smb_smb2_request_inbytes{operation="write"} 8.272958e+06
# HELP smb_smb2_request_outbytes Bytes replied for SMB2 requests
# TYPE smb_smb2_request_outbytes gauge
smb_smb2_request_outbytes{operation="break"} 0
smb_smb2_request_outbytes{operation="cancel"} 0
smb_smb2_request_outbytes{operation="close"} 170072
smb_smb2_request_outbytes{operation="create"} 490469
smb_smb2_request_outbytes{operation="find"} 91957
smb_smb2_request_outbytes{operation="flush"} 1020
smb_smb2_request_outbytes{operation="getinfo"} 142258
smb_smb2_request_outbytes{operation="ioctl"} 13849
smb_smb2_request_outbytes{operation="keepalive"} 68
smb_smb2_request_outbytes{operation="lock"} 0
smb_smb2_request_outbytes{operation="logoff"} 0
smb_smb2_request_outbytes{operation="negprot"} 268
smb_smb2_request_outbytes{operation="notify"} 0
smb_smb2_request_outbytes{operation="read"} 1.6145173e+07
smb_smb2_request_outbytes{operation="sesssetup"} 264
smb_smb2_request_outbytes{operation="setinfo"} 7668
smb_smb2_request_outbytes{operation="tcon"} 160
smb_smb2_request_outbytes{operation="tdis"} 0
smb_smb2_request_outbytes{operation="write"} 11600
# HELP smb_smb2_request_duration_microseconds_sum Execution time in microseconds of SMB2 requests
# TYPE smb_smb2_request_duration_microseconds_sum gauge
smb_smb2_request_duration_microseconds_sum{operation="break"} 0
smb_smb2_request_duration_microseconds_sum{operation="cancel"} 0
smb_smb2_request_duration_microseconds_sum{operation="close"} 431570
smb_smb2_request_duration_microseconds_sum{operation="create"} 7.244576e+06
smb_smb2_request_duration_microseconds_sum{operation="find"} 310193
smb_smb2_request_duration_microseconds_sum{operation="flush"} 149128
smb_smb2_request_duration_microseconds_sum{operation="getinfo"} 59480
smb_smb2_request_duration_microseconds_sum{operation="ioctl"} 14357
smb_smb2_request_duration_microseconds_sum{operation="keepalive"} 4
smb_smb2_request_duration_microseconds_sum{operation="lock"} 0
smb_smb2_request_duration_microseconds_sum{operation="logoff"} 0
smb_smb2_request_duration_microseconds_sum{operation="negprot"} 3.737457e+06
smb_smb2_request_duration_microseconds_sum{operation="notify"} 0
smb_smb2_request_duration_microseconds_sum{operation="read"} 30674
smb_smb2_request_duration_microseconds_sum{operation="sesssetup"} 16994
smb_smb2_request_duration_microseconds_sum{operation="setinfo"} 327027
smb_smb2_request_duration_microseconds_sum{operation="tcon"} 192715
smb_smb2_request_duration_microseconds_sum{operation="tdis"} 0
smb_smb2_request_duration_microseconds_sum{operation="write"} 73739
# HELP smb_vfs_io_bytes Number of bytes transferred via underlying VFS I/O layer
# TYPE smb_vfs_io_bytes gauge
smb_vfs_io_bytes{operation="asys_fsync"} 0
smb_vfs_io_bytes{operation="asys_pread"} 1.6126933e+07
smb_vfs_io_bytes{operation="asys_pwrite"} 8.256718e+06
smb_vfs_io_bytes{operation="pread"} 0
smb_vfs_io_bytes{operation="pwrite"} 0
# HELP smb_vfs_io_duration_microseconds_sum Execution time in microseconds of VFS I/O requests
# TYPE smb_vfs_io_duration_microseconds_sum gauge
smb_vfs_io_duration_microseconds_sum{operation="asys_fsync"} 148760
smb_vfs_io_duration_microseconds_sum{operation="asys_pread"} 26798
smb_vfs_io_duration_microseconds_sum{operation="asys_pwrite"} 68615
smb_vfs_io_duration_microseconds_sum{operation="pread"} 0
smb_vfs_io_duration_microseconds_sum{operation="pwrite"} 0
# HELP smb_vfs_io_total Total number of I/O calls to underlying VFS layer
# TYPE smb_vfs_io_total gauge
smb_vfs_io_total{operation="asys_fsync"} 15
smb_vfs_io_total{operation="asys_pread"} 228
smb_vfs_io_total{operation="asys_pwrite"} 145
smb_vfs_io_total{operation="pread"} 0
smb_vfs_io_total{operation="pwrite"} 0
# HELP smb_vfs_total Total number of calls to underlying VFS layer
# TYPE smb_vfs_total gauge
smb_vfs_total{operation="chdir"} 810
smb_vfs_total{operation="chmod"} 0
smb_vfs_total{operation="close"} 23138
smb_vfs_total{operation="closedir"} 2403
smb_vfs_total{operation="createfile"} 0
smb_vfs_total{operation="fallocate"} 0
smb_vfs_total{operation="fchmod"} 79
smb_vfs_total{operation="fchown"} 0
smb_vfs_total{operation="fdopendir"} 2403
smb_vfs_total{operation="fntimes"} 107
smb_vfs_total{operation="fstat"} 18204
smb_vfs_total{operation="fstatat"} 0
smb_vfs_total{operation="ftruncate"} 12
smb_vfs_total{operation="getwd"} 2
smb_vfs_total{operation="lchown"} 0
smb_vfs_total{operation="linkat"} 16
smb_vfs_total{operation="lseek"} 0
smb_vfs_total{operation="lstat"} 29
smb_vfs_total{operation="mkdirat"} 58
smb_vfs_total{operation="mknodat"} 0
smb_vfs_total{operation="open"} 0
smb_vfs_total{operation="openat"} 27689
smb_vfs_total{operation="opendir"} 0
smb_vfs_total{operation="readdir"} 24129
smb_vfs_total{operation="readlinkat"} 0
smb_vfs_total{operation="realpath"} 4
smb_vfs_total{operation="renameat"} 76
smb_vfs_total{operation="rewinddir"} 0
smb_vfs_total{operation="stat"} 2049
smb_vfs_total{operation="symlinkat"} 0
smb_vfs_total{operation="unlinkat"} 110
# HELP smb_vfs_duration_microseconds_sum Execution time in microseconds of VFS requests
# TYPE smb_vfs_duration_microseconds_sum gauge
smb_vfs_duration_microseconds_sum{operation="chdir"} 12524
smb_vfs_duration_microseconds_sum{operation="chmod"} 0
smb_vfs_duration_microseconds_sum{operation="close"} 55989
smb_vfs_duration_microseconds_sum{operation="closedir"} 63982
smb_vfs_duration_microseconds_sum{operation="createfile"} 0
smb_vfs_duration_microseconds_sum{operation="fallocate"} 0
smb_vfs_duration_microseconds_sum{operation="fchmod"} 2447
smb_vfs_duration_microseconds_sum{operation="fchown"} 0
smb_vfs_duration_microseconds_sum{operation="fdopendir"} 39712
smb_vfs_duration_microseconds_sum{operation="fntimes"} 16961
smb_vfs_duration_microseconds_sum{operation="fstat"} 157620
smb_vfs_duration_microseconds_sum{operation="fstatat"} 0
smb_vfs_duration_microseconds_sum{operation="ftruncate"} 180999
smb_vfs_duration_microseconds_sum{operation="getwd"} 6
smb_vfs_duration_microseconds_sum{operation="lchown"} 0
smb_vfs_duration_microseconds_sum{operation="linkat"} 11593
smb_vfs_duration_microseconds_sum{operation="lseek"} 0
smb_vfs_duration_microseconds_sum{operation="lstat"} 828
smb_vfs_duration_microseconds_sum{operation="mkdirat"} 40445
smb_vfs_duration_microseconds_sum{operation="mknodat"} 0
smb_vfs_duration_microseconds_sum{operation="open"} 0
smb_vfs_duration_microseconds_sum{operation="openat"} 790532
smb_vfs_duration_microseconds_sum{operation="opendir"} 0
smb_vfs_duration_microseconds_sum{operation="readdir"} 5.27503e+06
smb_vfs_duration_microseconds_sum{operation="readlinkat"} 0
smb_vfs_duration_microseconds_sum{operation="realpath"} 20
smb_vfs_duration_microseconds_sum{operation="renameat"} 253650
smb_vfs_duration_microseconds_sum{operation="rewinddir"} 0
smb_vfs_duration_microseconds_sum{operation="stat"} 131346
smb_vfs_duration_microseconds_sum{operation="symlinkat"} 0
smb_vfs_duration_microseconds_sum{operation="unlinkat"} 106974
```

When running with profile-per-share enabled ("smbd profiling share = on")
we get additional per-share metrics:

```console
smb_smb2_request_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="close",share="smbshare"} 7351
smb_smb2_request_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="create",share="smbshare"} 59801
smb_smb2_request_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="find",share="smbshare"} 5972
smb_smb2_request_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="flush",share="smbshare"} 10257
...
smb_smb2_request_inbytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="close",share="smbshare"} 1584
smb_smb2_request_inbytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="create",share="smbshare"} 3336
smb_smb2_request_inbytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="find",share="smbshare"} 392
smb_smb2_request_inbytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="flush",share="smbshare"} 88
...
smb_smb2_request_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="close",share="smbshare"} 18
smb_smb2_request_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="create",share="smbshare"} 18
smb_smb2_request_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="find",share="smbshare"} 4
smb_smb2_request_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="flush",share="smbshare"} 1
...
smb_vfs_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="chdir",share="smbshare"} 562
smb_vfs_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="chmod",share="smbshare"} 0
smb_vfs_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="close",share="smbshare"} 4279
smb_vfs_duration_microseconds_sum{client="192.168.122.25",netbiosname="smb-cephfs",operation="closedir",share="smbshare"} 46
...
smb_vfs_io_bytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="asys_fsync",share="smbshare"} 0
smb_vfs_io_bytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="asys_pread",share="smbshare"} 303
smb_vfs_io_bytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="asys_pwrite",share="smbshare"} 17
smb_vfs_io_bytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="pread",share="smbshare"} 0
smb_vfs_io_bytes{client="192.168.122.25",netbiosname="smb-cephfs",operation="pwrite",share="smbshare"} 0
...
smb_vfs_io_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="asys_fsync",share="smbshare"} 1
smb_vfs_io_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="asys_pread",share="smbshare"} 2
smb_vfs_io_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="asys_pwrite",share="smbshare"} 1
smb_vfs_io_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="pread",share="smbshare"} 0
smb_vfs_io_total{client="192.168.122.25",netbiosname="smb-cephfs",operation="pwrite",share="smbshare"} 0
...
```
