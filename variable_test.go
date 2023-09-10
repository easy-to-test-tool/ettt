package ettt

import "testing"

/*
TestReplaceSuccess Replace関数の正常系
*/
func TestReplaceSuccess(t *testing.T) {
	t.Run("Profile変数による置換", func(t *testing.T) {
		var profileVariables = make([]ProfileVariable, 0)
		profileVariables = append(profileVariables, ProfileVariable{
			Key:   "key1",
			Value: "value1",
		})
		var profile = Profile{
			Name:      "local",
			Variables: profileVariables,
		}
		var gc = GlobalContext{
			profile: profile,
		}
		var sc = ScenarioContext{}
		result, err := Replace(gc, sc, "abc${profile.key1}def")
		if err != nil {
			t.Fatalf("failed test %#v", err)
		}
		if result != "abcvalue1def" {
			t.Fatal("failed test")
		}
	})
	t.Run("Profile変数による置換 - 変数２つ", func(t *testing.T) {
		var profileVariables = make([]ProfileVariable, 0)
		profileVariables = append(profileVariables, ProfileVariable{
			Key:   "key1",
			Value: "value1",
		}, ProfileVariable{
			Key:   "key2",
			Value: "value2",
		})
		var profile = Profile{
			Name:      "local",
			Variables: profileVariables,
		}
		var gc = GlobalContext{
			profile: profile,
		}
		var sc = ScenarioContext{}
		result, err := Replace(gc, sc, "abc${profile.key1}def${profile.key2}")
		if err != nil {
			t.Fatalf("failed test %#v", err)
		}
		if result != "abcvalue1defvalue2" {
			t.Fatal("failed test")
		}
	})
	t.Run("Profile変数による置換 - 変数の中に変数", func(t *testing.T) {
		var profileVariables = make([]ProfileVariable, 0)
		profileVariables = append(profileVariables, ProfileVariable{
			Key:   "key1",
			Value: "value1",
		}, ProfileVariable{
			Key:   "key2",
			Value: "123${profile.key3}456",
		}, ProfileVariable{
			Key:   "key3",
			Value: "value3",
		})
		var profile = Profile{
			Name:      "local",
			Variables: profileVariables,
		}
		var gc = GlobalContext{
			profile: profile,
		}
		var sc = ScenarioContext{}
		result, err := Replace(gc, sc, "abc${profile.key1}def${profile.key2}_/?")
		if err != nil {
			t.Fatalf("failed test %#v", err)
		}
		if result != "abcvalue1def123value3456_/?" {
			t.Fatal("failed test")
		}
	})
}
