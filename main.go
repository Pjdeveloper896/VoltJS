
// main.go
package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/dop251/goja"
)

var vm *goja.Runtime
var intervals = make(map[int]*time.Ticker)
var timeouts = make(map[int]*time.Timer)
var nextID = 1

func setupConsole(vm *goja.Runtime) {
    console := vm.NewObject()
    console.Set("log", func(call goja.FunctionCall) goja.Value {
        for _, arg := range call.Arguments {
            fmt.Print(arg.Export(), " ")
        }
        fmt.Println()
        return goja.Undefined()
    })
    vm.Set("console", console)
}

func setupTimers(vm *goja.Runtime) {
    vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
        fn, _ := goja.AssertFunction(call.Arguments[0])
        delay := time.Duration(call.Arguments[1].ToInteger()) * time.Millisecond
        id := nextID
        nextID++
        timeouts[id] = time.AfterFunc(delay, func() {
            _, _ = fn(goja.Undefined())
        })
        return vm.ToValue(id)
    })

    vm.Set("clearTimeout", func(call goja.FunctionCall) goja.Value {
        id := int(call.Arguments[0].ToInteger())
        if t, ok := timeouts[id]; ok {
            t.Stop()
            delete(timeouts, id)
        }
        return goja.Undefined()
    })

    vm.Set("setInterval", func(call goja.FunctionCall) goja.Value {
        fn, _ := goja.AssertFunction(call.Arguments[0])
        delay := time.Duration(call.Arguments[1].ToInteger()) * time.Millisecond
        ticker := time.NewTicker(delay)
        id := nextID
        nextID++
        intervals[id] = ticker

        go func() {
            for range ticker.C {
                _, _ = fn(goja.Undefined())
            }
        }()

        return vm.ToValue(id)
    })

    vm.Set("clearInterval", func(call goja.FunctionCall) goja.Value {
        id := int(call.Arguments[0].ToInteger())
        if t, ok := intervals[id]; ok {
            t.Stop()
            delete(intervals, id)
        }
        return goja.Undefined()
    })
}

func setupFS(vm *goja.Runtime) {
    fs := vm.NewObject()
    fs.Set("readFileSync", func(call goja.FunctionCall) goja.Value {
        filename := call.Arguments[0].String()
        data, err := os.ReadFile(filename)
        if err != nil {
            panic(vm.ToValue(err.Error()))
        }
        return vm.ToValue(string(data))
    })
    fs.Set("writeFileSync", func(call goja.FunctionCall) goja.Value {
        filename := call.Arguments[0].String()
        content := call.Arguments[1].String()
        err := os.WriteFile(filename, []byte(content), 0644)
        if err != nil {
            panic(vm.ToValue(err.Error()))
        }
        return goja.Undefined()
    })
    vm.Set("fs", fs)
}

func setupHTTP(vm *goja.Runtime) {
    httpMod := vm.NewObject()
    httpMod.Set("createServer", func(call goja.FunctionCall) goja.Value {
        handler, _ := goja.AssertFunction(call.Arguments[0])
        go func() {
            http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                req := vm.NewObject()
                req.Set("url", r.URL.Path)
                req.Set("method", r.Method)

                res := vm.NewObject()
                res.Set("end", func(call goja.FunctionCall) goja.Value {
                    if len(call.Arguments) > 0 {
                        w.Write([]byte(call.Arguments[0].String()))
                    }
                    return goja.Undefined()
                })

                _, err := handler(goja.Undefined(), req, res)
                if err != nil {
                    fmt.Println("HTTP handler error:", err)
                }
            }))
        }()
        return goja.Undefined()
    })
    vm.Set("http", httpMod)
}

func setupProcess(vm *goja.Runtime) {
    proc := vm.NewObject()
    proc.Set("argv", os.Args)
    proc.Set("cwd", func(goja.FunctionCall) goja.Value {
        dir, _ := os.Getwd()
        return vm.ToValue(dir)
    })
    vm.Set("process", proc)
}

func setupRequire(vm *goja.Runtime, basePath string) {
    vm.Set("require", func(call goja.FunctionCall) goja.Value {
        modulePath := call.Arguments[0].String()
        fullPath := filepath.Join(basePath, "modules", modulePath+".js")
        data, err := os.ReadFile(fullPath)
        if err != nil {
            panic(vm.ToValue("Module not found: " + fullPath))
        }

        newVM := goja.New()
        setupConsole(newVM)
        setupTimers(newVM)
        setupFS(newVM)
        setupHTTP(newVM)
        setupProcess(newVM)

        _, err = newVM.RunString(string(data))
        if err != nil {
            panic(vm.ToValue("Error in module: " + err.Error()))
        }

        return newVM.Get("exports")
    })
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <script.js>")
        return
    }

    script := os.Args[1]
    data, err := os.ReadFile(script)
    if err != nil {
        fmt.Println("Error reading script:", err)
        return
    }

    vm = goja.New()
    setupConsole(vm)
    setupTimers(vm)
    setupFS(vm)
    setupHTTP(vm)
    setupProcess(vm)
    setupRequire(vm, ".")

    _, err = vm.RunString(string(data))
    if err != nil {
        fmt.Println("JavaScript error:", err)
    }

    select {}
}
