package compute

import "github.com/nextlag/in_memory_db/internal/errs"

func (a *Compute) Analyze(tokens []string) (Query, error) {
	if len(tokens) == 0 {
		return Query{}, errs.ErrInvalidQuery
	}

	command := tokens[0]
	commandID := commandNameToCommandID(command)

	if commandID == UnknowCommandID {
		return Query{}, errs.ErrInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])
	argumentsNum := commandArgumentsNumber(commandID)

	if len(query.Arguments()) != argumentsNum {
		return Query{}, errs.ErrInvalidQueryArguments
	}

	return query, nil
}
