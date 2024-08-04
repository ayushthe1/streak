# Streak

This is scalable chat website built using Golang, designed to provide real-time messaging, file sharing, semantic search, and AI chatbot functionalities. The application is containerized using Docker Compose and deployed on Google Cloud Platform (GCP). 


## Architecture

![Project Architecture](https://github.com/user-attachments/assets/0554dcd2-2860-4859-9769-c6b287423a94)

## Features

- **Real-time Messaging**: Powered by Golang on the backend, with Redis and RabbitMQ as the messaging queue, ensuring efficient pub/sub messaging.
- **Frontend**: Built with React, providing a responsive and interactive user interface.
- **Database**: PostgreSQL is used as the primary database for storing chat messages and user data.
- **Message Persistence**: Kafka is used for writing messages to the database, ensuring reliable and scalable message persistence.
- **Semantic Search**: Weaviate is integrated to enable semantic search in chat history, enhancing user experience.
- **File Sharing**: AWS Go SDK is used to enable file sharing between users using S3 buckets.
- **AI Chatbot**: Integrated with Dialogflow and Gemini to provide AI-driven chat functionalities.
- **Monitoring**: Prometheus and Grafana are used for monitoring and visualizing metrics.
- **Reverse Proxy**: Caddy is used as a reverse proxy to forward incoming requests to different services.

## Technologies Used

- **Backend**: Golang, Fiber framework
- **Frontend**: React, Chakra UI
- **Database**: PostgreSQL
- **Message Queue**: RabbitMQ, Redis
- **Message Persistence**: Kafka
- **Search**: Weaviate
- **File Storage**: AWS S3
- **AI Chatbot**: Dialogflow, Gemini
- **Containerization**: Docker Compose
- **Deployment**: Google Cloud Platform (GCP)
- **Monitoring**: Prometheus, Grafana
- **Reverse Proxy**: Caddy
