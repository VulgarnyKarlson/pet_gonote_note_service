package tests

type TestIntegrationService int

const (
	TestIntegrationServicePostgres TestIntegrationService = iota
	TestIntegrationServiceRabbitMQ
	TestIntegrationServiceRedis
)
