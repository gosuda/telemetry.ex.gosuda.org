atlas schema apply \
--url "postgres://postgres:pass@localhost:5432/sqlc?sslmode=disable" \
--dev-url "docker://postgres" \
--to "file://schema.sql"