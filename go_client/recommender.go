package main

import (
	"context"
	"fmt"
	// "log"
	"io"
	"encoding/json"
	"os"
	"net/http"
	"strings"
	"strconv"
	"html/template"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/edgedb/edgedb-go"
)

type Group struct {
	group_name string
}

// EdgeDB setup
var db *edgedb.Client
var ctx context.Context

// Global variables
var results []byte
var userwish = ""
var embedding []float64
var filterQueries []string
var groups []Group
var logicQ []string
var negate = false

func arrayToString(a []float64, delim string) string {
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func lastAndOr() string {
	element := logicQ[len(logicQ)-1]
	return element
}

func initDB() {
	var err error

	opts := edgedb.Options{Concurrency: 4}
	ctx = context.Background()
	db, err = edgedb.CreateClient(ctx, opts)
	if err != nil {
		panic(err)
	}
}

func parseQuality(k string, v string) {
	if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
		filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
	}
	if negate {
		filterQueries = append(filterQueries, "not")
	}


	if v == "<nil>" {
		filterQueries = append(filterQueries, fmt.Sprintf("(not exists .%v)", k))
	} else {
		filterQueries = append(filterQueries, 
							fmt.Sprintf("(.%v = %v)", 
							k, v))
	}
}

func parseQualityComplex(k string, v map[string]interface{}) {
    for key, val := range v {
		switch key {
			case "$in":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				strVal := fmt.Sprintf("%v", val)
				filterQueries = append(filterQueries, 
					fmt.Sprintf("(.%v = %v)", 
					k, strVal[1:len(strVal)-1]))
			case "$nin":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}


				strVal := fmt.Sprintf("%v", val)
				filterQueries = append(filterQueries, 
					fmt.Sprintf("(.%v != %v)", 
					k, strVal[1:len(strVal)-1]))
			
			case "$lte":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, 
					fmt.Sprintf(`(.%v <= to_int32("%v"))`, 
					k, val.(float64)))	
			case "$gte":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, 
					fmt.Sprintf(`(.%v >= to_int32("%v"))`, 
					k, val.(float64)))	
			case "$ne":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, 
					fmt.Sprintf("(.%v != %v)", 
					k, int(val.(float64))))	

			default:
				fmt.Println(val)
		}
    }
}

func parseNumProp(k string, v string) {
	if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
		filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
	}
	if negate {
		filterQueries = append(filterQueries, "not")
	}


	if v == "<nil>" {
		filterQueries = append(filterQueries, fmt.Sprintf("(not exists .%v)", k))
	} else {
		filterQueries = append(filterQueries, 
							fmt.Sprintf("(.%v >= <decimal>( select %v * 0.95 ) and .%v <= <decimal>( select %v * 1.05 ))", 
							k, v, k, v))
	}
}

func parseNumPropComplex(k string, v map[string]interface{}) {
    for key, val := range v {
		switch key {
			case "$in":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}


				strVal := fmt.Sprintf("%v", val)
				filterQueries = append(filterQueries, 
					fmt.Sprintf("(.%v >= <decimal>( select %v * 0.95 ) and .%v <= <decimal>( select %v * 1.05 ))", 
					k, strVal[1:len(strVal)-1], k, strVal[1:len(strVal)-1])[:len(lastAndOr())-2])
			case "$nin":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				strVal := fmt.Sprintf("%v", val)
				filterQueries = append(filterQueries, 
					fmt.Sprintf("(.%v < <decimal>( select %v * 0.95 ) or .%v > <decimal>( select %v * 1.05 ))", 
					k, strVal[1:len(strVal)-1], k, strVal[1:len(strVal)-1])[:len(lastAndOr())-2])
			
			case "$lte":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, 
					fmt.Sprintf(`(.%v <= to_decimal("%v"))`, 
					k, val.(float64)))	
			case "$gte":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, 
					fmt.Sprintf(`(.%v >= to_decimal("%v"))`, 
					k, val.(float64)))	
			case "$ne":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, 
					fmt.Sprintf("(.%v < <decimal>( select %v * 0.95 ) or .%v > <decimal>( select %v * 1.05 ))", 
					k, int(val.(float64)), k, int(val.(float64))))	

			default:
				fmt.Println(val)
		}
    }
}

func parseCountry(country string) {
	if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
		filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
	}
	if negate {
		filterQueries = append(filterQueries, "not")
	}

	filterQueries = append(filterQueries, "(.originID.country", "=", fmt.Sprintf("'%v')", country))	
}

