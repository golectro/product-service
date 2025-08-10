package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func BuildElasticQuery(params url.Values) map[string]any {
	from := 0
	size := 10

	if p := params.Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &from)
		if from > 0 {
			from = (from - 1) * size
		}
	}
	if l := params.Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &size)
	}

	boolQuery := map[string]any{
		"must":   []map[string]any{},
		"filter": []map[string]any{},
	}

	if name := params.Get("name"); name != "" {
		boolQuery["must"] = append(boolQuery["must"].([]map[string]any), map[string]any{
			"match": map[string]any{
				"name": map[string]any{
					"query":     name,
					"fuzziness": "AUTO",
				},
			},
		})
	}

	if cat := params.Get("category"); cat != "" {
		categories := strings.Split(cat, ",")
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]any), map[string]any{
			"terms": map[string]any{
				"category.keyword": categories,
			},
		})
	}

	if brand := params.Get("brand"); brand != "" {
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]any), map[string]any{
			"term": map[string]any{
				"brand.keyword": brand,
			},
		})
	}

	if color := params.Get("color"); color != "" {
		colors := strings.Split(color, ",")
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]any), map[string]any{
			"terms": map[string]any{
				"color.keyword": colors,
			},
		})
	}

	if price := params.Get("price"); price != "" {
		var pVal float64
		fmt.Sscanf(price, "%f", &pVal)
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]any), map[string]any{
			"term": map[string]any{
				"price": pVal,
			},
		})
	}

	priceRange := map[string]any{}
	if min := params.Get("min_price"); min != "" {
		var minVal float64
		fmt.Sscanf(min, "%f", &minVal)
		priceRange["gte"] = minVal
	}
	if max := params.Get("max_price"); max != "" {
		var maxVal float64
		fmt.Sscanf(max, "%f", &maxVal)
		priceRange["lte"] = maxVal
	}
	if len(priceRange) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]map[string]any), map[string]any{
			"range": map[string]any{
				"price": priceRange,
			},
		})
	}

	if specs := params.Get("specs"); specs != "" {
		var specsMap map[string]any
		if err := json.Unmarshal([]byte(specs), &specsMap); err == nil {
			for k, v := range specsMap {
				boolQuery["filter"] = append(boolQuery["filter"].([]map[string]any), map[string]any{
					"term": map[string]any{
						fmt.Sprintf("specs.%s.keyword", k): v,
					},
				})
			}
		}
	}

	sortField := "price"
	sortOrder := "asc"
	if s := params.Get("sort"); s != "" {
		parts := strings.Split(s, ":")
		if len(parts) == 2 {
			sortField = parts[0]
			sortOrder = parts[1]
		}
	}

	query := map[string]any{
		"from": from,
		"size": size,
		"sort": []map[string]string{
			{sortField: sortOrder},
		},
		"query": map[string]any{
			"bool": boolQuery,
		},
	}

	return query
}
