

ghz --insecure \
  --proto=./proto/greet.proto \
  --call unary.Greeter.greet \
  -c 50 \
  -n 100000 \
  -d '{"name":"Joe"}' \
  0.0.0.0:3456