func parseCountryComplex(country map[string]interface{}) {
    for key, val := range country {
		switch key {
			case "$ne":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, "(.originID.country", "!=", fmt.Sprintf("'%v')", val))
			default:
				fmt.Println(val)
		}
    }
}

func popGrops() {
	initDB()
	defer db.Close()

	query := `select Groups {group_name}
			order by .group_name`

	err := db.Query(ctx, query, &groups)
    if err != nil {
        panic(err)
	}
}

func parseGroupValues(anArray []interface{}) {
	ptr := &groups
	var vals []string
    for _, val := range anArray {
        switch concreteVal := val.(type) {
			case map[string]interface{}:
				parseGroups(val.(map[string]interface{}))
			case []interface{}:
				parseGroupValues(val.([]interface{}))
			default:
				v := fmt.Sprint((*ptr)[int(concreteVal.(float64))-1])
				vals = append(vals, fmt.Sprintf("'%v'", v[1:len(v)-1]))
        }
    }
	filterQueries = append(filterQueries, strings.Join(vals, ", ") + " }))")
}


func parseGroups(categories map[string]interface{}) {
    for key, val := range categories {
		switch key {
			case "$in":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, "(any(.groupID.group_name in {")
				parseGroupValues(val.([]interface{}))
			case "$nin":
				if len(filterQueries) > 0 && (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
					filterQueries = append(filterQueries, lastAndOr()[:len(lastAndOr())-2])
				}
				if negate {
					filterQueries = append(filterQueries, "not")
				}

				filterQueries = append(filterQueries, "(not any(.groupID.group_name in {")
				parseGroupValues(val.([]interface{}))
			default:
				fmt.Println(val)
		}
    }
}

func lookupGroupID(categories map[string]interface{}) {
	if len(groups) == 0 {
		popGrops()
		ptr := &groups
		// remove unknown group
		groups = append((*ptr)[:17], (*ptr)[18:]...)
	}
	parseGroups(categories)
}

func findReview(page int, rows int) {
	offset := (page-1)*rows
	query := ""
	if len(filterQueries) == 0 {
		query = fmt.Sprintf(
			`WITH
				maxcount := 20,
				page := %v,
				rows := %v,
				embedding := <REQUIRED TxEmbedding>[%s],
				dream := ( select Reviews {
					breed := (select Breed {name} filter Reviews.BreedReview.id = .id)
				    } 
				    order by ext::pgvector::cosine_distance(Reviews.embedding, embedding)
                    empty last
				    limit maxcount
                ),
				selected := (select Breed
					filter .name in dream.breed.name
				),
				select {
					total:= count(selected),
					rows := (
						select Breed {*}
						filter .id = selected.id
						order by .name
						offset page
						limit rows)
				}`,
				offset,
				rows,
				arrayToString(embedding, ", "))
	} else {
		filter := strings.Join(filterQueries, " ")
		conn := lastAndOr()
		query = fmt.Sprintf(
			`WITH
				maxcount := 89,
				page := %v,
				rows := %v,
				embedding := <REQUIRED TxEmbedding>[%s],
				dream := ( select Reviews {
					breed := (select Breed {name} filter Reviews.BreedReview.id = .id)
				    } 
				    order by ext::pgvector::cosine_distance(Reviews.embedding, embedding)
                    empty last
				    limit maxcount
                ),
				selected := (select Breed
					filter %v %v .name in dream.breed.name)
				),
				select {
					total:= count(selected),
					rows := (
						select Breed {*}
						filter .id = selected.id
						order by .name
						offset page
						limit rows)
				}`,
				offset,
				rows,
				arrayToString(embedding, ", "),
				filter,
				conn)
	}
	initDB()
	defer db.Close()

    err := db.QueryJSON(ctx, query, &results)
    if err != nil {
        panic(err)
    }
}

func findMatch(page int, rows int) {
	offset := (page-1)*rows
	filter := strings.Join(filterQueries, " ")

	initDB()
	defer db.Close()

	query := fmt.Sprintf(
			`WITH
				page := %v,
				rows := %v,
				selected := (select Breed
					filter %v
					),
				select {
					total:= count(selected),
					rows := (select Breed {*}
							filter .id = selected.id
							order by .name
							offset page
							limit rows)
				}`,
				offset,
				rows,
				filter)
    err := db.QueryJSON(ctx, query, &results)
    if err != nil {
        panic(err)
    }
}

