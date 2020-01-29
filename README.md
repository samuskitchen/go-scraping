SCRAPING

Sample scraping operation using Golang and CockroachDb (PostgreSQL)


## API ENDPOINTS

### Details by Address
- Path : `/scraping/address={address}`
- Method: `GET`
- Response: `200`

### All Address
- Path : `/scraping`
- Method: `GET`
- Response: `200`

## Required Packages
- Dependency management
    * [dep](https://github.com/golang/dep)
    * [goQuery](https://github.com/PuerkitoBio/goquery)
- Database
    * [CockroachDb](https://github.com/lib/pq)
- Routing
    * [chi](https://github.com/go-chi/chi)
    
## Quick Run Project
First clone the repo then go to scraping folder. After that build your image and run by docker. Make sure you have docker in your machine. 

```
git clone https://github.com/samuskitchen/go-scraping

cd go-scraping

chmod +x run.sh
./run.sh

docker-compose up / docker-compose down