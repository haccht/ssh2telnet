# ssh2telnet
Proxy ssh connection into telnet.

## Options

```
$ ssh2telnet -h
Usage:
  ssh2telnet [OPTIONS]

Application Options:
  -a, --addr=            Address to listen (default: localhost:2222)
  -k, --key=             Path to the host key
  -l, --login            Enable auto login
      --login-prompt=    Login prompt (default: "login: ")
      --password-prompt= Password prompt (default: "Password: ")

Help Options:
  -h, --help  Show this help message
```
## Basic Usage

Start a ssh server.

```
$ ssh2telnet -a :2222
Starting ssh server on :2222
```

Connect the server from another terminal.
The specified username will be interpeted into the hostname.
Now the proxied telnet session is attached to the target host.

```
$ ssh localhost -p 2222 -l 192.168.1.1


RP/0/RSP0/CPU0:R1#show clock
Thu Nov 11 00:00:00.000 JST
00:00:00:000 JST Thu Nov 11 2021

RP/0/RSP0/CPU0:R1#exit
Connection to localhost closed.
```

## Auto Login

ssh2telnet also comes with the auto login feature.
Start a ssh server with the `--login` option and specify `--login-prompt` and/or `--password-prompt` if necessary.

```
$ ssh2telnet -a :2222 -l --login-prompt 'Username: '
Starting ssh server on :2222
```

Connect the server from another terminal.
The specified username will be interpeted into 'username@hostname'.
Login password for the target host is also prompted afterward.

```
$ ssh localhost -p 2222 -l vagrant@192.168.1.1
vagrant@192.168.1.1@localhost's password:


RP/0/RSP0/CPU0:R1#show clock
Thu Nov 11 00:00:00.000 JST
00:00:00:000 JST Thu Nov 11 2021

RP/0/RSP0/CPU0:R1#exit
Connection to localhost closed.
```
