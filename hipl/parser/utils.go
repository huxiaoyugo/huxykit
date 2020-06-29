package parser

import "strings"

func getStringVal(lit string) string {
	if len(lit) == 0 {
		return lit
	}
	if lit[0] == '\'' {
		return strings.Trim(lit,"'")
	} else {
		return strings.Trim(lit,"\"")
	}
}


func lower(ch rune) rune     { return ('a' - 'A') | ch }

func upper(ch rune) rune     {

	return ('a' - 'A') | ch

}