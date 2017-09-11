/* 
 * Grafeas API
 *
 * An API to insert and retrieve annotations on cloud artifacts.
 *
 * OpenAPI spec version: 0.1
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

// A set of properties that uniquely identify a given Docker image.
type Fingerprint struct {

	// The layer-id of the final layer in the Docker image’s v1 representation. This field can be used as a filter in list requests.
	V1Name string `json:"v1Name,omitempty"`

	// The ordered list of v2 blobs that represent a given image.
	V2Blob []string `json:"v2Blob,omitempty"`

	// The name of the image’s v2 blobs computed via:   [bottom] := v2_blobbottom := sha256(v2_blob[N] + “ ” + v2_name[N+1]) Only the name of the final blob is kept. This field can be used as a filter in list requests. @OutputOnly
	V2Name string `json:"v2Name,omitempty"`
}
