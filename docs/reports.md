# Reports

## Sales Summary

```bash
curl \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/reports/sales-summary?restaurant_id=$RESTAURANT_ID"
```

## Low Stock CSV

```bash
curl \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/reports/low-stock.csv?restaurant_id=$RESTAURANT_ID"
```
