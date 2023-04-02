# Config Package

The **config** package provides a simple configuration struct for the application, including logger instances, email channels, and a validator instance.

### Features

**InfoLogger**: an instance of the `log.Logger` struct that is used to log informational messages.

**ErrorLogger**: an instance of the `log.Logger` struct that is used to log error messages.

**MailChan**: a channel of type `model.Email` that is used to send emails from the application.

**Validator**: an instance of the `validator.Validate` struct that is used to validate struct fields based on tags.

### Usage
To use the `AppConfig` struct in your application, simply import the `config` package and create a new instance of the `AppConfig` struct.

```go
import (
	"github.com/yusuf/track-space/pkg/config"
)

func main() {
	appConfig := &config.AppConfig{
		InfoLogger:  // initialize logger instance,
		ErrorLogger: // initialize logger instance,
		MailChan:    make(chan model.Email),
		Validator:   // initialize validator instance,
	}
  

	// ...
}
```

### Dependencies

The `config` package depends on the following external packages:

* **github.com/go-playground/validator/v10**: used to validate struct fields based on tags.
* **github.com/yusuf/track-space/pkg/model**: used to define the model.Email struct that is used in the MailChan channel.


