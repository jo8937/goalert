package main

// generated from [configure_source_generator.sh]
func LoadConfig(env string) DatabaseConfig {
	return DatabaseConfig{"user", "pass", "host", 3306, "dbname"}
}
