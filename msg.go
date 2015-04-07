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

// Call Lang to create the lookup function for the concerning language to use in your templates
// You can pass the value of the Accept-Language http header
// TODO: be more appreciative to the languages listed in the Accept-Language header;
// currently only the main group of the language first listed is considered
func Language (accept_language string) (msg func (string) string, lang languageType) {
  accept_language = strings.Split(accept_language, ",")[0]
  accept_language = strings.Split(accept_language, ";")[0]
  first_language := strings.Split(accept_language, "-")
  lang = languageType{
    Full: accept_language,
    Main: first_language[0],
  }
  if len(first_language) > 1 {
    lang.Sub = first_language[1]
  }
  msg = func (key string) (value string) {
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
  return
}
