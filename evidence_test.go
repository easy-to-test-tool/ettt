package ettt

import (
	"github.com/google/uuid"
	"os"
	"testing"
)

const EvidencesTmpDir string = "evidencesTmpDir"
const TextEvidenceFileName string = "dummy.txt"

func TestRegistrationTextEvidence(t *testing.T) {
	t.Run("テキスト形式のエビデンス登録が正常に行えること", func(t *testing.T) {
		var evidenceTmpDir, err = os.MkdirTemp("", EvidencesTmpDir)

		defer os.RemoveAll(evidenceTmpDir)

		var sc = ScenarioContext{
			scenarioName: "dummy",
			evidencesDir: evidenceTmpDir,
		}

		var ro = RegistrationTextEvidenceOptions{
			Name:         TextEvidenceFileName,
			TextContents: "{\"prop\": \"あ\"}",
		}

		// エビデンス登録
		err = registrationTextEvidence(&sc, ro)
		if err != nil {
			t.Fail()
		}

		// アサート
		if len(sc.evidences) != 1 {
			t.Fail()
		}
		for key, value := range sc.evidences {
			if value.Name != TextEvidenceFileName {
				t.Fail()
			}
			if value.Path != (evidenceTmpDir +
				string(os.PathSeparator) +
				key.String() +
				string(os.PathSeparator) +
				TextEvidenceFileName) {
				t.Fail()
			}
			if value.BelongCommandId != uuid.Nil {
				t.Fail()
			}
		}
	})
	t.Run("テキスト形式のエビデンス登録が値が空であるためエラーとなること", func(t *testing.T) {
		var evidenceTmpDir, err = os.MkdirTemp("", EvidencesTmpDir)

		defer os.RemoveAll(evidenceTmpDir)

		var sc = ScenarioContext{
			scenarioName: "dummy",
			evidencesDir: evidenceTmpDir,
		}

		var ro = RegistrationTextEvidenceOptions{
			Name:         TextEvidenceFileName,
			TextContents: "",
		}

		// エビデンス登録
		err = registrationTextEvidence(&sc, ro)
		if err == nil {
			t.Fail()
		}
		if err.Error() != "options.TextContents parameter is required" {
			t.Fail()
		}
	})
}

func TestRegistrationBinaryEvidence(t *testing.T) {
	t.Run("バイナリ形式のエビデンス登録が正常に行えること", func(t *testing.T) {
		var evidenceTmpDir, err = os.MkdirTemp("", EvidencesTmpDir)

		defer os.RemoveAll(evidenceTmpDir)

		var sc = ScenarioContext{
			scenarioName: "dummy",
			evidencesDir: evidenceTmpDir,
		}

		var ro = RegistrationBinaryEvidenceOptions{
			Name:           TextEvidenceFileName,
			BinaryContents: []byte("{\"prop\": \"あ\"}"),
		}

		// エビデンス登録
		err = registrationBinaryEvidence(&sc, ro)
		if err != nil {
			t.Fail()
		}

		// アサート
		if len(sc.evidences) != 1 {
			t.Fail()
		}
		for key, value := range sc.evidences {
			if value.Name != TextEvidenceFileName {
				t.Fail()
			}
			if value.Path != (evidenceTmpDir +
				string(os.PathSeparator) +
				key.String() +
				string(os.PathSeparator) +
				TextEvidenceFileName) {
				t.Fail()
			}
			if value.BelongCommandId != uuid.Nil {
				t.Fail()
			}
		}
	})
	t.Run("バイナリ形式のエビデンス登録が値が空であるためエラーとなること", func(t *testing.T) {
		var evidenceTmpDir, err = os.MkdirTemp("", EvidencesTmpDir)

		defer os.RemoveAll(evidenceTmpDir)

		var sc = ScenarioContext{
			scenarioName: "dummy",
			evidencesDir: evidenceTmpDir,
		}

		var ro = RegistrationBinaryEvidenceOptions{
			Name:           TextEvidenceFileName,
			BinaryContents: nil,
		}

		// エビデンス登録
		err = registrationBinaryEvidence(&sc, ro)
		if err == nil {
			t.Fail()
		}
		if err.Error() != "options.BinaryContents parameter is required" {
			t.Fail()
		}
	})
}
