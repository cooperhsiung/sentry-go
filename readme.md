# sentry-go

A handy metrics monitor.

### Install

```
git clone https://github.com/cooperhsiung/sentry-go.git
cd sentry-go
export GO111MODULE=on
export GOPROXY=https://goproxy.io
go mod tidy
```

### Run

```
go run main.go // web service
go run worker.go // timer for monitor
go run collector.go // timer for collector
```

### Directory

dev => \$HOME/go/src/gitlab.com/sentry.go

pord => /opt/srv/sentry-go
