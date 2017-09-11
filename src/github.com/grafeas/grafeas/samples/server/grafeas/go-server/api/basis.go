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

// Basis describes the base image portion (Note) of the DockerImage relationship.  Linked occurrences are derived from this or an equivalent image via:   FROM <Basis.resource_url> Or an equivalent reference, e.g. a tag of the resource_url.
type Basis struct {

	// The resource_url for the resource representing the basis of associated occurrence images.
	ResourceUrl string `json:"resourceUrl,omitempty"`

	// The fingerprint of the base image
	Fingerprint Fingerprint `json:"fingerprint,omitempty"`
}
