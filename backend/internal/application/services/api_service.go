package services

type ApiServices struct {
	RoutingTable RoutingTableService
	Device       DeviceService
	Environment  EnvironmentService
}
