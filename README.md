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
visible only when active SMB connections exists. Execute the folowing `curl`
command on the same machine where you run `smbmetrics` instance:

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
| `smb_openfiles_total`     | Number of currently open files                   |
| `smb_openfiles_access_rw` | Open files with `"RW"` access-mask set           |
| `smb_share_activity`      | Number of remote machines using each share       |
| `smb_share_byremote`      | Number of shares used by each remote machine     |



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
