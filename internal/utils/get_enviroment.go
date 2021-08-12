package utils

// GetDBEnvironments gets database variables from .env file.
func GetDBEnvironments(conf ConfigService) (user, pass, host, port, dbname string, err error) {
	user, err = conf.GetEnv("DB_USER")
	if err != nil {
		return "", "", "", "", "", ErrFindVariable
	}
	pass, err = conf.GetEnv("DB_PASSWORD")
	if err != nil {
		return "", "", "", "", "", ErrFindVariable
	}
	host, err = conf.GetEnv("DB_HOST")
	if err != nil {
		return "", "", "", "", "", ErrFindVariable
	}
	port, err = conf.GetEnv("DB_PORT")
	if err != nil {
		return "", "", "", "", "", ErrFindVariable
	}
	dbname, err = conf.GetEnv("DB_NAME")
	if err != nil {
		return "", "", "", "", "", ErrFindVariable
	}
	return user, pass, host, port, dbname, nil
}

// GetTokenTTL gets tokenTTL variable from .env file.
func GetTokenTTL(conf ConfigService) (string, error) {
	tokenTTL, err := conf.GetEnv("TOKEN_TTL")
	if err != nil {
		return "", ErrFindVariable
	}
	return tokenTTL, nil
}

// GetRabbitMQURL gets URL variable from .env file.
func GetRabbitMQURL(conf ConfigService) (string, error) {
	URL, err := conf.GetEnv("RABBITMQ_URL")
	if err != nil {
		return "", ErrFindVariable
	}
	return URL, nil
}
