package main

import (
	"flag"
	"log"
	"net"
	"runtime"

	"github.com/kataras/iris"
	"github.com/valyala/tcplisten"
)

var (
	addr = flag.String("addr", ":8080", "TCP address to listen to")
)

func getListener() net.Listener {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)

	listenerCfg := tcplisten.Config{
		ReusePort:   true,
		DeferAccept: true,
		FastOpen:    true,
	}

	l, err := listenerCfg.NewListener("tcp4", *addr)

	if err != nil {
		log.Fatal(err)
	}

	return l
}

func main() {
	db, err := initDb("tempo.db")

	if err != nil {
		log.Fatal(err)
	}

	app := iris.New()

	app.OnAnyErrorCode(func(ctx iris.Context) {
		if msg := ctx.Value("error"); msg != nil {
			ctx.JSON(iris.Map{"error": ctx.Value("error")})
		} else {
			ctx.JSON(iris.Map{"error": "unknown"})
		}
	})

	app.Get("/tasks", func(ctx iris.Context) {
		tasks, err := getTasks(*db)

		if err != nil {
			ctx.Values().Set("error", err.Error())
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		ctx.JSON(tasks)
	})

	app.Get("/tasks/{id:string}", func(ctx iris.Context) {
		id := ctx.Params().Get("id")

		task, err := getTask(*db, id)

		if err != nil {
			ctx.Values().Set("error", err.Error())
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		} else if task == nil {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}

		ctx.JSON(task)
	})

	app.Post("/tasks", func(ctx iris.Context) {
		task := &Task{}

		if err := ctx.ReadJSON(task); err != nil {
			ctx.Values().Set("error", err.Error())
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		dbTask, err := addTask(*db, *task)

		if err != nil {
			ctx.Values().Set("error", err.Error())
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		ctx.StatusCode(201)
		ctx.JSON(*dbTask)
	})

	flag.Parse()

	go startScheduler(*db, 10)

	app.Run(iris.Listener(getListener()))
}
