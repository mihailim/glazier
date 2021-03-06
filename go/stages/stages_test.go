// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package stages

import (
	"fmt"
	"testing"

	"golang.org/x/sys/windows/registry"
)

const (
	testStageRoot = `SOFTWARE\Glazier\Testing`
)

func createTestKeys(subKeys ...string) error {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, testStageRoot, registry.CREATE_SUB_KEY)
	if err != nil {
		return err
	}
	defer k.Close()
	for _, id := range subKeys {
		sk, _, err := registry.CreateKey(k, id, registry.CREATE_SUB_KEY)
		if err != nil {
			return err
		}
		defer sk.Close()
		k = sk
	}
	return nil
}

func cleanupTestKey() error {
	return registry.DeleteKey(registry.LOCAL_MACHINE, testStageRoot)
}

func TestGetActiveStageNoRootKey(t *testing.T) {
	testID := "TestGetActiveStageNoRootKey"
	stage, err := getActiveStage(testStageRoot + `\` + testID)
	if err != nil {
		t.Errorf("%s(): raised unexpected error %v", testID, err)
	}
	if stage != "0" {
		t.Errorf("%s(): got %s, want %s", testID, stage, "0")
	}
}

func TestGetActiveStageNoActiveKey(t *testing.T) {
	testID := "TestGetActiveStageNoActiveKey"
	if err := createTestKeys(testID); err != nil {
		t.Fatal(err)
	}
	defer cleanupTestKey()

	stage, err := getActiveStage(testStageRoot + `\` + testID)
	if err != nil {
		t.Errorf("%s(): raised unexpected error %v", testID, err)
	}
	if stage != "0" {
		t.Errorf("%s(): got %s, want %s", testID, stage, "0")
	}
}

func TestGetActiveStageInProgress(t *testing.T) {
	testID := "TestGetActiveStageInProgress"
	subKey := testStageRoot + `\` + testID

	if err := createTestKeys(testID); err != nil {
		t.Fatal(err)
	}
	defer cleanupTestKey()

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, subKey, registry.WRITE)
	if err != nil {
		t.Fatal(err)
	}
	if err = k.SetStringValue("_Active", "2"); err != nil {
		t.Fatal(err)
	}
	k.Close()

	stage, err := getActiveStage(subKey)
	if err != nil {
		t.Errorf("%s(): raised unexpected error %v", testID, err)
	}
	if stage != "2" {
		t.Errorf("%s(): got %s, want %s", testID, stage, "2")
	}
}

func TestGetActiveStageTypeError(t *testing.T) {
	testID := "TestGetActiveStageTypeHandling"
	subKey := testStageRoot + `\` + testID

	if err := createTestKeys(testID); err != nil {
		t.Fatal(err)
	}
	defer cleanupTestKey()

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, subKey, registry.WRITE)
	if err != nil {
		t.Fatal(err)
	}
	if err = k.SetDWordValue("_Active", 0); err != nil {
		t.Fatal(err)
	}
	k.Close()

	if _, err := getActiveStage(subKey); err == nil {
		t.Errorf("%s(): failed to raise expected error", testID)
	}
}

func TestGetActiveTime(t *testing.T) {
	testID := "TestGetActiveTime"
	testKey := fmt.Sprintf(`%s\%s`, testStageRoot, testID)
	stageKey := fmt.Sprintf(`%s\%d`, testKey, 5)

	if err := createTestKeys(testID, "5"); err != nil {
		t.Fatal(err)
	}
	defer cleanupTestKey()

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, stageKey, registry.WRITE)
	if err != nil {
		t.Fatal(err)
	}
	if err = k.SetStringValue("Start", "2019-11-06T17:37:43.279253"); err != nil {
		t.Fatal(err)
	}
	k.Close()
	_, err = getActiveTime(testKey, "5")
	if err != nil {
		t.Errorf("%s(): raised unexpected error %v", testID, err)
	}
}

func TestGetActiveTimeParseError(t *testing.T) {
	testID := "TestGetActiveTimeParseError"
	testKey := fmt.Sprintf(`%s\%s`, testStageRoot, testID)
	stageKey := fmt.Sprintf(`%s\%d`, testKey, 5)

	if err := createTestKeys(testID, "5"); err != nil {
		t.Fatal(err)
	}
	defer cleanupTestKey()

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, stageKey, registry.WRITE)
	if err != nil {
		t.Fatal(err)
	}
	if err = k.SetStringValue("Start", "20191106V17:37:43"); err != nil {
		t.Fatal(err)
	}
	k.Close()
	_, err = getActiveTime(testKey, "5")
	if err == nil {
		t.Errorf("%s(): failed to raise expected error", testID)
	}
}

func TestGetActiveTimeNoKey(t *testing.T) {
	testID := "TestGetActiveTimeNoKey"
	_, err := getActiveTime(testStageRoot, "3")
	if err == nil {
		t.Errorf("%s(): failed to raise expected error", testID)
	}
}
