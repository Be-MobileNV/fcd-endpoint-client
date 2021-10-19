Push GPS position example
=========================

This example will push 100 GPS position messages to the given FCD-endpoint server.

Execution instructions
----------------------
Make sure you have installed Golang on your system and put this repo in your gopath.
Build the Golang code (with a version that supports go modules, e.g. go1.13):

```
go build
```

Execute the binary with option -h to see the command-line parameters.

Example usage (change _XXX_ to your assigned provider id):

```
go run test/Golang/pushGPSPositions.go -address XXX-fcd-endpoint.be-mobile.biz -port 443 -username 666 -password 'accompanying password' -tls true
```