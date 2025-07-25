package core

import (
	"fmt"
	"infini.sh/framework/core/elastic"
	"infini.sh/framework/core/errors"
	"infini.sh/framework/core/orm"
	"infini.sh/framework/core/util"
	"reflect"
)

func SearchV2WithResultItemMapper(ctx *orm.Context, resultArray interface{}, qb *orm.QueryBuilder, itemMapFunc func(source map[string]interface{}, targetRef interface{}) error) (error, *orm.SimpleResult) {
	response, err := orm.SearchV2(ctx, qb)
	if err != nil {
		return err, nil
	}
	if response == nil {
		return errors.New("invalid response"), nil
	}

	bytes, ok := response.Payload.([]byte)
	if ok {
		// Validate that resultArray is a pointer to a slice
		arrayValue := reflect.ValueOf(resultArray)
		if arrayValue.Kind() != reflect.Ptr || arrayValue.Elem().Kind() != reflect.Slice {
			return fmt.Errorf("resultArray must be a pointer to a slice"), nil
		}

		sliceValue := arrayValue.Elem()
		elementType := sliceValue.Type().Elem() // Get the type of elements in the slice

		searchResponse := elastic.SearchResponse{}
		err := util.FromJSONBytes(bytes, &searchResponse)
		if err != nil {
			return err, nil
		}
		// Populate the resultArray with typed data
		for _, doc := range searchResponse.Hits.Hits {
			// Create a new instance of the target element type
			elem := reflect.New(elementType).Elem()

			//make sure id exists and always be _id
			doc.Source["id"] = doc.ID

			if itemMapFunc != nil {

				source := doc.Source
				// Map the document source into the element
				if err := itemMapFunc(source, elem.Addr().Interface()); err != nil { // Ensure passing a pointer to itemMapFunc
					return fmt.Errorf("failed to map document to struct: %w", err), nil
				}
			}

			// Append the populated element to the result slice
			sliceValue.Set(reflect.Append(sliceValue, elem))
		}

		result := orm.SimpleResult{}
		result.Total = searchResponse.GetTotal()
		result.Raw = util.MustToJSONBytes(searchResponse)
		return nil, &result
	}

	return errors.New("invalid response"), nil
}
