# simple_home

# build command
- X86_64: 
- LINUX ALPINE: go build --ldflags '-w -linkmode external -extldflags "-static"' -o gcsh-alpine main.go