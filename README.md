# Genesis test

## How to run

1. Install and run [Docker](https://www.docker.com/)
2. Run `docker build -t genesis .` for build docker image
3. For the next step, you need to get the app password in your Google account [Tutorial](https://support.google.com/accounts/answer/185833?hl=ru)
4. Run `docker run --rm -e PORT=8000 -e EMAIL=YOUR_EMAIL -e PASSWORD=YOUR_EMAIL_APP_PASSWORD -p 8000:8000 --name genesis_container genesis` for starting docker container
5. Go to http://localhost:8000/swagger/index.html
