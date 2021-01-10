![CI](https://github.com/vleedev/smtp-relay-rabbitmq/workflows/CI/badge.svg)
# Description
This app is used to listen to a queue of rabbitmq and send an email via SMTP connection
## Using
Use container image `vleedev/smtp-relay-rabbitmq:latest`

### Set environment variables
    RABBITMQ_URL: your rabbitmq url service
    QUEUE_NAME: a queue name to listen on that service
    SMTP_HOSTNAME: smtp host
    SMTP_PORT: smtp port
    SMTP_USERNAME: smtp username
    SMTP_PASSWORD: smtp password
    SMTP_DEFAULT_EMAIL: default email if sender is not set

### Docker-compose example
    version: '3.7'
    services:
        smtp-relay-rabbitmq:
            image: vleedev/smtp-relay-rabbitmq:latest
            environment:
                RABBITMQ_URL: amqps://wrfizrrb:Fz3wHpwrwLL39J31ekLR8kR_bBT0s8ruv@cougar.rmq.cloudamqp.com/wrfizrrb
                QUEUE_NAME: my_project_app_email_queue
                SMTP_HOSTNAME: email-smtp.us-west-2.amazonaws.com
                SMTP_PORT: 587
                SMTP_USERNAME: AKIATXPJWI7G5O3LIAFX
                SMTP_PASSWORD: BNP9c2nm8taH7dpaZUXLpM7MP0OO6vQlSjfRX1Yk43vb
                SMTP_DEFAULT_EMAIL: me@vlee.dev
### Test service by publishing your contents to queue
You must prepare your data by encoding json structure as below
    
    {
        "from": "me@vlee.dev",
        "to": [
            "admin@google.com",
            "ad@facebook.com"
        ],
        "subject": "My subject",
        "body_type": "text/html",
        "body": "<html><body><p>This one is a test email from smtp-relay-rabbitmq</p></body></html>",
        "attachment": [
            "https://i.imgur.com/UbUQWHO.jpeg"
        ]
    }

Note: publish to your queue with `content-type` = `application/json`

### Sign up a cloud amqp queue
Website: https://www.cloudamqp.com/

### Code Boilerplate
Look for your favourite language and implement it
https://www.rabbitmq.com/getstarted.html
