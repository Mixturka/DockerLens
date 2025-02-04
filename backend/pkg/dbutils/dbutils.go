package dbutils

import (
	"net"
	"net/url"
)

// Builds correct postgres connection string based on provided data
func BuildPostgresURL(user, password, host, port, db string) string {
	if port == "" {
		port = "5432"
	}

	hostPort := net.JoinHostPort(host, port)

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   hostPort,
		Path:   "/" + db,
	}

	return u.String()
}
