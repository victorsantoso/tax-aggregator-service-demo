### Tax Aggregator Service

this is a service to process tax calculation from cryptocurrencies company, this service will be deployed as a microservice service to handle monolithic calculation. This service also can be changed to cloud function to handle tax transaction, since it's quite rare to use this service to cut down the deployment cost.


### Running Application Server

```bash
    go run app/main.go start -p 3000 -c ./config/config.json
```

### Mockery Generate
```
    mockery --keeptree --all
```