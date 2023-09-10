package ettt

/*
Scenario
シナリオインタフェース
*/
type Scenario interface {
	/*
		Setup
		テストの事前準備.
	*/
	Setup(gc GlobalContext, sc *ScenarioContext) error
	/*
		Exercise
		テストの実行.
	*/
	Exercise(gc GlobalContext, context *ScenarioContext) error
	/*
		Verify
		テスト結果の確認.
	*/
	Verify(gc GlobalContext, context *ScenarioContext) error
	/*
		TearDown
		テストの後片付け.
	*/
	TearDown(gc GlobalContext, context *ScenarioContext) error
}
