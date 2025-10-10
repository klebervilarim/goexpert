
#!/bin/bash
IP="192.168.0.16"
TOKEN="abc123"

echo "Testando limite por IP..."
for i in {1..6}; do
    curl -s -w "%{http_code}\n" http://$IP:8080/test &
done
wait
echo "-------------------------"

echo "Testando limite por Token..."
for i in {1..11}; do
    curl -s -w "%{http_code}\n" -H "API_KEY: $TOKEN" http://$IP:8080/test &
done
wait
echo "-------------------------"

