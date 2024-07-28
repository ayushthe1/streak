FROM node:alpine as builder

WORKDIR /app

COPY ./client/package.json ./client/package-lock.json ./

RUN npm install

COPY ./client .

RUN npm run build

FROM node:alpine

WORKDIR /app

# Copy the built frontend from the previous stage
COPY --from=builder /app/build /app/build

RUN npm install -g serve

EXPOSE 4000

# Run serve to serve the built frontend
CMD ["serve", "-s", "build", "-l", "4000"]