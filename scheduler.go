package main

import (
	"database/sql"
	"log"
	"time"
)

func startScheduler(db sql.DB, tasksPerLoop int) {
	errCount := 0

	for {
		now := time.Now().UTC()

		tasks, err := getRunnableTasks(db, now.Add(5*time.Second), tasksPerLoop)

		if err != nil {
			errCount++

			if errCount > 10 {
				log.Fatal("Too many errors occurred, scheduler is stopping")
			}
		} else if len(tasks) > 0 {
			for _, task := range tasks {
				runTask(task)

				if err := updateTaskNextFireAt(db, task.Id, now.Add(30*time.Second), false); err != nil {
					log.Print("Update task failed ", task.Id, " ", task.Rule, " ", err)
				}
			}
		} else {
			log.Print("No tasks to run")
		}

		<-time.After(5 * time.Second)
	}
}

func runTask(task DbTask) *error {
	log.Print("Running task ", task.Id, " ", task.Rule)
	return nil
}
