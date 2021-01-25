# Copyright 2021 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "add routes"

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRoute.json | jq .

# idempotency test

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRoute.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteBlankID.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleFilterRoute.json | jq .

echo "invalid routes"

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteBadName.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteBadPluginName.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteNoReceiver.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteNoSender.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteNoApp.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteNoOrg.json | jq .

curl -X POST http://localhost:3000/ears/v1/routes --data @testdata/simpleRouteNoUser.json | jq .

echo "get routes"

curl -X GET http://localhost:3000/ears/v1/routes | jq .

curl -X GET http://localhost:3000/ears/v1/routes/r123 | jq .

curl -X GET http://localhost:3000/ears/v1/routes/foo | jq .

echo "delete routes"

curl -X DELETE http://localhost:3000/ears/v1/routes/94d5eff28471968e9bd946bc9db27847  | jq .

curl -X DELETE http://localhost:3000/ears/v1/routes/r123  | jq .

curl -X DELETE http://localhost:3000/ears/v1/routes/f123  | jq .

# idempotency test

curl -X DELETE http://localhost:3000/ears/v1/routes/r123  | jq .

echo "get routes"

curl -X GET http://localhost:3000/ears/v1/routes | jq .
