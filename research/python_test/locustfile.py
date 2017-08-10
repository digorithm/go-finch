from locust import HttpLocust, TaskSet, task

class UserBehavior(TaskSet):
    @task
    def index(self):
        t = 0
        t = t + 1
        print t
        self.client.get("/")

    @task
    def profile(self):
        self.client.get("/recipes/1")

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 0
    max_wait = 0
