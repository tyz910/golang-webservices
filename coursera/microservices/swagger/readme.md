go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger serve ../gateway/session/session.swagger.json -p 8082

rm -rf sess-client && \
    mkdir sess-client && \
    swagger generate client -f ../gateway/session/session.swagger.json -A sess-client/ -t ./sess-client/

mkdir sess-server
swagger generate server -f ../gateway/session/session.swagger.json -A sess-server/ -t ./sess-server/
