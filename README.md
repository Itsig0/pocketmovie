# PocketMovie

A small database to keep track of the movies you have seen or want to see.

I made this so I can directly see where I can watch a movie I put on my watchlist. Since I'm only subscribed to a hand full of services.

What you can do with it:

- Add movies you have seen or want to see
- Document wether you own the movie or not
- See where movies are streamed via flatrate in your watchlist. You can add or remove streaming services in your settings. 

## Installation

The project comes in a single binary. You execute it and the web server is running. It will create a data directory to store the sqllite database and the movie posters.

```bash
// just execute the binary
./pocketmovie

// You can also specify the port you want to use
./pocketmovie -p 6000
```
