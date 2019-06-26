## 数据库
```bash
docker run --name wx-pg -e POSTGRES_PASSWORD=password -e POSTGRES_USER=wx -e POSTGRES_DB=wx -p 0.0.0.0:5438:5432 -d postgres
```

## 后端
``` bash
go get github.com/lib/pq
go build main.go
```