// msg.New("Hello").
//   Add("en", "Hello, world").
//   Add("nl", "Hallo wereld")
// msg.New("Bye").
//   Add("en", "Farewell, cruel world").
//   Add("nl", "Vaarwel, wrede wereld")
package msg

import (
	"net/http"
	"strings"
)

type message map[string]string

func (m message) Add(language, translation string) message {
	language = strings.ToLower(language)
	m[language] = translation
	return m
}

var messageStore = make(map[string]message, 500)

func New(key string, numLang ...int) message {
	size := 2
	if len(numLang) > 0 {
		size = numLang[0]
	}
	m := make(message, size)
	messageStore[key] = m
	return m
}

type languageType struct {
	Full string
	Main string
	Sub  string
}

var languageCache = make(map[string]languageType, 100)

// TODO: be more appreciative to the languages listed in the Accept-Language header;
//   currently only the language first listed is considered
func Language(r *http.Request) (language languageType) {
	acceptLanguage := r.Header.Get("Accept-Language")
	acceptLanguage = strings.ToLower(acceptLanguage)
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

type msgFunc func(key string) (value string)

func Msg(r *http.Request) msgFunc {
	lang := Language(r)
	return func(key string) (value string) {
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
}
