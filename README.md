# NGORMQ
### Nginx + Go + gORM + rabbitMQ

A small project written in Go for learning how to use NGINX, GORM, and RabbitMQ. Why?

- To learn how to use NGINX as a load balancer (there are 2 HTTP server replicas on `docker-compose.yml`)
- To learn how to use an ORM, because my current knowledge is using raw SQL (based on github stats, [GORM](https://github.com/go-gorm/gorm) is more used compared to [SQLX](https://github.com/jmoiron/sqlx))
- To learn how to use RabbitMQ, because my current knowledge is using NSQ as a message broker (RabbitMQ is way more widely used)

Things I could do to make this codebase better:
- Separate handlers into message queue handlers and HTTP handlers
- Separation of concerns (handlers -> usecases -> repositories)
- Add unit tests