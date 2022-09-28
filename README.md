# Multifile

> Multifile is designed to be an Nginx-like static web server, but much more easier to use.

## Start

You need an Linux/Windows/MacOS as your server.

I strongly recommend using [Golang](https://go.dev/doc/install) to compile the binaries from source yourself.

### Build From Source

If you download the pre-build binary from the release page, you can skip this step.

```bash
git clone https://github.com/LordCasser/multifile.git
cd multifile
go build .
```

## Init Environment

For windows

```powershell
multifile.exe -init
```

For Linux/MacOS



```bash
./multifile -init
```

and you will find

```bash
├─multifile  #binary
├─resources  #used for storing certificates
└─static     #used for storing static web files
```

## Start Service

### Set certificates

Just put HTTPS certificates for nginx into `resources` folder, like

```
resources/
├── tls.crt
├── tls.csr
├── tls.key
└── tls.pem
```

### Set static web files (example by hexo)

There are already many tutorials on the use of hexo on the Internet, so I won't repeat them here. This user manual starts with the hexo already deployed.



Compile static file (in hexo floder)

```bash
hexo g 
```

Copy files to `static` folder

```bash
cp -r ./publish ../multifile/static/
```



### Start Server

with HTTPS (need HTTPS certificates)

```bash
./multifile -SSL
```

without HTTPS

```bash
./multifile
```

custom port

```bash
./multifile -Port 8080 #for example
```

## And enjoy it
