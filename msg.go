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

// Call Lang to create the lookup function for the concerning language to use in your templates
// You can pass the value of the Accept-Language http header
// TODO: be more appreciative to the languages listed in the Accept-Language header;
// currently only the main group of the language first listed is considered
func Language (language string) func (string) string {
  language = strings.Split(language, ",")[0]
  language = strings.Split(language, ";")[0]
  language = strings.Split(language, "-")[0]
  return func (key string) string {
    value, ok := messages[key][language]
    if ok {return value}
    return "X-" + key
  }
}
