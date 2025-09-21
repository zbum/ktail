# ktail

Kubernetes 로그를 실시간으로 추적하는 유틸리티로, 퍼지 파인더를 사용한 대화형 네임스페이스 및 파드 선택 기능을 제공합니다.

## 개요

`ktail`은 Kubernetes 파드 로그에 대한 `tail` 기능을 제공하는 도구입니다. 퍼지 파인더를 사용하여 네임스페이스와 파드를 대화형으로 선택할 수 있어 더 나은 사용자 경험을 제공하며, 클러스터의 모든 파드에서 로그를 쉽게 추적할 수 있습니다.

## 주요 기능

- 🎯 **대화형 선택**: 퍼지 파인더를 사용하여 네임스페이스와 파드를 대화형으로 선택
- 📊 **파드 상태 표시**: 파드 상태에 대한 시각적 표시 (Running, Pending, Failed 등)
- 🔄 **실시간 로그 스트리밍**: `tail -f` 동작으로 실시간 로그 추적
- 🎨 **컬러 출력**: 네임스페이스와 파드 이름을 초록색으로 표시하여 가독성 향상
- 🎛️ **유연한 옵션**: 사용자 정의 tail 라인, 컨테이너 선택, 색상 비활성화 등 지원
- 🚀 **사용하기 쉬움**: 합리적인 기본값을 가진 간단한 CLI 인터페이스

## 사전 요구사항

- Go 1.24 이상
- Kubernetes 클러스터 접근 권한
- 클러스터에 접근할 수 있도록 구성된 `kubectl`

**참고**: 외부 의존성이 필요하지 않습니다! 이 도구는 대화형 선택을 위해 `go-fuzzyfinder` 라이브러리를 사용하므로 `fzf`를 별도로 설치할 필요가 없습니다.

## 설치

### 옵션 1: GitHub 릴리즈에서 다운로드 (권장)

