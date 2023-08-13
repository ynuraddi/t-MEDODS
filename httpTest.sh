# Обычный запрос
echo "Ожидаемый запрос"
response=$(curl --silent --location --request POST 'http://localhost:8080/auth/1')

access_token=$(echo "$response" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
refresh_token=$(echo "$response" | sed -n 's/.*"refresh_token":"\([^"]*\)".*/\1/p')

echo "Response from login: $access_token"

responce=$(curl --location --request POST 'http://localhost:8080/refresh' \
--header "Authorization: Bearer $access_token" \
--header 'Content-Type: application/json' \
--data-raw "{
    \"refresh_token\":\"$refresh_token\"
}")

echo "Response from refresh: $response"
