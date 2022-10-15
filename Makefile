# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all:
	go build -o frontend/frontend ./frontend
	go build -o gifcreator/gifcreator ./gifcreator
	go build -o render/render ./render

proto:
	protoc -I ./proto/ \
	--go_out=./proto \
	--go_opt=paths=source_relative \
	--go-grpc_out=./proto \
	--go-grpc_opt=paths=source_relative \
	proto/*.proto

# from upstream
#
#proto: proto/gifcreator.pb.go proto/render.pb.go
#
#proto/gifcreator.pb.go proto/render.pb.go: proto/gifcreator.proto proto/render.proto
#	protoc $^ --go_out=plugins=grpc:.

clean:
	rm -f gifcreator/gifcreator render/render frontend/frontend

.PHONY: all proto clean
