package copier

var anonymize_postgres_map = map[AnonymizeFn]string{
	ANONYMIZE_STRING: "md5(random()::text)",
	ANONYMIZE_INT:    "",
}

func (a AnonymizeFn) for_postgres() string {
	return anonymize_postgres_map[a]
}
