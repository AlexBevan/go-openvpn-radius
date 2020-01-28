### Go-OpenVPN-Radius

As basic go app that can be used by applications such as OpenVPN to forward radius requests, add the following your openvpn server.conf:

```
script-security 2
tmp-dir /dev/shm
auth-user-pass-verify "go-openvpn-radius -secret '{{openvpn_radius_secret}}' -server {{openvpn_radius_address}} -file" via-file
```

To build a binary for linux:
```
to build for Linux : env GOOS=linux GOARCH=amd64 go build
```