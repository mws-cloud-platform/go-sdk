package sensitive_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"go.mws.cloud/go-sdk/pkg/apimodels/sensitive"
)

func Example() {
	secret := sensitive.New("secret")
	fmt.Println("redacted:", secret)
	fmt.Println("raw:", secret.Value())
	// Output:
	// redacted: ****
	// raw: secret
}

func Example_customFormat() {
	secret := sensitive.New("123456789", sensitive.WithFormat(func(v string) string {
		return v[0:3] + "..."
	}))
	fmt.Println(secret)
	// Output:
	// 123...
}

func Example_text() {
	secret := sensitive.New("secret")

	rawText, err := secret.MarshalText()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(rawText))
	// Output: ****
}

func Example_json() {
	secret := sensitive.New("secret")

	rawJSON, err := json.Marshal(secret)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(rawJSON))
	// Output: "****"
}

func Example_yaml() {
	secret := sensitive.New("secret")

	rawYAML, err := yaml.Marshal(secret)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(rawYAML))
	// Output: '****'
}

func Example_slog() {
	secret := sensitive.New("secret")

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey { // hide time attr
				return slog.Attr{}
			}
			return attr
		},
	}))

	slogger.Info("slog", "secret", secret)
	// Output:
	// {"level":"INFO","msg":"slog","secret":"****"}
}

func Example_zap() {
	secret := sensitive.New("secret")
	zapLogger := zap.NewExample()
	zapLogger.Info("zap", zap.Any("secret", secret))
	zapLogger.Sync()
	// Output:
	// {"level":"info","msg":"zap","secret":{"value":"****"}}
}
