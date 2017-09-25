create table tasks
(
  id VARCHAR(64) not null
    primary key,
  rule VARCHAR(64) not null,
  timeZone VARCHAR(32) not null,
  epsilon INT default 60 not null,
  nextFireAt DATETIME,
  nextRetryAt DATETIME,
  maxRetries INT default 3 not null,
  completed BOOLEAN default FALSE not null
)
;

create unique index task_id_uindex
  on tasks (id)
;

create table tasksRuns
(
  id INTEGER not null,
  taskId INT not null
    constraint tasksRuns_tasks_id_fk
    references tasks (id)
      on delete cascade,
  status INT not null,
  runAt DATETIME not null,
  finishedAt DATETIME
)
;

create unique index tasksRuns_id_uindex
  on tasksRuns (id)
;

