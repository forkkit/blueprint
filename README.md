# blueprint

Blueprint takes a template and turns it into a application.


# Building Protobuf 3 Locally

If you're on mac you'll need to `brew install automake` if you haven't already

```
  mkdir tmp
  cd tmp
  git clone https://github.com/google/protobuf
  cd protobuf
  ./autogen.sh
  ./configure --prefix /usr/local/Cellar/protobuf/3.0.0-dev
  make
  make install
  brew switch protobuf 3.0.0-dev
```


You're going to need the plugins installed globally (the binaries in your path)

```
  go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
  go get -u github.com/golang/protobuf/protoc-gen-go
```
