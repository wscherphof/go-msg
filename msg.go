/*
Package msg provides a means to manage translations of text labels ("messages") in a web application.

New messages are defined like this:
	msg.New("Hello").
	  Add("en", "Hello, world").
	  Add("nl", "Hallo wereld")
	msg.New("Hi").
	  Add("en", "Hi").
	  Add("nl", "Hoi")

When you ask for the translation of a certain message key, the user's language is determined from the "Accept-Language" request header.
Passing the http request pointer to Msg() renders a function to do the key->translation lookup:
	translation := Msg(r)("Hi")

You could include the function returned by Msg() to the FuncMap of your template:
	template.FuncMap{
		"Msg": msg.Msg(r),
	},
And then use the mapped Msg function inside the template:
	{{ Msg "Hi" }} {{ .Name }}
*/
package msg

import (
	"net/http"
	"strings"
)

type message map[string]string

/*
Add stores the translation of the message for the given language.
*/
func (m message) Add(language, translation string) message {
	language = strings.ToLower(language)
	m[language] = translation
	return m
}

var messageStore = make(map[string]message, 500)

var NumLang = 2

/*
New creates a new message, and stores it in memory under the given key.
*/
func New(key string) message {
	m := make(message, NumLang)
	messageStore[key] = m
	return m
}

// LanguageType defines a language.
type LanguageType struct {

	// e.g. "en-us"
	Full string

	// e.g. "en"
	Main string

	// e.g. "us"
	Sub  string
}

var languageCache = make(map[string]LanguageType, 100)

// Language provides the first language in the "Accept-Language" header in the
// given http request.
func Language(r *http.Request) (language LanguageType) {
	acceptLanguage := r.Header.Get("Accept-Language")
	acceptLanguage = strings.ToLower(acceptLanguage)
	if lang, ok := languageCache[acceptLanguage]; ok {
		language = lang
	} else {
		firstLanguage := strings.Split(acceptLanguage, ",")[0] // cut other languages
		firstLanguage = strings.Split(firstLanguage, ";")[0]   // cut the q parameter
		parts := strings.Split(firstLanguage, "-")
		lang = LanguageType{
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

/*
Msg returns a function that can lookup the translation for e certain message key.
The language to use is read from the "Accept-Language" header in the given http request.
*/
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
