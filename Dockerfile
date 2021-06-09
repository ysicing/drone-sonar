FROM hub.51talk.biz/library/goeff:1.15-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

ENV GOPRIVATE="code.51talk.com"

WORKDIR /go/src/

COPY go.mod go.mod

COPY go.sum go.sum

RUN go mod download

COPY . /go/src/

RUN go build

FROM koalaman/shellcheck-alpine:stable as shellcheck

FROM sonarsource/sonar-scanner-cli:4.6

LABEL MAINTAINER=ysicing

COPY --from=builder /go/src/drone-sonar /usr/bin/drone-sonar

COPY --from=shellcheck /bin/shellcheck /usr/bin/shellcheck

COPY hack/docker/entrypoint.sh /entrypoint.sh

ENV TZ=Asia/Shanghai

RUN chmod +x /usr/bin/drone-sonar /usr/bin/shellcheck /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]