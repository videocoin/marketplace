package datastore

import "strings"

func toTSQuery(q string) string {
	tsQuery := make([]string, 0)
	words := strings.Split(q, " ")
	for _, word := range words {
		w := strings.ReplaceAll(word, " ", "")
		if w != "" {
			tsQuery = append(tsQuery, w)
		}
	}
	return strings.Join(tsQuery, " & ")
}
