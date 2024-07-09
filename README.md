# Platform 48 functions toolkit
This is a library of Go functions which you might find useful when writing and/or debugging JellyFaas functions. 

## How to use
### Basics

This toolkit library works by creating a context object which contains information about the request which invoked this function. This object is the base of this library. 

To create it simply call the ``toolkit.FuncCtx(w, r)`` function with the request reader and response writer arguments of your functions.

```golang
import (
	tk "github.com/Platform48/p48-toolkit"
)

func yourJellyFaasFunction(w http.ResponseWriter, r *http.Request) {
	ctx := tk.FuncCtx(w, r)
	
	// The rest of your function
}
```

When the ctx object is created it automatically assings a span id to your request, and generates a ``context.Context`` object. You can access them through the ``ctx.SpanId`` and ``ctx.Context`` fields.

### Request information
You can then read the body of the request through the ``ctx.GetBody()`` and ``ctx.GetJsonBody(&result)`` methods.

```golang
bytes, err := ctx.GetBody() //  Reads Content-Length bytes from the request's body and returns them as a byte array.

var result yourResponseObject
err = ctx.GetJsonBody(&result)  //  Reads Content-Length bytes from the request's body and deserializes them into the given object through json.Unmarshal. This method doesn't check the Content-Type.
```

You can also check the query parameters and headers of the request through the ``ctx.HasParameter(name)``, ``ctx.GetParameter(name)``, and ``ctx.GetHeader(name)`` methods.
```golang
var exists bool = ctx.HasParameter("param") //  Returns true if the given query parameter exists in the url of the request, false otherwise

var queryValue string = ctx.GetParameter("param") //  Returns the value of the given query parameter as a string

var headerValue string = ctx.GetHeader("Content-Type") //  Returns the value of the given header as a string
```

You can also access the request information by manually accessing the request reader through the ``ctx.Request`` field.

### Logging

This toolkit library also adds extra information to your log messages such as the span id of the request, and the file the log statement was called in.

Logging is done through calling the ``ctx.Debug(msg)``, ``ctx.Info(msg)``, ``ctx.Warn(msg)``, and ``ctx.Error(msg)`` methods. 
These methods log the given message to the console with the Debug, Info, Warn and Error levels respectively. You can also use the ``ctx.Log(level, message)`` method to specify the log level as an argument.

Each of these methods also have an f variant (e.g.: ``ctx.Infof(format, ...args)``) which allows for formatting your log messages with extra values. For how to use the format versions of the functions [see this link](https://pkg.go.dev/fmt#hdr-Printing).

```golang
ctx.Debug("Debug level message")
ctx.Info("Info level message")
ctx.Warn("Warn level message")
ctx.Error("Error level message")

ctx.Log(tk.LogLevelDebug, " ") //  You can use tk.LogLevelDebug, tk.LogLevelInfo, tk.LogLevelWarn, and tk.LogLevelError to indicate the log level

ctx.Debugf("Number %v", 1)
ctx.Infof("Number %v", 2)
ctx.Warnf("Number %v", 3)
ctx.Errorf("Number %v", 4)

ctx.Logf(tk.LogLevelWarn, "Number %s", "5")
```



Sample output from the above example

```text
12:00PM DBG [5Wc6d2_Sg] Debug level message
12:00PM INF [5Wc6d2_Sg] Info level message
12:00PM WRN [5Wc6d2_Sg] Warn level message
12:00PM ERR [5Wc6d2_Sg] Error level message
12:00PM DBG [5Wc6d2_Sg]  
12:00PM DBG [5Wc6d2_Sg] Number 1
12:00PM INF [5Wc6d2_Sg] Number 2
12:00PM WRN [5Wc6d2_Sg] Number 3
12:00PM ERR [5Wc6d2_Sg] Number 4
12:00PM WRN [5Wc6d2_Sg] Number 5
```
When the function is deployed, the log messages will include extra information in the ``jsonPayload`` object in every log. This makes the logs more readable, while still allowing you to view the detailed information about each message.

You can also access the underlying ``zerolog.Logger`` object to change how the messages are shown, or to print out messages with extra information added to their ``jsonPayload``. This can be done by accessing the ``ctx.Logger`` field.

For more information about the ``zerolog.Logger`` object see [the zerolog repo](https://github.com/rs/zerolog/) and [this link](https://pkg.go.dev/github.com/rs/zerolog).

### Building your responses

When you are done with processing your request, the ctx object will provide you with methods for easier building and logging your responses.
You can call the ``ctx.OkResponse(format, bytes)`` method with the Content-Type, and the bytes for the response. This method automatically builds and sends the response for you.

This method takes in bytes as input, however most of your responses are likely to be in the Json format. Fortunately, you can use the ``ctx.OkResponseJson(obj)`` method to automatically serialize the given struct with the json.Marshal converter, and the ``application/json; charset=utf-8`` Content Type.
If you don't want to define a struct for every possible response, you can also use the ``toolkit.Json`` object to define your objects inline.

You can also add additional headers to your response with the ``ctx.SetResponseHeader(name, value)`` method.

```golang
bytes := make([]byte, 5)
ctx.OkResponse("application/octet-stream", bytes)

type SampleData struct {
    Data   string `json:"data"`
}

ctx.OkResponseJson(SampleData{Data:"Hello, World!"})

//  Alternatively you can call

ctx.OkResponseJson(tk.Json{"data":"Hello, World!"})
//  Ginkgo misinterprets tk.Json as an map[string]interface{}, which causes tests to fail. When comparing tk.Json in your tests you should call .AsMap() on the tk.Json object you are comparing your response to prevent this bug from occurring.
//   Expect(res.Data).To(Equal(tk.Json{"Foo": "Bar", "Heh": 1234.}.AsMap()))
```

The ctx object also has methods for returning and logging invalid responses and errors. You can use the ``ctx.FailResponse(code, message)`` and ``ctx.ErrResponse(code, error, message)`` methods to automatically build and send a Json object with the given error code and a message for the user. Both methods also add a message to the logs, and the ``ctx.ErrResponse`` method also accepts an error object as an argument, which will also be added to the logs.

```golang

ctx.FailResponse(http.StatusBadRequest, "Error message for the user")   //  Useful for when the request contains a problem, but an error wasn't thrown

ctx.ErrResponse(http.StatusInternalServerError, errors.New("Error thrown by another function"), "Error message for the user")    //  Useful for when an error occurred somewhere, and the request cannot be processed.
```
