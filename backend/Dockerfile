# Build the API and scheduler binaries in a dedicated Go image.
FROM golang:1.25-alpine AS build

WORKDIR /src

# Download dependencies separately so Docker can reuse this cached layer.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build static Linux binaries.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/scheduler ./cmd/scheduler

# Use a small runtime image that does not include the Go toolchain.
FROM alpine:3.22

# Run the application as an unprivileged user.
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy only the compiled binaries from the build stage.
COPY --from=build /out/api /app/api
COPY --from=build /out/scheduler /app/scheduler

USER appuser

# Document the API's default container port.
EXPOSE 8080

# Start the API by default; Compose can override this for the scheduler.
CMD [ "/app/api" ]
