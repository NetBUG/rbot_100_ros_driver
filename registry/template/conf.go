package template

func DefaultConfig() string {
	return `---
# API Endpoint
api:
  host: localhost
  port: 8800

# Database
db:
  # Available options: sqlite3 | postgres
  driver: sqlite3
  path: development.sqlite3
  # For postgres
  # path: host=%host port=%port dbname=%dbname user=%username password=%passwd %options

`
}
