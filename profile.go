package ettt

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

const (
	// ProfilePathDefault デフォルトのProfile読み込みパス.
	ProfilePathDefault string = "./profiles/"
	ProfileDefault     string = "default"
)

/*
Profile
Profile設定ファイルを保持する構造体
*/
type Profile struct {
	Name      string
	Variables []ProfileVariable
}

/*
ProfileVariable
Profile設定ファイル中のKey=Value保持
*/
type ProfileVariable struct {
	Key   string
	Value string
}

func ParseProfile(target string) (Profile, error) {
	slog.Info("ettt: parse profile start.", "source", target)

	// Read Yaml File
	var bytes, err = os.ReadFile(target)
	if err != nil {
		slog.Error("ettt: read profile failure.", err, "source", target)
		// 空とエラーを返却
		return Profile{}, err
	}

	// Yaml -> Struct
	profileVariables := Profile{}
	err = yaml.Unmarshal(bytes, &profileVariables)
	if err != nil {
		slog.Error("ettt: parse profile failure.", err, "source", target)
		// 空とエラーを返却
		return Profile{}, err
	}
	return profileVariables, nil
}
