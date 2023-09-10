package ettt

import (
	"github.com/google/uuid"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

/*
Engine
実行エンジン
*/
type Engine struct {
	GlobalContext
}

/*
New
実行エンジン生成
*/
func New(scenarios []Scenario,
	extensions []ExtensionContext,
	options Options) (Engine, error) {

	// Profileの解析＆変数保持
	profile, err := ParseProfile(resolveProfile(options))
	if err != nil {
		slog.Error("ettt: profile parse error occurred...", err)
		return Engine{}, err
	}

	// 拡張機能コンテキストのMap作成
	extensionMap := make(map[string]ExtensionContext, len(extensions))
	for _, e := range extensions {
		extensionMap[e.ExtensionKey()] = e
	}

	// 実行シナリオリストの作成
	var executeScenarios []ExecuteScenario
	for _, s := range scenarios {
		rv := reflect.ValueOf(s)
		sc := ScenarioContext{
			scenarioName: rv.Type().Name(),
		}
		executeScenarios = append(executeScenarios, ExecuteScenario{
			Scenario:        &s,
			ScenarioContext: &sc,
		})
	}

	// 全体コンテキストの作成
	globalContext := GlobalContext{
		options:    options,
		extensions: extensionMap,
		profile:    profile,
		scenarios:  executeScenarios,
	}

	return Engine{
		globalContext,
	}, nil
}

/*
Run
ツール実行
*/
func (engine *Engine) Run() error {
	var err error
	// 実行開始タイムスタンプの保持（for Report）
	engine.start = time.Now()

	// 結果ルートディレクトリ・結果ディレクトリの作成
	resultRootDir, err := engine.createResultRootDir()
	if err != nil {
		slog.Error("ettt: failure create result root dir.")
		return err
	}
	executionResultDir, err := engine.createDir(resultRootDir, engine.start.Format("20060102_150405"))
	if err != nil {
		slog.Error("ettt: failure create result dir.")
		return err
	}

	// 指定されたシナリオを随時実行
	// IDEA: 並列化対応するのであれば、このあたりから変更
	for i, v := range engine.scenarios {
		slog.Info("ettt: start scenario.", "index", i, "name", v.ScenarioContext.scenarioName)
		v.executionResultDir = executionResultDir
		engine.runScenario(v)
		slog.Info("ettt: end scenario.",
			"index", i,
			"name", v.ScenarioContext.scenarioName,
			"status", v.scenarioResultStatus)
	}

	// 実行終了タイムスタンプの保持（for Report）
	engine.end = time.Now()
	return nil
}

/*
createResultDir
実行結果のルートディレクトリの作成
*/
func (engine *Engine) createResultRootDir() (string, error) {
	// 絶対パスの作成
	var dirPath = engine.options.ResultPath
	if filepath.IsAbs(engine.options.ResultPath) {
		path, err := filepath.Abs("./")
		if err != nil {
			slog.Error("ettt: failure absolute file path.")
			return "", err
		}
		dirPath = path + engine.options.ResultPath
	}

	if f, err := os.Stat(dirPath); os.IsNotExist(err) || !f.IsDir() {
		// 存在しないため、新規作成を行う
		var fileInfo, err = os.Lstat("./")
		if err != nil {
			slog.Error("ettt: failure get file info.")
			return "", err
		}
		fileMode := fileInfo.Mode()
		unixPerms := fileMode & os.ModePerm

		err = os.Mkdir(dirPath, unixPerms)
		if err != nil {
			slog.Error("ettt: failure get file info.")
			return "", err
		}
	} else {
		slog.Info("ettt: already exists result dir.", "resultDir", dirPath)
	}
	return dirPath, nil
}

/*
createDir
与えられた親ディレクトリと子ディレクトリの名称を利用してディレクトリを作成.
*/
func (engine *Engine) createDir(parent string, target string) (string, error) {
	var targetPath = parent + string(os.PathSeparator) + target
	if f, err := os.Stat(targetPath); os.IsNotExist(err) || !f.IsDir() {
		// 存在しないため、新規作成を行う
		var fileInfo, err = os.Lstat("./")
		if err != nil {
			slog.Error("ettt: failure get file info.")
			return "", err
		}
		fileMode := fileInfo.Mode()
		unixPerms := fileMode & os.ModePerm

		err = os.Mkdir(targetPath, unixPerms)
		if err != nil {
			slog.Error("ettt: failure get file info.")
			return "", err
		}
	}
	return targetPath, nil
}

/*
RunScenario
シナリオ単位の実行関数.
*/
func (engine *Engine) runScenario(es ExecuteScenario) {

	var err error
	var scenario = *es.Scenario
	es.id = uuid.New()
	es.start = time.Now()

	// シナリオの結果ディレクトリ作成
	scenarioResultDir, err := engine.createDir(es.executionResultDir, es.id.String())
	if err != nil {
		slog.Error("ettt: failure create result dir.")
		es.end = time.Now()
		es.error = err
		es.scenarioResultStatus = ScenarioFailure
		return
	}
	es.scenarioResultDir = scenarioResultDir

	// Execute Scenario
	slog.Info("ettt: start Setup.")
	es.ScenarioContext.phase = ScenarioPhaseSetup
	err = scenario.Setup(engine.GlobalContext, es.ScenarioContext)
	if err != nil {
		slog.Info("ettt: error Setup.")
		es.end = time.Now()
		es.scenarioResultStatus = ScenarioFailure
		es.error = err
		return
	}
	slog.Info("ettt: end Setup.")

	slog.Info("ettt: start Exercise.")
	es.ScenarioContext.phase = ScenarioPhaseExercise
	err = scenario.Exercise(engine.GlobalContext, es.ScenarioContext)
	if err != nil {
		slog.Warn("ettt: error Exercise.")
		es.end = time.Now()
		es.scenarioResultStatus = ScenarioFailure
		es.error = err
		return
	}
	slog.Info("ettt: end Exercise.")

	slog.Info("ettt: start Verify.")
	es.ScenarioContext.phase = ScenarioPhaseVerify
	err = scenario.Verify(engine.GlobalContext, es.ScenarioContext)
	if err != nil {
		slog.Warn("ettt: error Verify.")
		es.end = time.Now()
		es.scenarioResultStatus = ScenarioFailure
		es.error = err
		return
	}
	slog.Info("ettt: end Verify.")

	slog.Info("ettt: start TearDown.")
	es.ScenarioContext.phase = ScenarioPhaseTearDown
	err = scenario.TearDown(engine.GlobalContext, es.ScenarioContext)
	if err != nil {
		slog.Info("ettt: error TearDown.")
		es.end = time.Now()
		es.scenarioResultStatus = ScenarioFailure
		es.error = err
		return
	}
	slog.Info("ettt: end TearDown.")

	es.scenarioResultStatus = JudgeScenarioResult(*es.ScenarioContext)
	es.end = time.Now()
	es.durationSeconds = es.end.Sub(es.start).Seconds()
}

/*
resolveProfile
Profile設定ファイルのパスを解決する.
オプションで指定がない場合は、デフォルトのProfile設定ファイルのパスを返却する.
*/
func resolveProfile(options Options) string {
	profile := ProfileDefault
	profilePath := ProfilePathDefault
	if "" != options.Profile {
		slog.Debug("ettt: use designation　profile", "profile", profile)
		profile = options.Profile
	}
	if "" != options.ProfilePath {
		slog.Debug("ettt: use designation　profilePath", "profilePath", profilePath)
		profilePath = options.ProfilePath
	}
	var profileFullPath = profilePath + profile + ".yaml"
	slog.Info("ettt: resolved profilePath finally", "profilePath", profileFullPath)
	return profileFullPath
}

/*
JudgeScenarioResult
シナリオ実行結果から、シナリオの実行結果コードを判定する.
アサーションエラーが１件でも含まれている場合は、アサーションエラーとする.
*/
func JudgeScenarioResult(sc ScenarioContext) ScenarioResultStatus {
	for _, v := range sc.setUpPhaseResults {
		if v.Result == CommandAssertionError {
			slog.Info("ettt: execute scenario assertion error on Setup.")
			return ScenarioAssertionError
		}
	}
	for _, v := range sc.exercisePhaseResults {
		if v.Result == CommandAssertionError {
			slog.Info("ettt: execute scenario assertion error on Exercise.")
			return ScenarioAssertionError
		}
	}
	for _, v := range sc.verifyPhaseResults {
		if v.Result == CommandAssertionError {
			slog.Info("ettt: execute scenario assertion error on Verify.")
			return ScenarioAssertionError
		}
	}
	for _, v := range sc.tearDownPhaseResults {
		if v.Result == CommandAssertionError {
			slog.Info("ettt: execute scenario assertion error on TearDown.")
			return ScenarioAssertionError
		}
	}
	slog.Info("ettt: execute scenario successful.")
	return ScenarioSuccess
}