1. [릴리즈 페이지](https://github.com/your-username/ktail/releases)로 이동
2. 플랫폼에 맞는 아카이브를 다운로드:
   - **Linux**: `ktail-linux.tar.gz` (amd64, arm64, 386, arm, ppc64, ppc64le, mips, mipsle, mips64, mips64le, riscv64, s390x 포함)
   - **macOS**: `ktail-darwin.tar.gz` (amd64, arm64 포함)
   - **Windows**: `ktail-windows.zip` (amd64, 386, arm64 포함)
   - **모든 플랫폼**: `ktail-all.tar.gz` (모든 바이너리 포함)

3. 아카이브를 추출하고 바이너리를 PATH에 추가:
```bash
# Linux/macOS용
tar -xzf ktail-linux.tar.gz
chmod +x ktail-linux-*
sudo mv ktail-linux-* /usr/local/bin/ktail

# Windows용
# ktail-windows.zip을 추출하고 PATH에 추가
```

4. 설치 확인:
```bash
ktail --help
```

### 옵션 2: 소스에서 빌드

1. 저장소 클론:
```bash
git clone <repository-url>
cd ktail
```

2. 프로젝트 빌드:
```bash
# Makefile 사용 (권장)
make build

# 또는 수동으로
go build -o ktail
```

3. 실행 가능하게 만들고 PATH에 추가 (선택사항):
```bash
chmod +x ktail
sudo mv ktail /usr/local/bin/
```

### 옵션 3: 모든 플랫폼용 빌드

```bash
# 모든 지원 플랫폼과 아키텍처에 대해 빌드
make build-all

# 특정 플랫폼만 빌드
make build-linux
make build-darwin
make build-windows

# 특정 아키텍처만 빌드
make build-arch GOOS=linux GOARCH=arm64
```

## 사용법

### 기본 사용법

```bash
# 대화형 모드 - 네임스페이스와 파드를 선택하여 로그 추적
ktail

# 특정 네임스페이스의 모든 파드 로그 추적
ktail -n my-namespace

# 특정 파드의 로그 추적
ktail -n my-namespace -p my-pod

# 도움말 보기
ktail --help
```

### 명령줄 옵션

```bash
사용법:
  ktail [flags]

플래그:
  -h, --help               ktail 도움말
  -m, --multi              멀티 선택 활성화 (기본값: true)
  -n, --namespace string   Kubernetes 네임스페이스 (제공되지 않으면 대화형으로 선택)
  -p, --pod string         파드 이름 (제공되지 않으면 모든 파드 선택)
  -t, --tail int           로그 끝에서 보여줄 라인 수 (기본값: 100)
  -c, --container string   컨테이너 이름
      --no-color           컬러 출력 비활성화
```

### 예제

```bash
# 대화형으로 네임스페이스와 파드 선택
ktail

# production 네임스페이스의 모든 파드 로그 추적
ktail -n production

# 특정 파드의 로그 추적
ktail -n production -p web-app-7d4f8b9c6-xyz12

# 마지막 500줄을 보여주고 추적
ktail -t 500

# tail 스타일 플래그 사용 (최근 1000줄 추적)
ktail -1000f

# 특정 네임스페이스에서 최근 200줄 추적
ktail -200f -n staging

# 컬러 출력 비활성화
ktail --no-color

# 특정 컨테이너의 로그 추적
ktail -n production -p web-app -c nginx
```

## 개발

### 사전 요구사항
- Go 1.24 이상
- Make (Makefile 사용을 위해)

### 테스트 실행
```bash
# 모든 테스트 실행
make test

# 커버리지와 함께 테스트 실행
make test-coverage

# 또는 수동으로
go test ./...
```

### 다양한 플랫폼용 빌드
```bash
# 모든 지원 플랫폼에 대해 빌드
make build-all

# 특정 플랫폼만 빌드
make build-linux
make build-darwin
make build-windows

# 특정 아키텍처만 빌드
make build-arch GOOS=linux GOARCH=arm64

# 지원하는 모든 플랫폼 목록 보기
make list-platforms
```

### 코드 품질
```bash
# 코드 포맷팅
make fmt

# 린터 실행
make lint

# go vet 실행
make vet

# 보안 검사
make security

# 모든 품질 검사 실행
make dev-setup
```

### 릴리즈 생성
```bash
# 새로운 릴리즈 태그 생성 및 GitHub에 푸시
make release-github

# 또는 수동으로
make tag
make release-tag
```

이렇게 하면:
1. git 태그 생성 (예: v1.0.0)
2. 태그를 GitHub에 푸시
3. GitHub Actions가 모든 플랫폼용 바이너리 빌드 트리거
4. 모든 바이너리가 첨부된 GitHub 릴리즈 생성

## 기여하기

1. 저장소 포크
2. 기능 브랜치 생성
3. 변경사항 적용
4. 해당하는 경우 테스트 추가
5. 풀 리퀘스트 제출

## 라이선스

이 프로젝트는 MIT 라이선스 하에 라이선스가 부여됩니다.

## 지원하는 플랫폼

- **Linux**: amd64, arm64, 386, arm, ppc64, ppc64le, mips, mipsle, mips64, mips64le, riscv64, s390x
- **macOS**: amd64, arm64
- **Windows**: amd64, 386, arm64

## 문제 해결

### 일반적인 문제

1. **kubectl 연결 오류**: `kubectl`이 올바르게 구성되어 있는지 확인
2. **권한 오류**: 클러스터에 대한 적절한 권한이 있는지 확인
3. **파드가 보이지 않음**: 네임스페이스 권한을 확인

### 로그 레벨 조정

환경 변수를 사용하여 로그 레벨을 조정할 수 있습니다:

```bash
export KTAIL_LOG_LEVEL=debug
./ktail
```

## 변경 로그

### v1.0.0
- 초기 릴리즈
- 대화형 네임스페이스 및 파드 선택
- 실시간 로그 스트리밍
- 다중 플랫폼 지원

---

**English**: [README.en.md](README.en.md)