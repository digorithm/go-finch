# Just a few notes and tips so I won't get lost in the future


## Running the experiment

- Always run the meal planner API first
- To run only one task in the debugging mode, run: `go build -o a.out boomer.go request_generator.go && ./a.out --run-tasks <TaskName>`- To run the environment completely run the locust master: `locust -f dummy.py --master --master-bind-host=127.0.0.1 --master-bind-port=5557` then run the GO slave: `go build -o a.out boomer.go request_generator.go && ./a.out --master-host=127.0.0.1 --master-port=5557`

## Important notes

- Increase max connections by editing /etc/postgresql/9.3/main/postgresql.conf
- restart by running /etc/init.d/postgresql restart
