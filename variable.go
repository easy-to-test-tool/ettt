package ettt

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

/*
変数を抽出する正規表現.
*/
var re = regexp.MustCompile(`\$\{([^\}]+)\}`)

/*
Resolve 変数を解決.
スコープ指定がされている場合は、該当スコープのみを走査して解決
スコープ指定がされていない場合に、Profile＞Storeの順序で走査を行い解決
*/
func Resolve(gc GlobalContext, sc ScenarioContext, target string) (string, error) {
	targetArray := strings.Split(target, VariableScopeSeparator)
	if 1 == len(targetArray) {
		slog.Debug("epion-t3: try resolve from all scope variable.", "target", target)
		for _, v := range gc.profile.Variables {
			if v.Key == targetArray[0] {
				return v.Value, nil
			}
		}
		for k, v := range sc.Store.Variables {
			if k == targetArray[0] {
				return v, nil
			}
		}
	} else if 2 == len(targetArray) {
		slog.Debug("epion-t3: try resolve scope variable.", "scope", targetArray[0], "target", targetArray[1])
		if ScopeNameProfile == targetArray[0] {
			// Profile変数から解決する
			for _, v := range gc.profile.Variables {
				if v.Key == targetArray[1] {
					return v.Value, nil
				}
			}
		} else if ScopeNameStore == targetArray[0] {
			// Store変数から解決する
			for k, v := range sc.Store.Variables {
				if k == targetArray[1] {
					return v, nil
				}
			}
		}
	}
	// 見つからない場合は、空文字+エラーを変役
	return "", fmt.Errorf("can not resolve target. target : %s", target)
}

/*
Replace 変数の置換処理.
引数で与えられた文字列から ${スコープ.変数名} 部分を全て置換する.
*/
func Replace(gc GlobalContext, sc ScenarioContext, target string) (string, error) {
	fss := re.FindStringSubmatch(target)
	if 0 != len(fss) {
		v, err := Resolve(gc, sc, fss[1])
		if err != nil {
			return target, err
		}
		target = strings.Replace(target, fss[0], v, -1)
		fss = re.FindStringSubmatch(target)
		if 0 != len(fss) {
			return Replace(gc, sc, target)
		} else {
			return target, nil
		}
	}
	return target, nil
}
