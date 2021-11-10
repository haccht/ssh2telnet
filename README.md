# ssh2telnet
Proxy ssh connection into telnet.

## Usage

```
$ ssh2telnet -h
Usage:
  ssh2telnet [OPTIONS]

Application Options:
  -a, --addr= Address to listen on (default: :2222)
  -k, --key=  Path to the host key

Help Options:
  -h, --help  Show this help message
```

Start a ssh server.

```
$ ssh2telnet -a :2222
Starting ssh server on :2222
```

Connect the ssh server from the other terminal.
The given username will be interpeted into the hostname you want to access.

```
$ ssh localhost -p 2222 -l 192.168.1.1


User Access Verification

Username: vagrant
Password:


RP/0/RSP0/CPU0:R1#exit
Connection to localhost closed.
```
