# Simple Strava Oauth

A simple [Fiber application](https://docs.gofiber.io/) allowing for creation of [Strava](https://developers.strava.com/docs/getting-started/) access tokens needed for advanced [api](http://developers.strava.com/docs/reference/) usage.

- [Simple Strava Oauth](#simple-strava-oauth)
  - [Setup](#setup)
    - [Register with Strava](#register-with-strava)
    - [Clone the repo](#clone-the-repo)
    - [Environment](#environment)
    - [Build the Binary or Run the Server](#build-the-binary-or-run-the-server)
    - [Update Strava Applications settings](#update-strava-applications-settings)
  - [TODO](#todo)

## Setup

### Register with Strava

1. If you have not already, sign up for a [Strava account](https://www.strava.com/register).
2. After you are logged in, [create an app](https://www.strava.com/settings/api).

### Clone the repo

    git@github.com:grahamplata/strava-oauth.git && cd strava-oauth

### Environment

Create a .env file or export the following variables using the supplied crendentials in your [settings](https://www.strava.com/settings/api)

    ENVIRONMENT=local
    PORT=3000
    STRAVA_CLIENT_ID=someclient
    STRAVA_SECRET=somesecret
    STRAVA_REDIRECT_URI=http://localhost:3000/strava-oauth
    STRAVA_SCOPE=read_all

### Build the Binary or Run the Server

    go build -o bin/strava-oauth -v .
    bin/strava-oauth
    # or
    go run main.go

### Update Strava Applications settings

Update your [Authorization Callback Domain](https://www.strava.com/settings/api): When building your app, change “Authorization Callback Domain” to localhost or any domain. When taking your app live, change “Authorization Callback Domain” to a real domain.

## TODO

- Add heroku button
- Add tests
- Do Basic Cleanup and Refactor
