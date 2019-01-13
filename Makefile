

default:
	go test -race -v ./application ./domain ./servers
	go build -v main.go


