# ktail

Kubernetes 파드 로그를 실시간으로 추적하는 간단하고 강력한 도구입니다.

## 주요 기능

- 🎯 **대화형 선택**: 퍼지 파인더를 사용하여 네임스페이스와 파드를 대화형으로 선택
- 🔄 **실시간 로그 스트리밍**: `tail -f` 동작으로 실시간 로그 추적
- 🎨 **컬러 출력**: 네임스페이스와 파드 이름을 초록색으로 표시하여 가독성 향상
- 👀 **Watch 모드**: 네임스페이스에서 새로 생성되는 파드의 로그를 자동으로 추적
- 🎛️ **유연한 옵션**: 사용자 정의 tail 라인, 컨테이너 선택, 색상 비활성화 등 지원

## 사전 요구사항

- Kubernetes 클러스터 접근 권한
- 클러스터에 접근할 수 있도록 구성된 `kubectl`

## 설치

### 소스에서 빌드

```bash
# 저장소 클론
git clone <repository-url>
cd ktail

# 빌드
go build -o ktail

# 실행 가능하게 만들고 PATH에 추가 (선택사항)
chmod +x ktail
sudo mv ktail /usr/local/bin/
```

## 사용법

### 기본 명령어

```bash
# 대화형 모드 - 네임스페이스와 파드를 선택하여 로그 추적
ktail

# 특정 네임스페이스의 모든 파드 로그 추적
ktail -n my-namespace

# 특정 파드의 로그 추적
ktail -n my-namespace -p my-pod

# Watch 모드 - 네임스페이스에서 새로 생성되는 파드의 로그를 자동으로 추적
ktail -n my-namespace -w
```

### 명령줄 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-n, --namespace` | Kubernetes 네임스페이스 | 대화형 선택 |
| `-p, --pod` | 파드 이름 | 모든 파드 |
| `-c, --container` | 컨테이너 이름 | 첫 번째 컨테이너 |
| `-t, --tail` | 로그 끝에서 보여줄 라인 수 | 100 |
| `-m, --multi` | 멀티 선택 활성화 | true |
| `-w, --watch` | Watch 모드 (네임스페이스만 선택 시) | false |
| `--no-color` | 컬러 출력 비활성화 | false |

### 사용 예제

#### 1. 대화형 모드
```bash
# 네임스페이스와 파드를 대화형으로 선택
ktail
```

#### 2. 특정 네임스페이스의 모든 파드
```bash
# production 네임스페이스의 모든 파드 로그 추적
ktail -n production
```

#### 3. 특정 파드
```bash
# 특정 파드의 로그 추적
ktail -n production -p web-app-7d4f8b9c6-xyz12
```

#### 4. Watch 모드
```bash
# 네임스페이스에서 새로 생성되는 파드의 로그를 자동으로 추적
ktail -n production -w
```

#### 5. 컨테이너 지정
```bash
# 특정 컨테이너의 로그 추적
ktail -n production -p web-app -c nginx
```

#### 6. 로그 라인 수 조정
```bash
# 마지막 500줄을 보여주고 추적
ktail -t 500

# tail 스타일 플래그 사용 (최근 1000줄 추적)
ktail -1000f
```

#### 7. 색상 비활성화
```bash
# 컬러 출력 비활성화
ktail --no-color -n production
```

## 문제 해결

### 일반적인 문제

1. **kubectl 연결 오류**: `kubectl`이 올바르게 구성되어 있는지 확인
2. **권한 오류**: 클러스터에 대한 적절한 권한이 있는지 확인
3. **파드가 보이지 않음**: 네임스페이스 권한을 확인

## 라이선스

이 프로젝트는 MIT 라이선스 하에 라이선스가 부여됩니다.

---

**English**: [README.en.md](README.en.md)