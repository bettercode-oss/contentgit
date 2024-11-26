# ContentGit

ContentGit은 컨텐츠 저장소로 Git 처럼 컨텐츠를 관리하며 REST API를 제공합니다.

## 시작하기

1. Go 설치
Go는 1.23.3 버전을 사용합니다.

https://go.dev/doc/install

2. Go 의존성 라이브러리 설치

```bash
go mod download
```

3. Postgres 설치
아래 도커 명령어로 Postgres 컨테이너를 실행합니다.

```bash
docker run -d --name pgmq-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 tembo.docker.scarf.sh/tembo/pg16-pgmq:latest
```

4. 데이터베이스 생성
psql 이나 pgAdmin 등으로 데이터베이스에 접속하여 아래 SQL을 실행합니다.
[create_database.sql](script/database/create_database.sql)

5. 실행
```bash
export DB_PASSWORD=postgres
go run main.go
```

## REST API 명세
아래 테스트 코드를 참고하세요.
[content_controller_test.go](ports/in/web/content_controller_test.go)

## 테스트

```bash
go test ./...
```