# Proto.Actor Go — Keecon 유지보수 버전

[English](README.md) | [한국어](README.ko.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/keecon/protoactor-go.svg)](https://pkg.go.dev/github.com/keecon/protoactor-go)
[![checks](https://github.com/keecon/protoactor-go/actions/workflows/checks.yml/badge.svg?branch=main)](https://github.com/keecon/protoactor-go/actions/workflows/checks.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/keecon/protoactor-go)](https://goreportcard.com/report/github.com/keecon/protoactor-go)

이 저장소는 [asynkron/protoactor-go](https://github.com/asynkron/protoactor-go)에서 파생되어 독립적으로
유지보수되는 버전입니다. 원 저장소의 Git 이력과 Apache-2.0 라이선스를 보존하며, 레거시 호환 브랜치와 최신
`main` 개발 라인을 별도로 관리합니다. 이 저장소는 Proto.Actor 원 저장소가 아닙니다.

이 모듈은 Go용 액터 런타임과 라우팅, 스케줄링, 스트림, 영속성, 원격 통신, 클러스터링 패키지를 제공합니다.
직렬화와 프로세스 간 프로토콜은 명시적으로 관리하며, 해당 영역에서는 Protocol Buffers와 gRPC를 사용합니다.

## 지원 및 릴리스 정책

| 버전 라인 | 브랜치 | 상태 | 유지보수 정책 |
| --- | --- | --- | --- |
| v0.5 이상 | `main` | 최신 개발 | 현재 Go 런타임, 의존성, 보안 수정 및 upstream 통합 |
| v0.3 | `v0.3.x` | 레거시 호환 | 명시적으로 요청된 경우에만 백포트 |
| v0.2 | `v0.2.x` | 레거시 호환 | 명시적으로 요청된 경우에만 백포트 |
| 릴리스 없음 | `origin/dev` | upstream 수신 | 내부 스테이징 전용이며 릴리스 생성 금지 |

`main`만 상시 지원합니다. 새 릴리스 태그는 `v0.5.0`부터 사용합니다. `v0.5.0` 태그가 만들어지기 전까지
`main`은 개발 라인이므로, 사용자는 브랜치 최신 상태를 계속 추적하지 말고 검토한 커밋을 고정해야 합니다.
릴리스 라인의 변경 내역은 [CHANGELOG.md](CHANGELOG.md)에 기록합니다.

`v0.2.x`와 `v0.3.x`는 아직 마이그레이션할 수 없는 프로젝트를 위해 유지합니다. 이 브랜치에는 런타임,
의존성, 기능 또는 보안 업데이트를 정기적으로 적용하지 않습니다. 백포트를 요청받으면 대상 브랜치 기준으로
코드를 검토하고, 필요하면 수정해 적용한 뒤 해당 브랜치에서 검증합니다. 호환성은 `main`이 아니라 대상 브랜치의
API, 의존성 및 Go 툴체인을 기준으로 판단합니다.

## upstream 통합

upstream 변경은 다음 경로로 통합합니다.

```text
asynkron/protoactor-go dev
          ↓
      origin/dev
          ↓
단계별 통합 및 충돌 검토
          ↓
         main
```

`origin/dev`는 `upstream/dev`를 처음 받는 브랜치입니다. 충돌, 로컬 동작, 테스트, 보안 영향 및 모듈 경로
차이를 검토한 뒤에만 `main`으로 통합합니다. 로컬 변경을 upstream으로 자동 반영하지 않습니다.

## 설치

모듈 경로는 `github.com/keecon/protoactor-go`입니다.

최신 개발 빌드를 사용하려면 검토를 마친 `main` 커밋을 선택해 고정합니다.

```bash
go get github.com/keecon/protoactor-go@<commit-sha>
```

레거시 라인은 필요한 릴리스 태그 또는 명시적으로 검토한 브랜치 리비전을 사용합니다. 예:

```bash
go get github.com/keecon/protoactor-go@v0.3.0
```

애플리케이션의 액터 생명주기, 메시지 계약, 영속성 및 사용하는 원격·클러스터 프로토콜을 검증하지 않은 채
레거시 버전을 `main`으로 교체하지 마십시오.

## Go 버전

최신 라인은 Go의 최신 major 릴리스 두 개를 지원하며, 그중 하위 릴리스를 언어·툴체인 기준선으로 사용합니다.
현재 기준선은 Go 1.26이고, [checks workflow](.github/workflows/checks.yml)는 Go 1.26과 Go 1.27의 최신 패치
릴리스를 검증합니다. 새 Go major 릴리스로 공식 지원 범위가 바뀌면 `main`의 기준선도 새로운 하위 릴리스로
올립니다. 루트 모듈과 독립된 예제 모듈에는 같은 기준선을 적용합니다.

`go.mod`의 `go` 지시문은 지원하는 최소 툴체인을 나타냅니다. Go 보안 및 중요 버그 수정을 받으려면 지원
릴리스 라인의 최신 패치 버전을 사용해야 합니다.

레거시 브랜치는 과거 Go 및 의존성 기준선을 유지합니다. 백포트 검증 버전은 각 브랜치의 workflow 파일을
기준으로 합니다.

## Hello actor

```go
package main

import (
	"fmt"

	"github.com/keecon/protoactor-go/actor"
)

type hello struct{ who string }

type helloActor struct{ done chan<- struct{} }

func (a *helloActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *hello:
		fmt.Printf("Hello %s\n", msg.who)
		a.done <- struct{}{}
	}
}

func main() {
	system := actor.NewActorSystem()
	done := make(chan struct{})
	props := actor.PropsFromProducer(func() actor.Actor {
		return &helloActor{done: done}
	})

	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &hello{who: "Proto.Actor"})
	<-done
	system.Root.Stop(pid)
}
```

실행 가능한 예제는 [`examples`](examples)에서 확인할 수 있습니다.

## 빌드 및 테스트

루트 모듈을 빌드합니다.

```bash
go build ./...
```

짧은 테스트와 정적 검사를 실행합니다.

```bash
make test-short
make vet
make lint
make vuln
```

일부 클러스터, 스케줄러, 원격 통신 및 영속성 통합 테스트에는 로컬 Consul을 포함해
[`docker-compose.yml`](docker-compose.yml)에 선언된 서비스가 필요합니다. 전체 테스트 target을 실행하기 전에
필요한 서비스를 시작합니다.

```bash
docker compose up -d
make test2
```

예제 디렉터리는 각각 독립된 Go 모듈이며 CI workflow에서 별도로 빌드합니다.

## 기여

개발 변경은 `main`을 대상으로 합니다. 명시적인 백포트 요청이 없다면 레거시 브랜치를 변경하지 않습니다.
브랜치, 테스트, 의존성 및 upstream 동기화 규칙은 [CONTRIBUTING.md](CONTRIBUTING.md)를 확인하십시오.

## 라이선스 및 저작자 표시

이 프로젝트는 [Apache License 2.0](LICENSE)으로 배포합니다. `asynkron/protoactor-go`의 원래 Git 이력,
저작권 고지 및 라이선스 표시는 유지됩니다. 이 저장소에 새로 기여한 코드에도 같은 라이선스를 적용합니다.
