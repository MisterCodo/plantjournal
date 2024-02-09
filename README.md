# Plant Journal

Plant Journal is a pet project meant to help take care of plants.

## How to Use

### Build and Run Go Application

1. Set environment variable `CGO_ENABLED` to 1.
2. Build application with `go build .\cmd\plantjournal\main.go`.
3. Launch the executable with `-h` to see help.

### Build and Run with Docker

1. Run `docker build -t plantjournal .` to build the Docker image.
2. Run `docker run -it -p 8080:8080 plantjournal:latest` to launch the image.
3. Access the application at [http://127.0.0.1:8080/](http://127.0.0.1:8080/)

## Features

- Set plant preferences and details
  - Watering
  - Lighting
  - Fertilizing
  - Toxicity
  - Additional notes
- Pictures
- Track quarantine
- Track pests and pest treatment
- Obituary section

## Out of Scope / Won't do

- User/account management
- Stats/metrics

## Resources

- Icons from [iconoir](https://iconoir.com/).
- CSS framework from [Bulma](https://bulma.io/).
