package compute

import "github.com/nextlag/in_memory_db/internal/server/entity"

const (
	UnknowCommandID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

var nameToID = map[string]int{
	entity.SetCommand: SetCommandID,
	entity.GetCommand: GetCommandID,
	entity.DelCommand: DelCommandID,
}

func commandNameToCommandID(command string) int {
	status, ok := nameToID[command]
	if !ok {
		return UnknowCommandID
	}

	return status
}

const (
	setCommandArgumentNum = 2
	getCommandArgumentNum = 1
	delCommandArgumentNum = 1
)

var argumentsNum = map[int]int{
	SetCommandID: setCommandArgumentNum,
	GetCommandID: getCommandArgumentNum,
	DelCommandID: delCommandArgumentNum,
}

func commandArgumentsNumber(commandID int) int {
	return argumentsNum[commandID]
}
