package ettt

import (
	"html/template"
	"log"
	"os"
)

func GlobalReport(globalContext GlobalContext) error {
	return nil
}

func ScenarioReport(globalContext GlobalContext, scenarioContext ScenarioContext) error {
	// カスタムテンプレートディレクトリが指定されていなければ、ツールオリジナルを利用する
	var templateDirPath string = ""
	if "" == globalContext.options.TemplateDirPath {
		templateDirPath = DefaultReportTemplateDirPath
	} else {
		templateDirPath = globalContext.options.TemplateDirPath
	}

	t, err := template.ParseGlob(templateDirPath + "/" + DefaultReportTemplateResultPath)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	if err := t.Execute(os.Stdout, scenarioContext); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
	if err != nil {
		log.Fatalf("template error: %v", err)
		return err
	}
	return nil
}
