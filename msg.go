package msg

import (
  "strings"
)

type translations map[string]string

var messages = make(map[string]translations)

// var m, a = msg.Init()
// m("Hello")
// a("en", "Hello, world")
// a("nl", "Hallo wereld")
// m("Bye")
// a("en", "Farewell, cruel world")
// a("nl", "Vaarwel, wrede wereld")
func Init () (func (string), func (string, string)) {
  var message *translations
  createMessage := func (key string) {
    t := make(translations)
    messages[key] = t
    message = &t
  }
  addTranslation := func (language string, translation string) {
    (*message)[language] = translation
  }
  return createMessage, addTranslation
}

type languageType struct {
  Full string
  Main string
  Sub string
}

var languages = map[string]languageType{}

// You can pass the value of the Accept-Language http header
// TODO: be more appreciative to the languages listed in the Accept-Language header;
//   currently only the language first listed is considered
func Language (accept_language string) (lang languageType) {
  var ok bool
  if lang, ok = languages[accept_language]; !(ok) {
    first_language := strings.Split(accept_language, ",")[0] // cut other languages
    first_language = strings.Split(first_language, ";")[0] // cut the q parameter
    parts := strings.Split(first_language, "-")
    lang = languageType{
      Full: first_language,
      Main: parts[0],
    }
    if len(parts) > 1 {
      lang.Sub = parts[1]
    }
    languages[accept_language] = lang
  }
  return
}

func Msg (lang languageType, key string) (value string) {
  var ok bool
  if value, ok = messages[key][lang.Main]; !(ok) {
    if value, ok = messages[key][lang.Sub]; !(ok) {
      if value, ok = messages[key][lang.Full]; !(ok) {
        value = "X-" + key
      }
    }
  }
  return
}
