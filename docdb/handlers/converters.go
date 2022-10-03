package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btubbs/datetime"
	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
)

var mapJSON2CQL = make(map[string](func(interface{}) (interface{}, error)))
var mapCQL2JSON = make(map[string](func(interface{}) (interface{}, error)))

func json2UUID(src interface{}) (interface{}, error) {
	var err error
	ret, ok := src.(gocql.UUID)
	if !ok {
		var cqlUUID gocql.UUID
		cqlUUID, err = gocql.ParseUUID(src.(string))
		ret = cqlUUID
	}

	return ret, err
}

func json2Decimal(src interface{}) (interface{}, error) {
	val := &inf.Dec{}
	err := val.UnmarshalText([]byte(src.(json.Number)))
	return val, err
}

func json2BigInt(src interface{}) (interface{}, error) {
	return strconv.ParseInt(src.(json.Number).String(), 10, 64)
}

func json2Int(src interface{}) (interface{}, error) {
	return strconv.ParseInt(src.(json.Number).String(), 10, 32)
}

func json2SmallInt(src interface{}) (interface{}, error) {
	return strconv.ParseInt(src.(json.Number).String(), 10, 16)
}

func json2TinyInt(src interface{}) (interface{}, error) {
	return strconv.ParseInt(src.(json.Number).String(), 10, 16)
}

func json2Float(src interface{}) (interface{}, error) {
	f, err := strconv.ParseFloat(src.(json.Number).String(), 32)
	return float32(f), err
}

func json2Double(src interface{}) (interface{}, error) {
	return strconv.ParseFloat(src.(json.Number).String(), 64)
}

func json2Timestamp(src interface{}) (interface{}, error) {
	return datetime.Parse(src.(string), time.UTC)
}

func json2Duration(src interface{}) (interface{}, error) {
	return time.ParseDuration(src.(string))
}

func json2Time(src interface{}) (interface{}, error) {
	buf2 := "Jan 2 00:00:00.000000000"
	buf := "Jan 2 " + src.(string)
	t, err := time.Parse(time.StampNano, buf)
	t2, _ := time.Parse(time.StampNano, buf2)
	return t.Sub(t2), err
}

func decimal2JSON(src interface{}) (interface{}, error) {
	val := src.(*inf.Dec)
	strVal := val.String()
	if strings.Contains(strVal, ".") {
		return strconv.ParseFloat(strVal, 32)
	}
	return strconv.ParseInt(strVal, 10, 32)
}

func date2JSON(src interface{}) (interface{}, error) {
	t := src.(time.Time)
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()), nil
}

func time2JSON(src interface{}) (interface{}, error) {
	return src.(time.Duration).String(), nil
}

func duration2JSON(src interface{}) (interface{}, error) {
	dur := src.(gocql.Duration)
	return time.Duration(dur.Nanoseconds).String(), nil
}

func timestamp2JSON(src interface{}) (interface{}, error) {
	t := src.(time.Time)
	timeWithoutZone := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%03d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1e6)
	_, offset := t.Zone()
	return fmt.Sprintf("%s%+05d", timeWithoutZone, offset/3600), nil
}

// InitConverters initializes the converters maps
// It must be called only once at startup
func InitConverters() {
	mapJSON2CQL["uuid"] = json2UUID
	mapJSON2CQL["timeuuid"] = json2UUID
	mapJSON2CQL["decimal"] = json2Decimal
	mapJSON2CQL["bigint"] = json2BigInt
	mapJSON2CQL["varint"] = json2BigInt
	mapJSON2CQL["int"] = json2Int
	mapJSON2CQL["smallint"] = json2SmallInt
	mapJSON2CQL["tinyint"] = json2TinyInt
	mapJSON2CQL["float"] = json2Float
	mapJSON2CQL["double"] = json2Double
	mapJSON2CQL["timestamp"] = json2Timestamp
	mapJSON2CQL["duration"] = json2Duration
	mapJSON2CQL["time"] = json2Time

	mapCQL2JSON["decimal"] = decimal2JSON
	mapCQL2JSON["date"] = date2JSON
	mapCQL2JSON["time"] = time2JSON
	mapCQL2JSON["duration"] = duration2JSON
	mapCQL2JSON["timestamp"] = timestamp2JSON
}

func getCollectionAndSubType(validator string) (string, string) {
	var subType = validator
	var mode string
	if strings.HasPrefix(subType, "map<") {
		// We assume that the key is always convertible
		subTypeArr := strings.Split(subType, ", ")
		subType = strings.TrimRight(subTypeArr[1], ">")
		mode = "map"
	}
	if strings.HasPrefix(subType, "set<") {
		subType = strings.TrimRight(subType, ">")
		subType = strings.TrimPrefix(subType, "set<")
		mode = "set"
	}
	if strings.HasPrefix(subType, "list<") {
		subType = strings.TrimRight(subType, ">")
		subType = strings.TrimPrefix(subType, "list<")
		mode = "list"
	}

	return mode, subType
}

func convertCollectionOrPrimitive(src interface{}, mode string, converter func(interface{}) (interface{}, error)) (interface{}, error) {
	var err error
	switch mode {
	case "map":
		dMap := src.(map[string]interface{})
		for key, val := range dMap {
			dMap[key], err = converter(val)
		}
		return dMap, err
	case "list":
		dList := src.([]interface{})
		for i, val := range dList {
			dList[i], err = converter(val)
		}
		return dList, err
	case "set":
		dSet := src.([]interface{})
		for i, val := range dSet {
			dSet[i], err = converter(val)
		}
		return dSet, err
	default: // convert the primitive value (no collection)
		return converter(src)
	}
}

// This function converts swagger (json) types into cql or native golang types
// All json primitive types are actually strings and do not marshall correctly into cassandra
// TODO:
// - benchmark me and find the best lookup order/strategy
func jsonToCql(dataMap *map[string]interface{}, cols *map[string]*gocql.ColumnMetadata) error {
	var err error
	for k, v := range *cols {
		// all data columns are not always present (no field is required except partition & clustering keys)
		data := *dataMap
		d, found := data[k]
		if !found {
			// set default value
			data[k] = nil
			continue
		}

		// Check if the type is a collection
		mode, subType := getCollectionAndSubType(v.Type)
		converter, found := mapJSON2CQL[subType]
		if found { // A conversion is necessary
			data[k], err = convertCollectionOrPrimitive(d, mode, converter)
		}

		if err != nil {
			break
		}
	}

	return err
}

// This function converts cql types into swagger/json types
// TODO:
// - benchmark me
// - handle collections properly
func cqlToJSON(dataMap *map[string]interface{}, cols *map[string]*gocql.ColumnMetadata) error {
	var err error
	for k, d := range *dataMap {
		v, found := (*cols)[k]
		if !found {
			return fmt.Errorf("Could not find column name in the table meta-data")
		}

		// Check if the type is a collection
		mode, subType := getCollectionAndSubType(v.Type)
		converter, found := mapCQL2JSON[subType]
		if found { // A conversion is necessary
			(*dataMap)[k], err = convertCollectionOrPrimitive(d, mode, converter)
		}

		if err != nil {
			break
		}
	}

	return err
}
