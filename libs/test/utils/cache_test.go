package utils_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"gitlab.com/youwol/platform/libs/go-libs/utils"
)

type CustomType struct {
	Indexes []int
	Name    string
	HowLong time.Duration
}

// TODO: test redis cache

func TestSingleLocalCacheStringValue(t *testing.T) {
	err := utils.InitCache("toto", utils.Local, 10*time.Minute)
	if err != nil {
		t.Errorf("Error initializing cache (%v)", err)
	}

	utils.SetCache(context.Background(), "toto", "gloups", "glamuk")
	var result string = "plip"
	bFound := utils.GetCache(context.Background(), "toto", "gloups", &result)
	if !bFound {
		t.Errorf("Did not find requested value in cache")
	}

	if result != "glamuk" {
		t.Errorf("Wrong value returned by the cache engine")
	}
}

func TestSingleLocalCacheInterfaceValue(t *testing.T) {
	err := utils.InitCache("toto", utils.Local, 10*time.Minute)
	if err != nil {
		t.Errorf("Error initializing cache (%v)", err)
	}

	ref := CustomType{
		Indexes: []int{1, 2, 3, 4, 6},
		Name:    "gloups",
		HowLong: 10 * time.Second,
	}

	utils.SetCache(context.Background(), "toto", "gloups", ref)
	result := &CustomType{}
	bFound := utils.GetCache(context.Background(), "toto", "gloups", &result)
	if !bFound {
		t.Errorf("Did not find requested value in cache")
	}

	if !reflect.DeepEqual(*result, ref) {
		t.Errorf("Wrong value returned by the cache engine")
	}
}

func TestMultipleLocalCacheStringValue(t *testing.T) {
	err := utils.InitCache("toto", utils.Local, 10*time.Minute)
	if err != nil {
		t.Errorf("Error initializing cache (%v)", err)
	}
	err = utils.InitCache("titi", utils.Local, 10*time.Minute)
	if err != nil {
		t.Errorf("Error initializing cache (%v)", err)
	}

	utils.SetCache(context.Background(), "toto", "gloups", "glamuk")
	utils.SetCache(context.Background(), "titi", "gloups", "glamuk2")
	var result string = "plip"
	bFound := utils.GetCache(context.Background(), "toto", "gloups", &result)
	if !bFound {
		t.Errorf("Did not find requested value in cache")
	}
	if result != "glamuk" {
		t.Errorf("Wrong value returned by the cache engine")
	}

	bFound = utils.GetCache(context.Background(), "titi", "gloups", &result)
	if !bFound {
		t.Errorf("Did not find requested value in cache")
	}
	if result != "glamuk2" {
		t.Errorf("Wrong value returned by the cache engine")
	}
}
