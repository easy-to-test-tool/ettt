package ettt

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
)

/*
EvidenceType
エビデンスの種別を定義.
*/
type EvidenceType string

const (
	// TextType テキストエビデンス
	TextType = EvidenceType("TextType")
	// BinaryType バイナリエビデンス
	BinaryType = EvidenceType("BinaryType")
)

/*
Evidence
エビデンスを表す構造体
*/
type Evidence struct {
	Id              uuid.UUID
	BelongCommandId uuid.UUID
	EvidenceType    EvidenceType
	Name            string
	Path            string
}

/*
RegistrationTextEvidenceOptions
エビデンス登録を行うためのオプション
*/
type RegistrationTextEvidenceOptions struct {
	Name         string
	TextContents string
	Command      Command
}

/*
RegistrationBinaryEvidenceOptions
エビデンス登録を行うためのオプション
*/
type RegistrationBinaryEvidenceOptions struct {
	Name           string
	BinaryContents []byte
	Command        Command
}

/*
registrationTextEvidence
エビデンス登録受付関数.
テキスト形式のエビデンスを結果ディレクトリへ保存し、ScenarioContextへ登録.
*/
func registrationTextEvidence(sc *ScenarioContext, options RegistrationTextEvidenceOptions) error {
	if options.Name == "" || len(options.Name) == 0 {
		return fmt.Errorf("options.Name parameter is requireds")
	}
	if options.TextContents == "" || len(options.TextContents) == 0 {
		return fmt.Errorf("options.TextContents parameter is required")
	}
	var evidenceId = uuid.New()
	var registrationFilePath = sc.evidencesDir +
		string(filepath.Separator) +
		evidenceId.String() +
		string(filepath.Separator) +
		options.Name
	slog.Debug("registration text evidence file", "path", registrationFilePath)
	data := []byte(options.TextContents)
	return registrationEvidence(
		sc, TextType, evidenceId, registrationFilePath, options.Name, data, options.Command,
	)
}

/*
registrationBinaryEvidence
エビデンス登録受付関数.
バイナリ形式のエビデンスを結果ディレクトリへ保存し、ScenarioContextへ登録.
*/
func registrationBinaryEvidence(sc *ScenarioContext, options RegistrationBinaryEvidenceOptions) error {
	if options.Name == "" || len(options.Name) == 0 {
		return fmt.Errorf("options.Name parameter is requireds")
	}
	if options.BinaryContents == nil || reflect.ValueOf(options.BinaryContents).IsNil() {
		return fmt.Errorf("options.BinaryContents parameter is required")
	}
	var evidenceId = uuid.New()
	var registrationFilePath = sc.evidencesDir +
		string(filepath.Separator) +
		evidenceId.String() +
		string(filepath.Separator) +
		options.Name
	slog.Debug("registration binary evidence file", "path", registrationFilePath)
	return registrationEvidence(
		sc, BinaryType, evidenceId, registrationFilePath, options.Name, options.BinaryContents, options.Command,
	)
}

/*
registrationEvidence
バイナリ形式のエビデンスを結果ディレクトリへ保存し、ScenarioContextへ登録.
登録処理の本処理は、この関数で実装する.
*/
func registrationEvidence(
	sc *ScenarioContext,
	evidenceType EvidenceType,
	evidenceId uuid.UUID,
	registrationFilePath string,
	name string, contents []byte,
	command Command) error {

	// 既に登録済みである場合はエラーとする
	f, err := os.Stat(registrationFilePath)
	if os.IsExist(err) {
		slog.Error("already exists evidence file.", err, "path", registrationFilePath)
		return fmt.Errorf("already exists evidence file")
	}
	if f != nil && f.IsDir() {
		slog.Error("already exists and directory.", err, "path", registrationFilePath)
		return fmt.Errorf("already exists and directory")
	}

	// エビデンス登録ディレクトリの作成
	var eachEvidenceDir = sc.evidencesDir + string(filepath.Separator) + evidenceId.String()
	err = os.MkdirAll(eachEvidenceDir, os.ModePerm)
	if err != nil {
		slog.Error("create evidence registration directory.", err, "path", eachEvidenceDir)
		return fmt.Errorf("create evidence registration directory")
	}

	// エビデンスファイルの登録（ファイルの作成）
	file, err := os.Create(registrationFilePath)
	writeByteCount, err := file.Write(contents)
	if err != nil {
		// エビデンスの登録失敗
		slog.Error("registration failure evidence file.", err, "path", registrationFilePath)
		return fmt.Errorf("registration failure evidence file")
	}
	if writeByteCount != len(contents) {
		// エビデンスの登録失敗（全てが書き込めなかった場合）
		slog.Error("registration failure evidence file.", "path", registrationFilePath)
		return fmt.Errorf("registration failure evidence file")
	}

	// 属するコマンドIDの解決（必須ではない）
	var commandId uuid.UUID
	if command != nil && !reflect.ValueOf(command).IsNil() {
		commandId = command.GetId()
	}

	// ScenarioContextへエビデンス情報の登録
	var evidence = Evidence{
		Id:              evidenceId,
		BelongCommandId: commandId,
		EvidenceType:    evidenceType,
		Name:            name,
		Path:            registrationFilePath,
	}
	if sc.evidences == nil {
		sc.evidences = make(map[uuid.UUID]Evidence)
	}
	sc.evidences[evidenceId] = evidence
	return nil
}

/*
getEvidenceFromCommandId
コマンドIDからエビデンスを取得.
*/
func getEvidenceFromCommandId(sc *ScenarioContext, command Command) (Evidence, error) {
	if sc.evidences == nil {
		return Evidence{}, fmt.Errorf("not found evidence")
	}
	for _, v := range sc.evidences {
		if command.GetId() == v.BelongCommandId {
			return v, nil
		}
	}
	return Evidence{}, fmt.Errorf("not found evidence")
}

/*
getEvidenceFromEvidenceId
エビデンスIDからエビデンスを取得.
*/
func getEvidenceFromEvidenceId(sc *ScenarioContext, evidenceId uuid.UUID) (Evidence, error) {
	if sc.evidences == nil {
		return Evidence{}, fmt.Errorf("not found evidence")
	}
	for k, v := range sc.evidences {
		if evidenceId == k {
			return v, nil
		}
	}
	return Evidence{}, fmt.Errorf("not found evidence")
}
