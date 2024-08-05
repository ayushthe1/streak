# Streak

This is scalable chat website built using Golang, designed to provide real-time messaging, file sharing, semantic search, and AI chatbot functionalities. The application is containerized using Docker Compose and deployed on Google Cloud Platform (GCP). 
Try it out here : https://streak.ayushsharma.co.in/


## Architecture

![Project Architecture](https://github.com/user-attachments/assets/9c6bf006-80c9-4919-8623-73bfd5f7dab1)

## Chat Interface
![Chat interface](https://github.com/user-attachments/assets/dc481370-b226-45fe-96ac-19695eb9ac4d)


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
- **Database**: PostgreSQL, Redis (will add caching later)
- **Message Queue**: RabbitMQ
- **Message Persistence**: Kafka
- **Semantic Search**: Weaviate
- **File Storage**: AWS S3 (using their Go SDK)
- **AI Chatbot**: Dialogflow ((using their Go SDK), Gemini
- **Containerization**: Docker Compose
- **Deployment**: Google Cloud Platform (GCP)
- **Monitoring**: Prometheus, Grafana
- **Reverse Proxy**: Caddy

## Grafana dashboard
  ![grafana_dashboard](https://github.com/user-attachments/assets/1684e67a-5472-4f57-abaf-9c7b38227c60)


#### More features will be added soon. All contributions are welcome :)
