package models

type Response struct {
	Cmds     []string `pg:",array"`
	Messages []string `pg:",array"`
}
