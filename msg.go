package msg

import (
	"net/http"
	"strings"
)

type message map[string]string

var messageStore = make(map[string]message, 500)

// var m, a = msg.Definition()
// m("Hello")
// a("en", "Hello, world")
// a("nl", "Hallo wereld")
// m("Bye")
// a("en", "Farewell, cruel world")
// a("nl", "Vaarwel, wrede wereld")
func Definition() (createMessage func(string), addTranslation func(string, string)) {
	var msg message
	createMessage = func(key string) {
		msg = make(message, 2)
		messageStore[key] = msg
	}
	addTranslation = func(language string, translation string) {
		msg[language] = translation
	}
	return
}

type languageType struct {
	Full string
	Main string
	Sub  string
}

var languageCache = make(map[string]languageType, 100)

func getAcceptLanguage(r *http.Request) string {
	acceptLanguage := r.Header.Get("Accept-Language")
	return strings.ToLower(acceptLanguage)
}

// TODO: be more appreciative to the languages listed in the Accept-Language header;
//   currently only the language first listed is considered
func Language(r *http.Request) (language languageType) {
	acceptLanguage := getAcceptLanguage(r)
	if lang, ok := languageCache[acceptLanguage]; ok {
		language = lang
	} else {
		firstLanguage := strings.Split(acceptLanguage, ",")[0] // cut other languages
		firstLanguage = strings.Split(firstLanguage, ";")[0]   // cut the q parameter
		parts := strings.Split(firstLanguage, "-")
		lang = languageType{
			Full: firstLanguage,
			Main: parts[0],
		}
		if len(parts) > 1 {
			lang.Sub = parts[1]
		}
		languageCache[acceptLanguage] = lang
		language = lang
	}
	return
}

type msgFuncType func(key string) (value string)

var msgFuncCache = make(map[string]msgFuncType, 100)

func Msg(r *http.Request) (msgFunc msgFuncType) {
	acceptLanguage := getAcceptLanguage(r)
	if f, ok := msgFuncCache[acceptLanguage]; ok {
		msgFunc = f
	} else {
		lang := Language(r)
		f = func(key string) (value string) {
			if val, ok := messageStore[key][lang.Full]; ok {
				value = val
			} else if val, ok := messageStore[key][lang.Sub]; ok {
				value = val
			} else if val, ok := messageStore[key][lang.Main]; ok {
				value = val
			} else {
				value = "X-" + key
			}
			return
		}
		msgFuncCache[acceptLanguage] = f
		msgFunc = f
	}
	return
}
