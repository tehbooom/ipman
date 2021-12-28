# Project Name Generator

API  where you submit a GET request which responds with a random project name in the form of ADJECTIVE-NOUN

## Dockerfile

Build the docker file with `docker build -t words-db ./`

Then run the container with `docker run -d --name words-db -p 5432:5432 words-db`

## Setting enviroment variables

set `DB_USERNAME`,  `DB_PASSWORD`, and `DB_NAME`

```bash
export DB_USERNAME=postgres
export DB_PASSWORD=<password in docekrfile>
export DB_NAME=words
```
## Nouns

List of 15,000 nouns

[Nouns](https://greenopolis.com/list-of-nouns/)

## Adjectives

List of 15,000 adjectives

[Adjectives](https://greenopolis.com/adjectives-list/)
