package ettt

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

/*
ScenarioResultStatus シナリオ実行結果ステータス.
*/
type ScenarioResultStatus string

const (
	ScenarioSuccess        = ScenarioResultStatus("ScenarioSuccess")
	ScenarioFailure        = ScenarioResultStatus("ScenarioFailure")
	ScenarioAssertionError = ScenarioResultStatus("ScenarioAssertionError")
)

type ScenarioPhase string

const (
	ScenarioPhaseSetup    = ScenarioPhase("SetUp")
	ScenarioPhaseExercise = ScenarioPhase("Exercise")
	ScenarioPhaseVerify   = ScenarioPhase("Verify")
	ScenarioPhaseTearDown = ScenarioPhase("TearDown")
)

/*
ExtensionContext 拡張機能用の情報を保持するコンテキスト.
*/
type ExtensionContext interface {
	ExtensionKey() string
}

/*
GlobalContext テスト実行時の全体の情報を保持するコンテキスト.
*/
type GlobalContext struct {
	// オプション
	options Options
	// 拡張機能コンテキスト
	extensions map[string]ExtensionContext
	// プロファイル
	profile Profile
	// 実行シナリオリスト
	scenarios []ExecuteScenario
	// 開始時間
	start time.Time
	// 終了時間
	end time.Time
}

/*
RegistrationExtensionContext
拡張機能コンテキストを登録する.
外部から変更不可としている.
キーはその拡張コンテキストが知っている自身のキー.
TODO : for duplicate key action...
*/
func (gc *GlobalContext) RegistrationExtensionContext(name string, extensionContext ExtensionContext) {
	gc.extensions[name] = extensionContext
}

/*
GetExtensionContext
拡張機能コンテキストを取得.
キーはその拡張コンテキストが知っている自身のキー
TODO : for duplicate key action...
*/
func (gc *GlobalContext) GetExtensionContext(name string) ExtensionContext {
	return gc.extensions[name]
}

/*
Options
実行オプション
*/
type Options struct {
	Profile     string
	ProfilePath string
	// 結果出力パス.
	ResultPath string
	// テンプレートディレクトリパス.()
	TemplateDirPath string
}

func DefaultOptions() Options {
	return Options{}
}

/*
ExecuteScenario
実行シナリオ情報.
単一シナリオ実行時の単位にラップした構造体.
*/
type ExecuteScenario struct {
	*ScenarioContext
	Scenario *Scenario
}

/*
ScenarioContext
テスト実行時のシナリオ毎の情報を保持するコンテキスト.
*/
type ScenarioContext struct {
	// 実行ID
	id uuid.UUID
	// 開始時間
	start time.Time
	// 終了時間
	end time.Time
	// 実行時間（秒）
	durationSeconds float64
	// 実行毎の結果ディレクトリ
	executionResultDir string
	// シナリオの結果ディレクトリ
	scenarioResultDir string
	// エビデンス格納ディレクトリ
	evidencesDir string
	// 詳細レポート格納ディレクトリ
	detailsDir string
	// シナリオ名
	scenarioName string
	// シナリオステータス
	scenarioResultStatus ScenarioResultStatus
	// エラー
	error error
	// Store変数
	Store StoreVariables
	// 現在のPhase
	phase ScenarioPhase
	// SetUpフェーズのCommand実行結果
	setUpPhaseResults []CommandResult
	// ExerciseフェーズのCommand実行結果
	exercisePhaseResults []CommandResult
	// VerifyフェーズのCommand実行結果
	verifyPhaseResults []CommandResult
	// TearDownフェーズのCommand実行結果
	tearDownPhaseResults []CommandResult
}

/*
RegistrationCommandResult
コマンド実行結果を登録
*/
func (sc *ScenarioContext) RegistrationCommandResult(commandResult CommandResult) {
	switch sc.phase {
	case ScenarioPhaseSetup:
		sc.setUpPhaseResults = append(sc.setUpPhaseResults, commandResult)
	case ScenarioPhaseExercise:
		sc.exercisePhaseResults = append(sc.exercisePhaseResults, commandResult)
	case ScenarioPhaseVerify:
		sc.verifyPhaseResults = append(sc.verifyPhaseResults, commandResult)
	case ScenarioPhaseTearDown:
		sc.tearDownPhaseResults = append(sc.tearDownPhaseResults, commandResult)
	default:
		fmt.Errorf("unknown scenario phase %s", sc.phase)
	}
}

/*
CurrentPhase
シナリオの現在Phaseを取得
*/
func (sc ScenarioContext) CurrentPhase() ScenarioPhase {
	return sc.phase
}

/*
StoreVariables
ストア変数.
*/
type StoreVariables struct {
	Variables map[string]string
}
