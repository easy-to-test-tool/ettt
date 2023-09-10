package ettt

import "github.com/google/uuid"

/*
CommandResultStatus コマンド結果ステータス.
定数定義と組み合わせてEnum的に扱う.
属性が1つのみなのが微妙か...
*/
type CommandResultStatus string

const (
	// CommandSuccess コマンド正常終了
	CommandSuccess = CommandResultStatus("CommandSuccess")
	// CommandFailure コマンド異常終了
	CommandFailure = CommandResultStatus("CommandFailure")
	// CommandAssertionError コマンドアサーションエラー
	CommandAssertionError = CommandResultStatus("CommandAssertionError")
)

/*
CommandResult コマンド結果
*/
type CommandResult struct {
	Id               uuid.UUID
	Result           CommandResultStatus
	Message          string
	CustomReportPath string
	Error            error
}

/*
Command
各コマンドが実装するインタフェース
*/
type Command interface {
	/*
		GetId
		コマンドID
	*/
	GetId() uuid.UUID
	/*
		Execute
		コマンド実行
	*/
	Execute(gc GlobalContext, sc *ScenarioContext)
}