func hfModel() {
	ic := huggingface.NewInferenceClient(os.Getenv("HUGGINGFACEHUB_API_TOKEN"))

	feat, err := ic.FeatureExtractionWithAutomaticReduction(context.Background(), &huggingface.FeatureExtractionRequest{
		Inputs: []string{userwish},
		Model:  "sentence-transformers/all-mpnet-base-v2",
		Options: huggingface.Options{
			WaitForModel: huggingface.PTR(true),
		},
	})
	if err != nil {
		panic(err)
	}

	embedding = feat[0]

}

func parseMap(aMap map[string]interface{}, page int, rows int) {
    for key, val := range aMap {
        switch concreteVal := val.(type) {
			case map[string]interface{}:
				// fmt.Println(key)
				switch key {
					case "category":
						lookupGroupID(concreteVal)
					case "country":
						parseCountryComplex(concreteVal)
					case "age",
						 "weight",
						 "height":
						if key == "age" {
							parseNumPropComplex("life_expectancy", concreteVal)
						} else {
							parseNumPropComplex(key, concreteVal)
						}
					case "playfulness",
		  				 "barking",
					 	 "drooling",
					 	 "grooming",
					 	 "good_with_children",
					 	 "good_with_strangers",
					 	 "good_with_other_dogs",
						 "shedding",
						 "coat_length",
						 "protectiveness",
						 "trainability",
						 "energy":
						parseQualityComplex(key, concreteVal)
		
					default:
            			parseMap(val.(map[string]interface{}), page, rows)
				}
        case []interface{}:
			switch key {
				case "$and":
					if len(logicQ)>0 {
						// dequeue
						if len(logicQ)>1 {
							logicQ = logicQ[:len(logicQ)-1]
						}
						if len(filterQueries) > 0 {
							qry := fmt.Sprintf("%v", strings.Join(filterQueries, " "))
							filterQueries = make([]string, 0)
							filterQueries = append(filterQueries, qry)
							if (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
								filterQueries = append(filterQueries, lastAndOr())
							}
						}
					}
					logicQ = append(logicQ, "and (")
					parseArray(val.([]interface{}), page, rows)
				case "$or":
					if len(logicQ)>0 {
						// dequeue
						if len(logicQ)>1 {
							logicQ = logicQ[:len(logicQ)-1]
						}
						if len(filterQueries) > 0 {
								qry := fmt.Sprintf("%v", strings.Join(filterQueries, " "))
							filterQueries = make([]string, 0)
							filterQueries = append(filterQueries, qry)
							if (filterQueries[len(filterQueries)-1] != "and (" && filterQueries[len(filterQueries)-1] != "or (") {
								filterQueries = append(filterQueries, lastAndOr())
							}
						}
					}
					logicQ = append(logicQ, "or (")
					parseArray(val.([]interface{}), page, rows)
				case "$nor":
					negate = !negate
					parseArray(val.([]interface{}), page, rows)
			}
		
        default:
			switch key {
				case "userwish":
					userwish = fmt.Sprintf("%v", concreteVal)
				case "country":
					parseCountry(fmt.Sprintf("%v", concreteVal))
				case "age",
					 "weight",
					 "height":
					if key == "age" {
						parseNumProp("life_expectancy", fmt.Sprintf("%v", concreteVal))
					} else {
						parseNumProp(key, fmt.Sprintf("%v", concreteVal))
					}
				case "playfulness",
					 "barking",
					 "drooling",
					 "grooming",
					 "good_with_children",
					 "good_with_strangers",
					 "good_with_other_dogs",
					 "shedding",
					 "coat_length",
					 "protectiveness",
					 "trainability",
					 "energy":
					parseQuality(key, fmt.Sprintf("%v", concreteVal))
			}
        }
    }
}

func parseArray(anArray []interface{}, page int, rows int) {
    for i, val := range anArray {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
			beforeCount := len(filterQueries)
			boolAndOrExist := false
			for _, token := range filterQueries {
				if token == "and (" || token == "or (" {
					boolAndOrExist = true
				}
			}
			
            parseMap(val.(map[string]interface{}), page, rows)

			// dequeue
			if len(anArray)>1 && i == len(anArray)-1 && len(logicQ)>1 {
				logicQ = logicQ[:len(logicQ)-1]
			}
			// close parenthesis
			if len(anArray)>1 && i == len(anArray)-1 && filterQueries[len(filterQueries)-1] != ")" && len(filterQueries)>beforeCount && boolAndOrExist {
				filterQueries = append(filterQueries, ")")
			} else {
				if i == len(anArray)-1 && boolAndOrExist {
					filterQueries = append(filterQueries, ")")
				}
			}
        case []interface{}:
            parseArray(val.([]interface{}), page, rows)
        default:
            fmt.Println("Index", i, ":", concreteVal)

        }
    }
}

