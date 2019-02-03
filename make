stringer -type URoles ./vstate/auth.go
stringer -type Event ./vstate/event.go
stringer -type State ./vstate/state.go 
go build -race .
