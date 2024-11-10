package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

var (
	env     *Environment
	envOnce sync.Once
)

type EnvVar struct {
	Name     string
	Required bool
}

type Environment struct {
	Port string `env:"PORT"`
}

func GetEnv() *Environment {
	envOnce.Do(
		func() {
			var err error
			env, err = loadEnv()
			if err != nil {
				log.Fatal(err)
			}
		},
	)

	return env
}

func loadEnv() (*Environment, error) {
	env := &Environment{}
	v := reflect.ValueOf(env).Elem()
	t := v.Type()

	missingVars := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")

		if tag == "" {
			continue
		}

		parts := splitTag(tag)
		name := parts[0]
		require := len(parts) > 1 && parts[1] == "required"

		value := os.Getenv(name)
		if require && value == "" {
			missingVars = append(missingVars, name)
		}
		v.Field(i).SetString(value)
	}
	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required variables:\n%s", strings.Join(missingVars, "\n"))
	}

	return env, nil
}

func splitTag(tag string) []string {
	result := []string{}
	current := ""
	for _, ch := range tag {
		if ch == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(ch)
		}
	}

	result = append(result, current)

	return result
}
