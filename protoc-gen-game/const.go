package main

import (
	"github.com/jhump/goprotoc/plugins"
)

const (
	isActionSetExtensionFieldNumber             = 50000
	actionExtensionFieldNumber                  = 50001
	isGameStateExtensionFieldNumber             = 50002
	isActionServiceResponseExtensionFieldNumber = 50003
	isEnumKeyFieldNumber                        = 50004

	protoPackageName       = "game_engine"
	responseFieldName      = "response"
	responseStateFieldName = "state"
	responseErrorFieldName = "error"
	errorTypeName          = "Error"

	stateVariable    = "state"
	inputVariable    = "in"
	responseVariable = "res"
)

var (
	goNames plugins.GoNames

	responseFieldNameCamelCase      = goNames.CamelCase(responseFieldName)
	responseStateFieldNameCamelCase = goNames.CamelCase(responseStateFieldName)
	responseErrorFieldNameCamelCase = goNames.CamelCase(responseErrorFieldName)
)
