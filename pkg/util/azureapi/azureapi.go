package azureapi

// BlankARMTemplate is used to validate template parameters
const BlankARMTemplate = `
{
	"mode": "Incremental",
	"template": {
		"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"resources": []
	}
}`
