package handler

// Path names
const (
	PING   string = "PING"
	ECHO   string = "ECHO"
	SET    string = "SET"
	GET    string = "GET"
	EXISTS string = "EXISTS"
	DEL    string = "DEL"
	INCR   string = "INCR"
	DECR   string = "DECR"
	LRANGE string = "LRANGE"
	LPUSH  string = "LPUSH"
	RPUSH  string = "RPUSH"
)

var WRITE_COMMANDS = []string{
	SET, DEL, INCR, DECR, LPUSH, RPUSH,
}
