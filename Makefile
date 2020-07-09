handler:
	oapi-codegen -package=monitorAPI -generate="types,chi-server,spec" api/openapi/api.yaml > pkg/api/openapi/monitor.gen.go