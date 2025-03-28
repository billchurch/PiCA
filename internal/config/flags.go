package config

import (
	"flag"
	"reflect"
	"strconv"
)

// RegisterFlags registers command-line flags based on the Config struct tags
func (cfg *Config) RegisterFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("pica", flag.ExitOnError)

	t := reflect.TypeOf(*cfg)
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		flagName := field.Tag.Get("flag")
		description := field.Tag.Get("desc")
		if description == "" {
			description = "Configuration for " + field.Name
		}

		if flagName != "" {
			fieldValue := v.Field(i)
			switch fieldValue.Kind() {
			case reflect.String:
				current := fieldValue.String()
				fs.StringVar((*string)(fieldValue.Addr().UnsafePointer()), flagName, current, description)
			case reflect.Int:
				current := int(fieldValue.Int())
				fs.IntVar((*int)(fieldValue.Addr().UnsafePointer()), flagName, current, description)
			case reflect.Bool:
				current := fieldValue.Bool()
				fs.BoolVar((*bool)(fieldValue.Addr().UnsafePointer()), flagName, current, description)
			case reflect.Float64:
				current := fieldValue.Float()
				fs.Float64Var((*float64)(fieldValue.Addr().UnsafePointer()), flagName, current, description)
			}
		}
	}

	return fs
}

// ParseFlags parses command-line flags and updates the config
func (cfg *Config) ParseFlags(args []string) error {
	fs := cfg.RegisterFlags()
	return fs.Parse(args)
}

// LoadFlagDefaults loads default values from struct tags into flags
func (cfg *Config) LoadFlagDefaults(fs *flag.FlagSet) {
	t := reflect.TypeOf(*cfg)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		flagName := field.Tag.Get("flag")
		defaultVal := field.Tag.Get("default")

		if flagName != "" && defaultVal != "" {
			flag := fs.Lookup(flagName)
			if flag != nil {
				flag.DefValue = defaultVal
				switch field.Type.Kind() {
				case reflect.String:
					flag.Value.Set(defaultVal)
				case reflect.Int:
					if intVal, err := strconv.Atoi(defaultVal); err == nil {
						flag.Value.Set(strconv.Itoa(intVal))
					}
				case reflect.Bool:
					if boolVal, err := strconv.ParseBool(defaultVal); err == nil {
						flag.Value.Set(strconv.FormatBool(boolVal))
					}
				case reflect.Float64:
					if floatVal, err := strconv.ParseFloat(defaultVal, 64); err == nil {
						flag.Value.Set(strconv.FormatFloat(floatVal, 'f', -1, 64))
					}
				}
			}
		}
	}
}
