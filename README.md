# Game Integration API - Go Microservice Demo

## How to run

## About
### Architecture
### Technologies

- Web Server Framework: [echo]()
- ORM: [bun]()
- Database: `PostgreSQL`
- OTL: `Promethaus`

## Challenges

### Running wallet-client on Linux (Fedora) amd64

0- install dependencies (once)

```sh
sudo dnf install qemu-user qemu-user-binfmt
```

1- setting up the emulators

```sh
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

2- run the container with requested arch
```sh
docker run --platform=linux/arm64 --rm -it -p 8000:8000 "kentechsp/wallet-client"                                                                                                       17s

```