// POST /recommendBreed
func recommendBreed(writer http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
    if err != nil {
        panic(err)
    }
	rows, err := strconv.Atoi(req.URL.Query().Get("rows"))
    if err != nil {
        panic(err)
    }

	var jsonBody map[string]interface{}

	body, _ := io.ReadAll(req.Body)
	json.Unmarshal([]byte(body), &jsonBody)

	userwish = ""
	filterQueries = make([]string, 0)
	logicQ = make([]string, 0)
	negate = false

	parseMap(jsonBody, page, rows)
	
	// dequeue
	if len(logicQ)>1 {
		logicQ = logicQ[:len(logicQ)-1]
	}
	if userwish == "" {
		findMatch(page, rows)
	} else {
		hfModel()
		findReview(page, rows)
	}


	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(string(results[:]))

}

// GET /getBreedDetail
func getBreedDetail(writer http.ResponseWriter, req *http.Request) {
	var name = req.URL.Query().Get("name")

	initDB()
	defer db.Close()

	query := fmt.Sprintf(
			`WITH
				breedname := <REQUIRED str>"%s",
				select Breed {
					*,
					imageID:  { image_link },
					originID: { country },
					reviewID: { review },
					groupID: { group_name }
				}				  
				filter .name in breedname`, 
				name)
    err := db.QueryJSON(ctx, query, &results)
    if err != nil {
        panic(err)
    }

    m := map[string]interface{}{}
    if err := json.Unmarshal([]byte(results[1:len(results)-1]), &m); err != nil {
        panic(err)
    }

	const templ = `
			<div onclick="togglePopup()" class="close-btn">
            X
	        </div>

			<h2>{{ .name }}</h2>
		
			<table class="dv-table" border="0" style="width:100%;">
				<tr>
					<td rowspan="6" style="width:480px">
						{{- range $key, $val := .imageID }}
							{{- range $k, $v := $val }}
								<img src={{ $v }} style="width:480px;margin-right:20px" />
							{{- end }}
						{{- end }}
					</td>
					<td class="dv-label">Avg. height:</td>
					<td>{{ .height }}</td>
					<td class="dv-label">Avg. weight:</td>
					<td>{{ .weight }}</td>
					<td class="dv-label">Barking:</td>
					<td>{{ .barking }}</td>
				</tr>
				<tr>
					<td class="dv-label">Coat length:</td>
					<td>{{ .coat_length }}</td>
					<td class="dv-label">Drooling:</td>
					<td>{{ .drooling }}</td>
					<td class="dv-label">Energy:</td>
					<td>{{ .energy }}</td>
				</tr>
				<tr>
					<td class="dv-label">Grooming:</td>
					<td>{{ .grooming }}</td>
					<td class="dv-label">Life Expectancy:</td>
					<td>{{ .life_expectancy }}</td>
					<td class="dv-label">Playfulness:</td>
					<td>{{ .playfulness }}</td>
				</tr>
				<tr>
					<td class="dv-label">Good with children:</td>
					<td>{{ .good_with_children }}</td>
					<td class="dv-label">Good with other dogs:</td>
					<td>{{ .good_with_other_dogs }}</td>
					<td class="dv-label">Good with strangers:</td>
					<td>{{ .good_with_strangers }}</td>
				</tr>
				<tr>
					<td class="dv-label">Protectiveness:</td>
					<td>{{ .protectiveness }}</td>
					<td class="dv-label">Shedding:</td>
					<td>{{ .shedding }}</td>
					<td class="dv-label">Trainability:</td>
					<td>{{ .trainability }}</td>
				</tr>
				<tr>
					<td class="dv-label">Country of Origin:</td>
					<td>{{ .originID.country }}</td>
					<td class="dv-label">Group Name(s):</td>
					<td colspan="2"><ul class="dv-ul">
					{{- range $key, $val := .groupID }}
						{{- range $k, $v := $val }}
							<li>{{ $v }}</li>
						{{- end }}
					{{- end }}
					</ul></td>
				</tr>
				<tr>
					<td class="dv-label">Review(s):</td>
				</tr>
				<tr>
					<td colspan="7">
			{{- range $key, $val := .reviewID }}
				{{- range $k, $v := $val }}
						<textarea class="dv-textarea">{{ $v }}</textarea>
				{{- end }}
			{{- end }}
					</td>
				</tr>
			</table>`

    t := template.Must(template.New("").Parse(templ))

    if err := t.Execute(writer, m); err != nil {
        panic(err)
    }

}
